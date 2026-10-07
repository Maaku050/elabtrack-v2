import { ApiRequestError } from '@/lib/api-error'

export const SESSION_CHANNEL = 'elabtrack_v2.session.v1'
export const SESSION_LOCK = 'elabtrack_v2.session-cookie.v1'
// Axios bounds the active HTTP request to 15s. Waiting for ownership is also
// bounded; never steal a live lock while its response can still change cookies.
export const SESSION_WAIT_MS = 20_000
export const SESSION_EVENT_AGE_MS = 60_000

export type RefreshFailure = 'network' | 'server' | 'unauthenticated' | 'ambiguous-denial' | 'forbidden' | 'coordination-timeout'
type EventType = 'refresh-started' | 'refresh-succeeded' | 'logout' | 'session-invalidated' | 'auth-state-changed'
interface EventBase { version: 1; tabId: string; attemptId: string; epoch: number; timestamp: number }
export type SessionEvent = { [K in EventType]: EventBase & { type: K } }[EventType] | (EventBase & { type: 'refresh-failed'; failure: RefreshFailure })
export type TerminalEvent = Extract<EventType, 'logout' | 'session-invalidated' | 'auth-state-changed'>
export interface SessionAttempt { attemptId: string; epoch: number; exclusive: boolean }
export interface SessionChannel {
  onmessage: ((event: MessageEvent<unknown>) => void) | null
  postMessage(message: SessionEvent): void
  close(): void
}
export interface SessionLocks {
  request<T>(name: string, options: { mode: 'exclusive'; signal: AbortSignal }, callback: () => Promise<T>): Promise<T>
}
interface CoordinatorOptions {
  channel?: SessionChannel
  locks?: SessionLocks
  tabId?: string
  now?: () => number
  waitMs?: number
}

const uuid = (value: unknown): value is string => typeof value === 'string' && /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/.test(value)
const types = new Set(['refresh-started', 'refresh-succeeded', 'refresh-failed', 'logout', 'session-invalidated', 'auth-state-changed'])
const failures = new Set<RefreshFailure>(['network', 'server', 'unauthenticated', 'ambiguous-denial', 'forbidden', 'coordination-timeout'])

/** Strict allowlist: even a recognized event carrying extra credential fields
 * is rejected. Messages are hints for local lifecycle, never authorization. */
export function parseSessionEvent(value: unknown, now: number): SessionEvent | null {
  if (!value || typeof value !== 'object' || Array.isArray(value)) return null
  const m = value as Record<string, unknown>
  const keys = ['version', 'type', 'tabId', 'attemptId', 'epoch', 'timestamp', ...(m.type === 'refresh-failed' ? ['failure'] : [])]
  if (Object.keys(m).length !== keys.length || Object.keys(m).some((k) => !keys.includes(k))) return null
  if (m.version !== 1 || typeof m.type !== 'string' || !types.has(m.type) || !uuid(m.tabId) || !uuid(m.attemptId)) return null
  if (!Number.isSafeInteger(m.epoch) || (m.epoch as number) <= 0 || typeof m.timestamp !== 'number' || !Number.isFinite(m.timestamp)) return null
  if (m.timestamp < now - SESSION_EVENT_AGE_MS || m.timestamp > now + 5_000 || Math.abs((m.epoch as number) / 1_000 - m.timestamp) > SESSION_EVENT_AGE_MS) return null
  if (m.type === 'refresh-failed' && !failures.has(m.failure as RefreshFailure)) return null
  return m as unknown as SessionEvent
}

/** Native origin-scoped Web Locks provide atomic ownership. BroadcastChannel
 * only carries bounded, non-secret outcomes; no token or account crosses it.
 * A missed completion cannot block a waiter: the browser manages lock release. */
export class SessionCoordinator {
  readonly tabId: string
  private options: CoordinatorOptions
  private epoch = 0
  private latest: { epoch: number; attemptId: string; rank: number; type: SessionEvent['type'] } | null = null
  private listeners = new Set<(event: SessionEvent & { type: TerminalEvent }) => void>()
  private waiting = new Set<AbortController>()
  private disposed = false

  constructor(options: CoordinatorOptions = {}) {
    this.options = options
    this.tabId = options.tabId ?? crypto.randomUUID()
    if (options.channel) options.channel.onmessage = (event) => this.receive(event.data)
  }
  get exclusive(): boolean { return !!this.options.locks }
  private now(): number { return this.options.now?.() ?? performance.timeOrigin + performance.now() }
  private attempt(): SessionAttempt {
    this.epoch = Math.max(this.epoch + 1, Math.floor(this.now() * 1_000))
    return { attemptId: crypto.randomUUID(), epoch: this.epoch, exclusive: this.exclusive }
  }
  subscribe(listener: (event: SessionEvent & { type: TerminalEvent }) => void): () => void {
    this.listeners.add(listener)
    return () => { this.listeners.delete(listener) }
  }
  private advance(message: SessionEvent): boolean {
    const rank = message.type === 'refresh-started' ? 0 : message.type === 'refresh-failed' || message.type === 'refresh-succeeded' ? 1 : 2
    const prev = this.latest
    this.epoch = Math.max(this.epoch, message.epoch)
    if (prev?.attemptId === message.attemptId && prev.type === 'refresh-succeeded' && (message.type === 'refresh-failed' || message.type === 'session-invalidated')) return false
    if (prev && (message.epoch < prev.epoch || message.epoch === prev.epoch && (message.attemptId < prev.attemptId || message.attemptId === prev.attemptId && rank <= prev.rank))) return false
    this.latest = { epoch: message.epoch, attemptId: message.attemptId, rank, type: message.type }
    return true
  }
  private receive(value: unknown): void {
    const event = parseSessionEvent(value, this.now())
    if (this.disposed || !event || event.tabId === this.tabId || !this.advance(event)) return
    if (event.type === 'logout' || event.type === 'session-invalidated' || event.type === 'auth-state-changed') {
      for (const listener of this.listeners) listener(event)
    }
  }
  private send(event: SessionEvent): void {
    if (this.disposed || !this.advance(event)) return
    try { this.options.channel?.postMessage(event) } catch { /* Local/server safeguards remain if the channel fails. */ }
  }
  private base(attempt: SessionAttempt): EventBase {
    return { version: 1, tabId: this.tabId, attemptId: attempt.attemptId, epoch: attempt.epoch, timestamp: this.now() }
  }
  announce(type: TerminalEvent): void {
    this.send({ ...this.base(this.attempt()), type })
  }
  capture(): SessionAttempt { return this.attempt() }
  invalidate(attempt: SessionAttempt): void {
    this.send({ ...this.base(attempt), type: 'session-invalidated' })
  }
  succeeded(attempt: SessionAttempt): void {
    this.send({ ...this.base(attempt), type: 'refresh-succeeded' })
  }
  failed(attempt: SessionAttempt, failure: RefreshFailure): void {
    this.send({ ...this.base(attempt), type: 'refresh-failed', failure })
    if (failure === 'unauthenticated' && attempt.exclusive) {
      this.send({ ...this.base(attempt), type: 'session-invalidated' })
    }
  }
  /** All cookie-changing endpoints share this lock. Refresh callbacks must
   * retain their local generation check before sending and accepting replies. */
  async run<T>(refresh: boolean, operation: (attempt: SessionAttempt) => Promise<T>): Promise<T> {
    if (this.disposed) throw new ApiRequestError('Session coordination ended.', 0, 'SESSION_CHANGED')
    const execute = () => {
      if (this.disposed) throw new ApiRequestError('Session coordination ended.', 0, 'SESSION_CHANGED')
      const attempt = this.attempt()
      if (refresh) this.send({ ...this.base(attempt), type: 'refresh-started' })
      return operation(attempt)
    }
    if (!this.options.locks) return execute()
    const controller = new AbortController()
    this.waiting.add(controller)
    let timeout = false
    const timer = setTimeout(() => { timeout = true; controller.abort() }, this.options.waitMs ?? SESSION_WAIT_MS)
    try {
      return await this.options.locks.request(SESSION_LOCK, { mode: 'exclusive', signal: controller.signal }, () => {
        clearTimeout(timer)
        this.waiting.delete(controller)
        return execute()
      })
    } catch (error) {
      if (timeout) throw new ApiRequestError('Session recovery is busy. Please try again.', 0, 'COORDINATION_TIMEOUT')
      if (controller.signal.aborted) throw new ApiRequestError('Session coordination ended.', 0, 'SESSION_CHANGED')
      throw error
    } finally { clearTimeout(timer); this.waiting.delete(controller) }
  }
  dispose(): void {
    if (this.disposed) return
    this.disposed = true
    for (const controller of this.waiting) controller.abort()
    this.waiting.clear()
    this.listeners.clear()
    if (this.options.channel) { this.options.channel.onmessage = null; this.options.channel.close() }
  }
}

export function classifyRefreshFailure(error: ApiRequestError, exclusive: boolean): RefreshFailure {
  if (error.code === 'COORDINATION_TIMEOUT') return 'coordination-timeout'
  if (error.status === 401) return exclusive ? 'unauthenticated' : 'ambiguous-denial'
  if (error.status === 403) return 'forbidden'
  return error.status === 0 ? 'network' : 'server'
}

function browserCoordinator(): SessionCoordinator {
  let channel: SessionChannel | undefined
  try { if (typeof BroadcastChannel !== 'undefined') channel = new BroadcastChannel(SESSION_CHANNEL) } catch { /* Unsupported or denied browser primitive. */ }
  const locks = typeof navigator !== 'undefined' && navigator.locks ? navigator.locks : undefined
  return new SessionCoordinator({ channel, locks })
}

// One application coordinator per document, independent of React remounts.
export const sessionCoordinator = browserCoordinator()
if (import.meta.hot) import.meta.hot.dispose(() => sessionCoordinator.dispose())
