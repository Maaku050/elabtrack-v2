import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { AxiosError, type AxiosAdapter, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios'
import { ApiClient } from './api-client'
import { ApiRequestError } from './api-client'
import { QueryClient } from '@tanstack/react-query'
import { useAuthStore } from '@/stores/auth-store'
import { queryClient } from '@/app/query-client'
import type { BrowserSession } from '@/types/common'

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
afterEach(() => { vi.restoreAllMocks(); queryClient.clear() })

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
    const client = new ApiClient(adapter)
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
    const client = new ApiClient(adapter)
    // Clearing the private cache deliberately cancels its pending query.
    await expect(queryClient.fetchQuery({ queryKey: ['protected'], queryFn: () => client.get('/protected') })).rejects.toMatchObject({ message: 'CancelledError' })
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

describe('logout and session lifecycle races', () => {
  it('logout clears memory/cache immediately and remains logged out on network failure', async () => {
    useAuthStore.getState().setSession(session); queryClient.setQueryData(['private'], { owner: session.user.id })
    const adapter = vi.fn<AxiosAdapter>(async (config) => { throw new AxiosError('Network unavailable', 'ERR_NETWORK', config) })
    const client = new ApiClient(adapter); const logout = client.logout()
    expect(useAuthStore.getState().accessToken).toBeNull(); expect(useAuthStore.getState().user).toBeNull()
    expect(queryClient.getQueryData(['private'])).toBeUndefined()
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
