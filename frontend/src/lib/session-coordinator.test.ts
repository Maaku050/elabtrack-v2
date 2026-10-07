import { afterEach, describe, expect, it, vi } from 'vitest'
import { ApiRequestError } from './api-error'
import { classifyRefreshFailure, parseSessionEvent, SessionCoordinator, SESSION_WAIT_MS, type SessionEvent, type SessionChannel, type SessionLocks, type SessionAttempt } from './session-coordinator'

const now = 1_790_000_000_000
const tab = '11111111-1111-4111-8111-111111111111'
const attempt = '22222222-2222-4222-8222-222222222222'
const base = { version: 1 as const, tabId: tab, attemptId: attempt, epoch: now * 1_000, timestamp: now }
function deferred<T>() {
  let resolve!: (value: T) => void
  const promise = new Promise<T>((r) => { resolve = r })
  return { promise, resolve }
}

// Test ports model independent document runtimes with a shared FIFO browser
// lock service. No production storage, sleeps or fabricated token sharing.
export function coordinationFixture() {
  const channels: SessionChannel[] = []
  const sent: SessionEvent[] = []
  const queue: Array<() => void> = []
  let active = false
  const pump = () => { if (!active) queue.shift()?.() }
  const locks: SessionLocks = {
    request: <T,>(_name: string, { signal }: { mode: 'exclusive'; signal: AbortSignal }, callback: () => Promise<T>) => new Promise<T>((resolve, reject) => {
      let started = false
      const cancel = () => {
        if (started) return
        const index = queue.indexOf(start)
        if (index >= 0) queue.splice(index, 1)
        reject(new DOMException('Aborted', 'AbortError'))
      }
      const start = () => {
        if (signal.aborted) { cancel(); pump(); return }
        started = true; active = true; signal.removeEventListener('abort', cancel)
        void Promise.resolve().then(callback).then(resolve, reject).finally(() => { active = false; pump() })
      }
      signal.addEventListener('abort', cancel, { once: true })
      queue.push(start); pump()
    }),
  }
  const channel = (): SessionChannel => {
    const c: SessionChannel = {
      onmessage: null,
      postMessage: (message) => { sent.push(message); for (const peer of channels) if (peer !== c) peer.onmessage?.({ data: structuredClone(message) } as MessageEvent) },
      close: () => { channels.splice(channels.indexOf(c), 1) },
    }
    channels.push(c)
    return c
  }
  return { locks, channel, sent, channels, crash: () => { active = false; pump() } }
}

const owned: SessionCoordinator[] = []
function coordinator(options: ConstructorParameters<typeof SessionCoordinator>[0] = {}) {
  const c = new SessionCoordinator({ now: () => now, ...options }); owned.push(c); return c
}
afterEach(() => { for (const c of owned.splice(0)) c.dispose(); vi.useRealTimers() })

describe('non-secret session message contract', () => {
  it.each(['refresh-started', 'refresh-succeeded', 'logout', 'session-invalidated', 'auth-state-changed'] as const)('accepts %s', (type) => {
    expect(parseSessionEvent({ ...base, type }, now)).toEqual({ ...base, type })
  })
  it.each(['network', 'server', 'unauthenticated', 'ambiguous-denial', 'forbidden', 'coordination-timeout'] as const)('accepts refresh-failed %s', (failure) => {
    expect(parseSessionEvent({ ...base, type: 'refresh-failed', failure }, now)).not.toBeNull()
  })
  it.each([null, [], 'logout', {}, { ...base, type: 'unknown' }, { ...base, type: 'logout', version: 2 }, { ...base, type: 'logout', tabId: 'token' }, { ...base, type: 'logout', epoch: Infinity }, { ...base, type: 'logout', timestamp: now - 60_001 }, { ...base, type: 'logout', timestamp: now + 5_001 }, { ...base, type: 'refresh-failed', failure: 'SQL reason' }])('rejects malformed or stale input %#', (value) => {
    expect(parseSessionEvent(value, now)).toBeNull()
  })
  it.each(['access_token', 'refresh_token', 'token_hash', 'Authorization', 'Cookie', 'password', 'reset_token', 'user', 'role'])('rejects the extra %s field even on valid types', (key) => {
    expect(parseSessionEvent({ ...base, type: 'logout', [key]: 'synthetic-secret' }, now)).toBeNull()
  })
})

describe('independent tabs with atomic browser ownership', () => {
  it.each([2, 3, 5])('serializes %i tabs, each receiving only its own server token', async (count) => {
    const f = coordinationFixture()
    const tabs = Array.from({ length: count }, () => coordinator({ locks: f.locks, channel: f.channel() }))
    const firstStarted = deferred<void>(), gate = deferred<void>()
    let inFlight = 0, max = 0, cookie = 0
    const memory = Array<string | null>(count).fill(null), used: number[] = []
    const tasks = tabs.map((c, i) => c.run(true, async (a) => {
      max = Math.max(max, ++inFlight); used.push(cookie)
      if (i === 0) { firstStarted.resolve(); await gate.promise }
      cookie++; memory[i] = `synthetic-own-access-${i}`; inFlight--; c.succeeded(a)
    }))
    await firstStarted.promise
    expect(used).toEqual([0]); expect(memory.slice(1).every((v) => v === null)).toBe(true)
    gate.resolve(); await Promise.all(tasks)
    expect(max).toBe(1); expect(used).toEqual(Array.from({ length: count }, (_, i) => i))
    expect(new Set(memory).size).toBe(count)
    expect(JSON.stringify(f.sent)).not.toContain('synthetic-own-access')
    expect(f.sent.filter((m) => m.type === 'refresh-started')).toHaveLength(count)
    expect(f.sent.filter((m) => m.type === 'refresh-succeeded')).toHaveLength(count)
  })
  it.each(['network', 'server'] as const)('owner %s failure releases ownership without peer logout', async (failure) => {
    const f = coordinationFixture(), a = coordinator({ locks: f.locks, channel: f.channel() }), b = coordinator({ locks: f.locks, channel: f.channel() })
    const clear = vi.fn(); b.subscribe(clear)
    const start = deferred<void>(), gate = deferred<void>()
    const first = a.run(true, async (attempt) => { start.resolve(); await gate.promise; a.failed(attempt, failure); throw new ApiRequestError('Failed', failure === 'network' ? 0 : 503) })
    const rejection = expect(first).rejects.toBeInstanceOf(ApiRequestError)
    await start.promise
    const next = b.run(true, async (attempt) => { b.succeeded(attempt); return 'own-access' })
    gate.resolve(); await rejection
    await expect(next).resolves.toBe('own-access'); expect(clear).not.toHaveBeenCalled()
  })
  it('an exclusive authentication denial publishes generic invalidation to every waiting peer', async () => {
    const f = coordinationFixture(), a = coordinator({ locks: f.locks, channel: f.channel() }), b = coordinator({ locks: f.locks, channel: f.channel() }), c = coordinator({ locks: f.locks, channel: f.channel() })
    const bClear = vi.fn(), cClear = vi.fn(); b.subscribe(bClear); c.subscribe(cClear)
    await a.run(true, async (attempt) => a.failed(attempt, 'unauthenticated'))
    expect(bClear).toHaveBeenCalledOnce(); expect(cClear).toHaveBeenCalledOnce()
    expect(f.sent.at(-1)).toMatchObject({ type: 'session-invalidated' })
    expect(f.sent.at(-1)).not.toHaveProperty('reason')
  })
  it('crash without any completion message automatically releases the browser lock', async () => {
    const f = coordinationFixture(), a = coordinator({ locks: f.locks, channel: f.channel() }), b = coordinator({ locks: f.locks, channel: f.channel() })
    const started = deferred<void>()
    void a.run(true, async () => { started.resolve(); return new Promise<void>(() => {}) })
    await started.promise
    const next = b.run(true, async (attempt) => { b.succeeded(attempt); return 'own-access' })
    a.dispose(); f.crash()
    await expect(next).resolves.toBe('own-access')
    expect(f.sent.filter((m) => m.tabId === a.tabId).map((m) => m.type)).toEqual(['refresh-started'])
  })
  it('bounds a suspended live owner wait; retry recovers after its browser lock is released', async () => {
    vi.useFakeTimers()
    const f = coordinationFixture(), a = coordinator({ locks: f.locks }), b = coordinator({ locks: f.locks })
    const started = deferred<void>(), next = vi.fn(async () => 'own-access')
    void a.run(true, async () => { started.resolve(); return new Promise<void>(() => {}) })
    await started.promise
    const waiting = b.run(true, next)
    const rejection = expect(waiting).rejects.toMatchObject({ code: 'COORDINATION_TIMEOUT', status: 0 })
    await vi.advanceTimersByTimeAsync(SESSION_WAIT_MS); await rejection
    expect(next).not.toHaveBeenCalled() // Never steal a live cookie mutation.
    f.crash(); await expect(b.run(true, next)).resolves.toBe('own-access')
  })
  it('does not depend on completion broadcasts when a channel is unavailable', async () => {
    const f = coordinationFixture(), a = coordinator({ locks: f.locks }), b = coordinator({ locks: f.locks })
    await expect(Promise.all([a.run(true, async () => 'A'), b.run(true, async () => 'B')])).resolves.toEqual(['A', 'B'])
  })
  it('without Web Locks a denial is ambiguous and cannot globally invalidate a peer', async () => {
    const f = coordinationFixture(), a = coordinator({ channel: f.channel() }), b = coordinator({ channel: f.channel() })
    const clear = vi.fn(); b.subscribe(clear)
    await a.run(true, async (attempt) => a.failed(attempt, classifyRefreshFailure(new ApiRequestError('Denied', 401), attempt.exclusive)))
    expect(f.sent.at(-1)).toMatchObject({ type: 'refresh-failed', failure: 'ambiguous-denial' })
    expect(clear).not.toHaveBeenCalled()
  })
})

describe('lifecycle and stale-event fencing', () => {
  it('a delayed authority denial retains its request-start epoch and cannot invalidate a newer peer', async () => {
    const f = coordinationFixture(), a = coordinator({ channel: f.channel() }), b = coordinator({ channel: f.channel(), now: () => now + 1 })
    const clear = vi.fn(); b.subscribe(clear)
    const request = a.capture()
    await b.run(true, async (attempt) => b.succeeded(attempt))
    a.invalidate(request)
    expect(clear).not.toHaveBeenCalled()
  })
  it.each(['logout', 'session-invalidated', 'auth-state-changed'] as const)('%s promptly notifies multiple peers without requesting refresh', (type) => {
    const f = coordinationFixture(), a = coordinator({ channel: f.channel() }), b = coordinator({ channel: f.channel() }), c = coordinator({ channel: f.channel() })
    const bClear = vi.fn(), cClear = vi.fn(); b.subscribe(bClear); c.subscribe(cClear)
    a.announce(type)
    expect(bClear).toHaveBeenCalledOnce(); expect(cClear).toHaveBeenCalledOnce()
    expect(f.sent).toHaveLength(1); expect(f.sent[0].type).toBe(type)
  })
  it('older failure/invalidation cannot clear a newer successful memory session', async () => {
    const f = coordinationFixture(), a = coordinator({ channel: f.channel() }), b = coordinator({ channel: f.channel() })
    const clear = vi.fn(); b.subscribe(clear)
    let old!: SessionAttempt
    await a.run(true, async (attempt) => { old = attempt })
    await b.run(true, async (attempt) => b.succeeded(attempt))
    f.channels[1].onmessage?.({ data: { ...base, tabId: a.tabId, attemptId: old.attemptId, epoch: old.epoch, type: 'refresh-failed', failure: 'unauthenticated' } } as MessageEvent)
    f.channels[1].onmessage?.({ data: { ...base, tabId: a.tabId, attemptId: old.attemptId, epoch: old.epoch, type: 'session-invalidated' } } as MessageEvent)
    expect(clear).not.toHaveBeenCalled()
  })
  it('a completed successful attempt cannot later report a contradictory invalidation', async () => {
    const f = coordinationFixture(), a = coordinator({ channel: f.channel() }), b = coordinator({ channel: f.channel() })
    const clear = vi.fn(); b.subscribe(clear)
    await a.run(true, async (attempt) => { a.succeeded(attempt); a.failed(attempt, 'unauthenticated') })
    expect(clear).not.toHaveBeenCalled()
  })
  it('malformed, old, duplicate and self messages have no lifecycle effect', () => {
    const c = fakedChannel(), b = coordinator({ channel: c }), clear = vi.fn(); b.subscribe(clear)
    c.onmessage?.({ data: { ...base, type: 'logout', password: 'secret' } } as MessageEvent)
    c.onmessage?.({ data: { ...base, type: 'logout', timestamp: now - 60_001 } } as MessageEvent)
    c.onmessage?.({ data: { ...base, type: 'logout', tabId: b.tabId } } as MessageEvent)
    expect(clear).not.toHaveBeenCalled()
    c.onmessage?.({ data: { ...base, type: 'logout' } } as MessageEvent)
    c.onmessage?.({ data: { ...base, type: 'logout' } } as MessageEvent)
    expect(clear).toHaveBeenCalledOnce()
  })
  it('closes channel/listeners and cancels queued ownership when disposed', async () => {
    const f = coordinationFixture(), a = coordinator({ locks: f.locks }), channel = fakedChannel(), b = coordinator({ locks: f.locks, channel })
    const gate = deferred<void>(), started = deferred<void>()
    const active = a.run(true, async () => { started.resolve(); await gate.promise })
    await started.promise
    const waiting = b.run(true, async () => 'forbidden-after-disposal')
    const rejection = expect(waiting).rejects.toMatchObject({ code: 'SESSION_CHANGED' })
    const listener = vi.fn(); b.subscribe(listener); b.dispose(); b.dispose()
    await rejection; expect(channel.onmessage).toBeNull(); expect(channel.close).toHaveBeenCalled()
    gate.resolve(); await active
    await expect(b.run(true, async () => undefined)).rejects.toMatchObject({ code: 'SESSION_CHANGED' })
  })
  it.each([[0, 'NETWORK_ERROR', 'network'], [503, 'INTERNAL_ERROR', 'server'], [403, 'FORBIDDEN', 'forbidden'], [401, 'UNAUTHORIZED', 'unauthenticated'], [0, 'COORDINATION_TIMEOUT', 'coordination-timeout']] as const)('classifies status %i / %s as %s', (status, code, expected) => {
    expect(classifyRefreshFailure(new ApiRequestError('Safe', status, code), true)).toBe(expected)
  })
})

function fakedChannel(): SessionChannel { return { onmessage: null, postMessage: vi.fn(), close: vi.fn() } }
