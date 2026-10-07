import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { AxiosError, type AxiosAdapter, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios'
import { ApiClient } from './api-client'
import { ApiRequestError } from './api-client'
import { QueryClient } from '@tanstack/react-query'
import { useAuthStore } from '@/stores/auth-store'
import { queryClient } from '@/app/query-client'
import type { BrowserSession } from '@/types/common'
import { SessionCoordinator, type SessionChannel, type SessionLocks, type TerminalEvent, SESSION_WAIT_MS } from './session-coordinator'

const synchronousLocks: SessionLocks = { request: (_name, _options, operation) => operation() }
const clients: ApiClient[] = []
const coordinators: SessionCoordinator[] = []
function coordinatedClient(adapter: AxiosAdapter, locks: SessionLocks | undefined = synchronousLocks) {
  const channel: SessionChannel = { onmessage: null, postMessage: vi.fn(), close: vi.fn() }
  const coordinator = new SessionCoordinator({ channel, locks })
  const client = new ApiClient(adapter, coordinator)
  clients.push(client); coordinators.push(coordinator)
  const notify = (type: TerminalEvent) => {
    const timestamp = performance.timeOrigin + performance.now()
    channel.onmessage?.({ data: { version: 1, type, tabId: crypto.randomUUID(), attemptId: crypto.randomUUID(), epoch: Math.ceil(timestamp * 1_000), timestamp } } as MessageEvent)
  }
  return { client, coordinator, channel, notify }
}

const session: BrowserSession = {
  access_token: 'synthetic-renewed-access', expires_at: '2026-10-06T12:00:00Z', token_type: 'Bearer',
  user: { id: '00000000-0000-0000-0000-000000000001', email: 'synthetic@example.invalid', name: 'Current Account', role: 'user', is_active: true },
}
function response(config: InternalAxiosRequestConfig, status: number, data: unknown = session): AxiosResponse {
  const res = { config, status, statusText: String(status), headers: {}, data: status < 400 ? { success: true, data, message: 'Success', meta: null, error: null } : { success: false, message: 'Authentication required.', error: { code: 'UNAUTHORIZED' } } }
  if (status >= 400) throw new AxiosError('Request failed', 'ERR_BAD_RESPONSE', config, undefined, res)
  return res
}
function deferred<T>() {
  let resolve!: (value: T) => void
  const promise = new Promise<T>((done) => { resolve = done })
  return { promise, resolve }
}
beforeEach(() => {
  useAuthStore.getState().clear()
  useAuthStore.setState({ status: 'idle', generation: 0 })
  localStorage.clear(); sessionStorage.clear(); queryClient.clear()
})
afterEach(() => {
  for (const client of clients.splice(0)) client.dispose()
  for (const coordinator of coordinators.splice(0)) coordinator.dispose()
  vi.restoreAllMocks(); vi.unstubAllGlobals(); vi.useRealTimers(); queryClient.clear()
})

describe('memory and cookie browser sessions', () => {
  it('login keeps access only in memory and accepts current safe account state', async () => {
    localStorage.setItem('elabtrack_v2.access_token', 'legacy-access')
    localStorage.setItem('elabtrack_v2.refresh_token', 'legacy-refresh')
    sessionStorage.setItem('elabtrack_v2.access_token', 'legacy-access')
    const get = vi.spyOn(Storage.prototype, 'getItem')
    const set = vi.spyOn(Storage.prototype, 'setItem')
    const adapter = vi.fn<AxiosAdapter>(async (config) => response(config, 200, { ...session, user: { ...session.user, password: 'unexpected-field' } }))
    const client = new ApiClient(adapter)
    const result = await client.authenticate('/auth/login', { email: 'synthetic@example.invalid', password: 'synthetic-password' })
    expect(get).not.toHaveBeenCalled()
    expect(set.mock.calls.filter(([key]) => /(?:access|refresh)_token/.test(key))).toEqual([])
    expect(result).toEqual(session)
    expect(useAuthStore.getState().accessToken).toBe(session.access_token)
    expect(useAuthStore.getState().user?.role).toBe('user')
    expect(useAuthStore.getState()).not.toHaveProperty('refreshToken')
    expect(result).not.toHaveProperty('refresh_token')
    expect(useAuthStore.getState().user).not.toHaveProperty('password')
    expect(adapter.mock.calls[0][0].headers.get('Authorization')).toBeUndefined()
    expect(adapter.mock.calls[0][0].withCredentials).toBe(true)
    get.mockRestore()
    expect(localStorage.getItem('elabtrack_v2.access_token')).toBeNull()
    expect(localStorage.getItem('elabtrack_v2.refresh_token')).toBeNull()
    expect(sessionStorage.getItem('elabtrack_v2.access_token')).toBeNull()
  })
  it('purges legacy keys without reading values and bootstraps exactly once, including repeated mounts', async () => {
    localStorage.setItem('elabtrack_v2.access_token', 'poison-access')
    localStorage.setItem('elabtrack_v2.refresh_token', 'poison-refresh')
    localStorage.setItem('elabtrack_v2.theme', 'dark')
    const get = vi.spyOn(Storage.prototype, 'getItem')
    const gate = deferred<void>()
    const adapter = vi.fn<AxiosAdapter>(async (config) => { await gate.promise; return response(config, 200) })
    const client = new ApiClient(adapter)
    const first = client.bootstrap(); const second = client.bootstrap()
    expect(first).toBe(second)
    expect(useAuthStore.getState().status).toBe('bootstrapping')
    expect(useAuthStore.getState().accessToken).toBeNull()
    gate.resolve(); await Promise.all([first, second, client.bootstrap()])
    expect(adapter).toHaveBeenCalledTimes(1)
    const config = adapter.mock.calls[0][0]
    expect(config.url).toBe('/auth/refresh'); expect(config.data).toBeUndefined()
    expect(config.headers.get('Authorization')).toBeUndefined()
    expect(get).not.toHaveBeenCalled()
    expect(useAuthStore.getState().status).toBe('authenticated')
    get.mockRestore()
    expect(localStorage.getItem('elabtrack_v2.theme')).toBe('dark')
    // A new document/store loses memory; the following attempt restores from cookie again.
    useAuthStore.getState().clear(); useAuthStore.setState({ status: 'idle' })
    expect(useAuthStore.getState().accessToken).toBeNull()
    await new ApiClient(adapter).bootstrap()
    expect(adapter).toHaveBeenCalledTimes(2)
  })
  it.each([401, 403])('bootstrap auth failure %s becomes unauthenticated without a retry', async (status) => {
    const adapter = vi.fn<AxiosAdapter>(async (config) => response(config, status))
    const client = new ApiClient(adapter)
    await Promise.all([client.bootstrap(), client.bootstrap()]); await client.bootstrap()
    expect(adapter).toHaveBeenCalledTimes(1)
    expect(useAuthStore.getState().status).toBe('unauthenticated')
    expect(useAuthStore.getState().accessToken).toBeNull()
    expect(useAuthStore.getState().user).toBeNull()
  })
  it.each(['network', 'server'])('bootstrap %s failure is recoverable only by deliberate retry', async (kind) => {
    let failing = true
    const adapter = vi.fn<AxiosAdapter>(async (config) => {
      if (failing && kind === 'network') throw new AxiosError('Network failed', 'ERR_NETWORK', config)
      return response(config, failing ? 503 : 200)
    })
    const client = new ApiClient(adapter)
    await client.bootstrap(); await client.bootstrap()
    expect(adapter).toHaveBeenCalledTimes(1)
    expect(useAuthStore.getState().status).toBe('error')
    expect(useAuthStore.getState().accessToken).toBeNull()
    failing = false; await client.retryBootstrap()
    expect(adapter).toHaveBeenCalledTimes(2)
    expect(useAuthStore.getState().status).toBe('authenticated')
  })
})

describe('single-flight and bounded authentication retries', () => {
  it('N simultaneous 401s issue one refresh and retry with the replacement memory token', async () => {
    useAuthStore.getState().setSession({ ...session, access_token: 'old-access' })
    const gate = deferred<void>(); const started = deferred<void>()
    const adapter = vi.fn<AxiosAdapter>(async (config) => {
      if (config.url === '/auth/refresh') { started.resolve(); await gate.promise; return response(config, 200) }
      return response(config, config.headers.get('Authorization') === 'Bearer old-access' ? 401 : 200, { ok: true })
    })
    const { client } = coordinatedClient(adapter)
    const requests = Array.from({ length: 8 }, (_, i) => client.get(`/protected/${i}`))
    await started.promise; gate.resolve()
    await expect(Promise.all(requests)).resolves.toEqual(Array.from({ length: 8 }, () => ({ ok: true })))
    const calls = adapter.mock.calls.map(([config]) => config)
    expect(calls.filter((config) => config.url === '/auth/refresh')).toHaveLength(1)
    expect(calls.filter((config) => config.headers.get('Authorization') === `Bearer ${session.access_token}`)).toHaveLength(8)
    expect(calls).toHaveLength(17)
  })
  it('a late 401 reuses an already renewed token instead of rotating the cookie again', async () => {
    useAuthStore.getState().setSession({ ...session, access_token: 'old-access' })
    const late = deferred<void>(); const sent = deferred<void>()
    const adapter = vi.fn<AxiosAdapter>(async (config) => {
      if (config.url === '/auth/refresh') return response(config, 200)
      const old = config.headers.get('Authorization') === 'Bearer old-access'
      if (config.url === '/late' && old) { sent.resolve(); await late.promise }
      return response(config, old ? 401 : 200, { ok: true })
    })
    const client = new ApiClient(adapter); const delayed = client.get('/late'); await sent.promise
    await client.get('/first'); late.resolve(); await delayed
    expect(adapter.mock.calls.filter(([config]) => config.url === '/auth/refresh')).toHaveLength(1)
  })
  it.each([401, 503])('failed refresh %s drains waiters and clears memory without another refresh', async (status) => {
    useAuthStore.getState().setSession({ ...session, access_token: 'old-access' })
    const gate = deferred<void>(); const started = deferred<void>()
    const adapter = vi.fn<AxiosAdapter>(async (config) => {
      if (config.url === '/auth/refresh') { started.resolve(); await gate.promise; return response(config, status) }
      return response(config, 401)
    })
    const client = new ApiClient(adapter)
    const requests = Promise.allSettled(Array.from({ length: 6 }, (_, i) => client.get(`/protected/${i}`)))
    await started.promise; gate.resolve(); const results = await requests
    expect(results.every((result) => result.status === 'rejected')).toBe(true)
    expect(adapter.mock.calls.filter(([config]) => config.url === '/auth/refresh')).toHaveLength(1)
    expect(adapter).toHaveBeenCalledTimes(7)
    expect(useAuthStore.getState().accessToken).toBeNull(); expect(useAuthStore.getState().user).toBeNull()
    expect(useAuthStore.getState().status).toBe(status === 401 ? 'unauthenticated' : 'error')
  })
  it('401 → refresh → retried 401 stops after one retry, including Query retry policy', async () => {
    useAuthStore.getState().setSession({ ...session, access_token: 'old-access' })
    const adapter = vi.fn<AxiosAdapter>(async (config) => response(config, config.url === '/auth/refresh' ? 200 : 401))
    const { client } = coordinatedClient(adapter)
    // Clearing the private cache deliberately cancels its pending query.
    await expect(queryClient.fetchQuery({ queryKey: ['auth', 'protected'], queryFn: () => client.get('/protected') })).rejects.toMatchObject({ message: 'CancelledError' })
    expect(adapter).toHaveBeenCalledTimes(3)
    expect(adapter.mock.calls.filter(([config]) => config.url === '/auth/refresh')).toHaveLength(1)
    expect(useAuthStore.getState().status).toBe('unauthenticated')
  })
  it.each([401, 403])('Query never restarts auth failure %s', async (status) => {
    const client = new QueryClient({ defaultOptions: queryClient.getDefaultOptions() })
    const query = vi.fn(async () => { throw new ApiRequestError('Authentication failed', status) })
    await expect(client.fetchQuery({ queryKey: ['denied'], queryFn: query })).rejects.toMatchObject({ status })
    expect(query).toHaveBeenCalledTimes(1)
    client.clear()
  })
  it.each(['/auth/login', '/auth/register', '/auth/refresh', '/auth/logout'])('%s never recursively refreshes', async (path) => {
    useAuthStore.getState().setSession(session)
    const adapter = vi.fn<AxiosAdapter>(async (config) => response(config, 401))
    await expect(new ApiClient(adapter).post(path)).rejects.toMatchObject({ status: 401 })
    expect(adapter).toHaveBeenCalledTimes(1)
  })
  it.each([403, 422, 500])('does not refresh for non-401 failure %s', async (status) => {
    useAuthStore.getState().setSession(session)
    const adapter = vi.fn<AxiosAdapter>(async (config) => response(config, status))
    await expect(new ApiClient(adapter).get('/protected')).rejects.toMatchObject({ status })
    expect(adapter).toHaveBeenCalledTimes(1)
  })
  it('respects attachAuth=false and does not restore a logged-out public request', async () => {
    useAuthStore.getState().setSession(session)
    const adapter = vi.fn<AxiosAdapter>(async (config) => response(config, 401))
    await expect(new ApiClient(adapter).get('/public', { attachAuth: false })).rejects.toMatchObject({ status: 401 })
    expect(adapter).toHaveBeenCalledTimes(1); expect(adapter.mock.calls[0][0].headers.get('Authorization')).toBeUndefined()
  })
})

describe('coordinator integration with memory, Query and HTTP fencing', () => {
  it('simultaneous peer logout cannot cancel every backend revocation', async () => {
    useAuthStore.getState().setSession(session)
    const channels: SessionChannel[] = []
    const adapter = vi.fn<AxiosAdapter>(async (config) => response(config, 204))
    const tabs = Array.from({ length: 3 }, () => {
      const channel: SessionChannel = { onmessage: null, close: vi.fn(), postMessage: message => { for (const peer of channels) if (peer !== channel) peer.onmessage?.({ data: message } as MessageEvent) } }
      channels.push(channel)
      const coordinator = new SessionCoordinator({ channel, locks: synchronousLocks })
      const client = new ApiClient(adapter, coordinator); clients.push(client); coordinators.push(coordinator)
      return client
    })
    await Promise.all(tabs.map(client => client.logout()))
    expect(adapter).toHaveBeenCalledTimes(3)
    expect(adapter.mock.calls.every(([config]) => config.url === '/auth/logout')).toBe(true)
    expect(useAuthStore.getState().status).toBe('unauthenticated')
  })
  it('login announces both intent and successful completion without conveying credentials', async () => {
    const adapter = vi.fn<AxiosAdapter>(async (config) => response(config, 200))
    const { client, channel } = coordinatedClient(adapter)
    await client.authenticate('/auth/login', { password: 'synthetic-password' })
    const hints = vi.mocked(channel.postMessage).mock.calls.map(([m]) => m)
    expect(hints.map(m => m.type)).toEqual(['auth-state-changed', 'auth-state-changed'])
    expect(hints[1].epoch).toBeGreaterThan(hints[0].epoch)
    expect(JSON.stringify(hints)).not.toMatch(/synthetic-password|synthetic-renewed-access/)
  })
  it('a late denial for an older retried token cannot clear a newer local recovery', async () => {
    useAuthStore.getState().setSession({ ...session, access_token: 'old-access' })
    const started = deferred<void>(), gate = deferred<void>()
    let rotations = 0
    const adapter = vi.fn<AxiosAdapter>(async (config) => {
      if (config.url === '/auth/refresh') return response(config, 200, { ...session, access_token: `renewed-${++rotations}` })
      if (config.headers.get('Authorization') === 'Bearer old-access') return response(config, 401)
      started.resolve(); await gate.promise; return response(config, 401)
    })
    const { client, channel } = coordinatedClient(adapter)
    const pending = client.get('/protected'); const denial = expect(pending).rejects.toMatchObject({ status: 401 })
    await started.promise; await client.refreshSession(); gate.resolve(); await denial
    expect(useAuthStore.getState().accessToken).toBe('renewed-2')
    expect(vi.mocked(channel.postMessage).mock.calls.map(([m]) => m.type)).not.toContain('session-invalidated')
    expect(adapter.mock.calls.filter(([c]) => c.url === '/protected')).toHaveLength(2)
  })
  it('a late current-account 403 for a superseded access token cannot invalidate a newer recovery', async () => {
    useAuthStore.getState().setSession(session)
    const started = deferred<void>(), gate = deferred<void>()
    const adapter = vi.fn<AxiosAdapter>(async (config) => {
      if (config.url === '/auth/refresh') return response(config, 200, { ...session, access_token: 'newest' })
      started.resolve(); await gate.promise; return response(config, 403)
    })
    const { client, channel } = coordinatedClient(adapter)
    const pending = client.get('/auth/me'); const denial = expect(pending).rejects.toMatchObject({ status: 403 })
    await started.promise; await client.refreshSession(); gate.resolve(); await denial
    expect(useAuthStore.getState().accessToken).toBe('newest')
    expect(vi.mocked(channel.postMessage).mock.calls.map(([m]) => m.type)).not.toContain('session-invalidated')
  })
  it.each(['logout', 'session-invalidated', 'auth-state-changed'] as const)('peer %s clears private memory/cache, cancels active work and preserves public data', async (event) => {
    useAuthStore.getState().setSession(session)
    queryClient.setQueryData(['auth', 'me'], session.user)
    queryClient.setQueryData(['foundation', 'health'], { status: 'ok' })
    const adapter = vi.fn<AxiosAdapter>(async (config) => response(config, 200))
    const { notify } = coordinatedClient(adapter)
    let signal!: AbortSignal
    const started = deferred<void>()
    const pending = queryClient.fetchQuery({ queryKey: ['future-private'], meta: { authenticated: true }, queryFn: ({ signal: current }) => { signal = current; started.resolve(); return new Promise(() => {}) } })
    const rejection = expect(pending).rejects.toMatchObject({ message: 'CancelledError' })
    await started.promise
    notify(event); await rejection
    expect(signal.aborted).toBe(true)
    expect(useAuthStore.getState()).toMatchObject({ accessToken: null, user: null, status: 'unauthenticated' })
    expect(queryClient.getQueryData(['auth', 'me'])).toBeUndefined()
    expect(queryClient.getQueryCache().find({ queryKey: ['future-private'] })).toBeUndefined()
    expect(queryClient.getQueryData(['foundation', 'health'])).toEqual({ status: 'ok' })
    expect(adapter).not.toHaveBeenCalled()
  })
  it('peer logout fences a pending successful refresh without another recovery', async () => {
    useAuthStore.getState().setSession(session)
    const gate = deferred<void>(), started = deferred<void>()
    const adapter = vi.fn<AxiosAdapter>(async (config) => { started.resolve(); await gate.promise; return response(config, 200) })
    const { client, notify } = coordinatedClient(adapter)
    const refresh = client.refreshSession()
    const rejection = expect(refresh).rejects.toMatchObject({ code: 'SESSION_CHANGED' })
    await started.promise; notify('logout'); gate.resolve(); await rejection
    expect(useAuthStore.getState().accessToken).toBeNull()
    expect(adapter).toHaveBeenCalledOnce()
  })
  it('authoritative current-account 403 invalidates; ordinary resource 403 does not', async () => {
    useAuthStore.getState().setSession(session)
    const adapter = vi.fn<AxiosAdapter>(async (config) => response(config, 403))
    const { client, channel } = coordinatedClient(adapter)
    await expect(client.get('/admin-resource')).rejects.toMatchObject({ status: 403 })
    expect(useAuthStore.getState().status).toBe('authenticated')
    expect(channel.postMessage).not.toHaveBeenCalled()
    await expect(client.get('/auth/me')).rejects.toMatchObject({ status: 403 })
    expect(useAuthStore.getState().status).toBe('unauthenticated')
    expect(channel.postMessage).toHaveBeenCalledWith(expect.objectContaining({ type: 'session-invalidated' }))
    expect(adapter).toHaveBeenCalledTimes(2)
  })
  it('a lock timeout settles bootstrap as recoverable error and deliberate retry can restore', async () => {
    vi.useFakeTimers()
    let suspended = true
    const locks: SessionLocks = { request: (_name, { signal }, operation) => suspended ? new Promise((_resolve, reject) => { signal.addEventListener('abort', () => reject(new DOMException('Aborted', 'AbortError')), { once: true }) }) : operation() }
    const adapter = vi.fn<AxiosAdapter>(async (config) => response(config, 200))
    const { client } = coordinatedClient(adapter, locks)
    const boot = client.bootstrap()
    await vi.advanceTimersByTimeAsync(SESSION_WAIT_MS); await boot
    expect(useAuthStore.getState().status).toBe('error'); expect(adapter).not.toHaveBeenCalled()
    suspended = false; await client.retryBootstrap()
    expect(useAuthStore.getState().status).toBe('authenticated'); expect(adapter).toHaveBeenCalledOnce()
  })
  it('coordinated login/refresh/logout never write auth or coordination values to any browser store', async () => {
    const setLocal = vi.spyOn(Storage.prototype, 'setItem')
    const indexed = { open: vi.fn(), deleteDatabase: vi.fn() }; vi.stubGlobal('indexedDB', indexed)
    const adapter = vi.fn<AxiosAdapter>(async (config) => response(config, config.url === '/auth/logout' ? 204 : 200))
    const { client, channel } = coordinatedClient(adapter)
    await client.authenticate('/auth/login', { password: 'synthetic-input' }); await client.refreshSession(); await client.logout()
    expect(setLocal.mock.calls.every(([key, value]) => key === '__storage_test__' && value === '1')).toBe(true)
    expect(localStorage.length).toBe(0); expect(sessionStorage.length).toBe(0)
    expect(indexed.open).not.toHaveBeenCalled(); expect(indexed.deleteDatabase).not.toHaveBeenCalled()
    expect(JSON.stringify(vi.mocked(channel.postMessage).mock.calls)).not.toMatch(/synthetic-input|synthetic-renewed-access|refresh_token|Authorization|Cookie|password/)
  })
  it('disposing the client prevents a pending reply from restoring auth and removes its peer subscription', async () => {
    const gate = deferred<void>(), started = deferred<void>()
    const adapter = vi.fn<AxiosAdapter>(async (config) => { started.resolve(); await gate.promise; return response(config, 200) })
    const { client, notify } = coordinatedClient(adapter)
    const refresh = client.refreshSession(); const rejection = expect(refresh).rejects.toMatchObject({ code: 'SESSION_CHANGED' })
    await started.promise; client.dispose(); gate.resolve(); await rejection
    useAuthStore.getState().setSession(session); notify('logout')
    expect(useAuthStore.getState().status).toBe('authenticated')
  })
})

describe('logout and session lifecycle races', () => {
  it('logout clears memory/cache immediately and remains logged out on network failure', async () => {
    useAuthStore.getState().setSession(session); queryClient.setQueryData(['auth', 'private'], { owner: session.user.id })
    const adapter = vi.fn<AxiosAdapter>(async (config) => { throw new AxiosError('Network unavailable', 'ERR_NETWORK', config) })
    const client = new ApiClient(adapter); const logout = client.logout()
    expect(useAuthStore.getState().accessToken).toBeNull(); expect(useAuthStore.getState().user).toBeNull()
    expect(queryClient.getQueryData(['auth', 'private'])).toBeUndefined()
    await expect(logout).rejects.toMatchObject({ status: 0 })
    expect(adapter).toHaveBeenCalledTimes(1); expect(adapter.mock.calls[0][0].url).toBe('/auth/logout')
    expect(adapter.mock.calls[0][0].data).toBeUndefined()
    expect(useAuthStore.getState().status).toBe('unauthenticated')
    expect(localStorage.length).toBe(0)
  })
  it('pending refresh cannot resurrect logout, and logout waits to revoke the replacement cookie', async () => {
    useAuthStore.getState().setSession(session)
    const gate = deferred<void>(); const started = deferred<void>(); const order: string[] = []
    const adapter = vi.fn<AxiosAdapter>(async (config) => {
      order.push(config.url ?? '')
      if (config.url === '/auth/refresh') { started.resolve(); await gate.promise; return response(config, 200) }
      return response(config, 204)
    })
    const client = new ApiClient(adapter)
    const refresh = client.refreshSession().catch(() => undefined); await started.promise
    const logout = client.logout(); expect(useAuthStore.getState().isAuthenticated).toBe(false)
    expect(order).toEqual(['/auth/refresh'])
    gate.resolve(); await Promise.all([refresh, logout])
    expect(order).toEqual(['/auth/refresh', '/auth/logout'])
    expect(useAuthStore.getState().status).toBe('unauthenticated')
    expect(useAuthStore.getState().accessToken).toBeNull()
  })
  it('pending login cannot resurrect logout', async () => {
    const gate = deferred<void>(); const started = deferred<void>()
    const adapter = vi.fn<AxiosAdapter>(async (config) => {
      if (config.url === '/auth/login') { started.resolve(); await gate.promise; return response(config, 200) }
      return response(config, 204)
    })
    const client = new ApiClient(adapter)
    const login = client.authenticate('/auth/login', {}).catch(() => undefined); await started.promise
    const logout = client.logout(); gate.resolve(); await Promise.all([login, logout])
    expect(useAuthStore.getState().status).toBe('unauthenticated')
    expect(adapter.mock.calls.map(([config]) => config.url)).toEqual(['/auth/login', '/auth/logout'])
  })
  it('a late protected 401 cannot trigger refresh after logout', async () => {
    useAuthStore.getState().setSession(session)
    const gate = deferred<void>(); const started = deferred<void>()
    const adapter = vi.fn<AxiosAdapter>(async (config) => {
      if (config.url === '/protected') { started.resolve(); await gate.promise; return response(config, 401) }
      return response(config, 204)
    })
    const client = new ApiClient(adapter); const request = client.get('/protected').catch(() => undefined)
    await started.promise; await client.logout(); gate.resolve(); await request
    expect(adapter.mock.calls.some(([config]) => config.url === '/auth/refresh')).toBe(false)
    expect(useAuthStore.getState().status).toBe('unauthenticated')
  })
})
