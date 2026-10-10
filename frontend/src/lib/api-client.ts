import axios, {
  AxiosError,
  type AxiosAdapter,
  type AxiosInstance,
  type AxiosRequestConfig,
  type InternalAxiosRequestConfig,
} from 'axios'
import type { ApiResponse } from '@/types/api'
import type { BrowserSession, AuthUser } from '@/types/common'
import { useAuthStore } from '@/stores/auth-store'
import { useSessionActionStore } from '@/stores/session-action-store'
import { clearLegacyAuthStorage } from '@/lib/storage'
import { ApiRequestError, normalizeApiError, normalizeApiResponseError } from '@/lib/api-error'
import { SessionCoordinator, sessionCoordinator, classifyRefreshFailure, type SessionAttempt } from '@/lib/session-coordinator'
export { ApiRequestError } from '@/lib/api-error'

const baseURL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api/v1'

export interface RequestOptions extends Omit<AxiosRequestConfig, 'auth'> {
  attachAuth?: boolean
  retryAuth?: boolean
}
interface SessionRequest extends InternalAxiosRequestConfig {
  attachAuth?: boolean
  retryAuth?: boolean
  sessionRetry?: boolean
  sessionGeneration?: number
  sentAccessToken?: string | null
  sessionAttempt?: SessionAttempt
}

// One transport and one refresh flight per application instance. Session
// mutations are serialized so logout revokes the cookie from a pending reply.
export class ApiClient {
  private axios: AxiosInstance
  private refreshPromise: Promise<BrowserSession> | null = null
  private bootstrapPromise: Promise<void> | null = null
  private sessionTask: Promise<unknown> = Promise.resolve()
  private coordinator: SessionCoordinator
  private unsubscribe: () => void
  private disposed = false
  private authChange = 0

  constructor(adapter?: AxiosAdapter, coordinator = new SessionCoordinator()) {
    this.coordinator = coordinator
    this.unsubscribe = coordinator.subscribe((event) => {
      if (event.type === 'auth-state-changed') this.authChange++
      // Notice metadata does not participate in session authority or fencing.
      const live = useAuthStore.getState().isAuthenticated
      useSessionActionStore.getState().setSessionNotice(live && event.type !== 'logout' ? 'invalidated' : null)
      useAuthStore.getState().clear()
    })
    this.axios = axios.create({ baseURL, timeout: 15_000, headers: { 'Content-Type': 'application/json' }, withCredentials: true, adapter })
    this.axios.interceptors.request.use(this.attachToken)
    this.axios.interceptors.response.use((response) => response, (error: AxiosError<ApiResponse>) => this.onError(error))
  }
  dispose(): void { this.disposed = true; this.unsubscribe() }
  get<T>(url: string, opts: RequestOptions = {}): Promise<T> { return this.request<T>('GET', url, undefined, opts) }
  post<T>(url: string, body?: unknown, opts: RequestOptions = {}): Promise<T> { return this.request<T>('POST', url, body, opts) }
  patch<T>(url: string, body?: unknown, opts: RequestOptions = {}): Promise<T> { return this.request<T>('PATCH', url, body, opts) }
  put<T>(url: string, body?: unknown, opts: RequestOptions = {}): Promise<T> { return this.request<T>('PUT', url, body, opts) }
  delete<T>(url: string, opts: RequestOptions = {}): Promise<T> { return this.request<T>('DELETE', url, opts.data, opts) }

  async download(url: string, signal?: AbortSignal): Promise<Blob> {
    try { const response = await this.axios.get<Blob>(url, { responseType: 'blob', signal }); return response.data }
    catch (error) { throw this.normalize(error) }
  }

  private async request<T>(method: string, url: string, body: unknown, opts: RequestOptions): Promise<T> {
    try {
      const response = await this.axios.request<ApiResponse<T>>({ ...opts, method, url, data: body })
      if (response.status === 204) return undefined as T
      if (response.data?.success !== true) throw normalizeApiResponseError(response.status, response.data, response.headers['x-request-id'])
      return response.data.data
    } catch (error) { throw this.normalize(error) }
  }
  private attachToken = (original: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
    if (this.disposed) throw new ApiRequestError('Session ended', 0, 'SESSION_CHANGED')
    const config = original as SessionRequest
    const state = useAuthStore.getState()
    if (config.sessionRetry && config.sessionGeneration !== state.generation) {
      throw new ApiRequestError('Session ended', 401)
    }
    config.sessionGeneration = state.generation
    config.sessionAttempt = this.coordinator.capture()
    config.sentAccessToken = config.attachAuth === false ? null : state.accessToken
    if (config.sentAccessToken) config.headers.set('Authorization', `Bearer ${config.sentAccessToken}`)
    else config.headers.delete('Authorization')
    return config
  }
  private isSessionEndpoint(config: SessionRequest): boolean {
    const path = config.url?.split('?')[0].replace(/\/$/, '')
    return ['/auth/login', '/auth/register', '/auth/refresh', '/auth/logout'].includes(path ?? '')
  }
  private async onError(error: AxiosError<ApiResponse>): Promise<unknown> {
    if (this.disposed) throw error
    const config = error.config as SessionRequest | undefined
    // These retained current-account routes have no resource-specific 403:
    // the authoritative resolver denied current active/recognized account state.
    if (config && error.response?.status === 403 && config.sentAccessToken && config.sentAccessToken === useAuthStore.getState().accessToken && config.sessionGeneration === useAuthStore.getState().generation && ['/auth/me', '/users/me'].includes(config.url?.split('?')[0] ?? '')) {
      useSessionActionStore.getState().setSessionNotice('invalidated')
      useAuthStore.getState().clear()
      if (config.sessionAttempt) this.coordinator.invalidate(config.sessionAttempt)
    }
    if (!config || error.response?.status !== 401 || config.attachAuth === false || config.retryAuth === false || this.isSessionEndpoint(config) || !config.sentAccessToken) throw error
    const state = useAuthStore.getState()
    if (config.sessionGeneration !== state.generation) throw error
    if (config.sessionRetry) {
      // A persistent denial ends this session; Query must not restart it.
      if (config.sentAccessToken === state.accessToken) {
        useSessionActionStore.getState().setSessionNotice('invalidated')
        state.clear()
        if (config.sessionAttempt) this.coordinator.invalidate(config.sessionAttempt)
      }
      throw error
    }
    config.sessionRetry = true
    try {
      // A late 401 for the previous access token can use the already renewed
      // token. It must not consume the newly rotated cookie a second time.
      if (config.sentAccessToken === state.accessToken) await this.refreshSession()
      const current = useAuthStore.getState()
      if (current.generation !== config.sessionGeneration || !current.accessToken) throw error
      return this.axios.request(config)
    } catch { throw error }
  }
  private scheduleSession<T>(operation: () => Promise<T>): Promise<T> {
    const task = this.sessionTask.then(operation)
    this.sessionTask = task.then(() => undefined, () => undefined)
    return task
  }
  private acceptSession(data: BrowserSession): BrowserSession {
    const user = data?.user
    if (!data || typeof data.access_token !== 'string' || !data.access_token || data.token_type !== 'Bearer' || typeof data.expires_at !== 'string' || !user || typeof user.id !== 'string' || !user.id || typeof user.email !== 'string' || typeof user.name !== 'string' || !['BORROWER', 'STAFF', 'ADMIN'].includes(user.role) || user.is_active !== true) {
      throw new ApiRequestError('Invalid session response', 502)
    }
    // Explicit safe projection; never decode JWTs or retain arbitrary fields.
    const account: AuthUser = { id: user.id, email: user.email, name: user.name, role: user.role, is_active: user.is_active }
    return { access_token: data.access_token, expires_at: data.expires_at, token_type: 'Bearer', user: account }
  }
  private failSession(error: unknown, generation: number): void {
    if (useAuthStore.getState().generation !== generation) return
    const status = this.normalize(error).status
    if (status === 401 || status === 403) {
      if (useAuthStore.getState().isAuthenticated) useSessionActionStore.getState().setSessionNotice('invalidated')
    } else useSessionActionStore.getState().setSessionNotice(null)
    useAuthStore.getState().clear(status === 401 || status === 403 ? 'unauthenticated' : 'error')
  }
  refreshSession(): Promise<BrowserSession> {
    if (this.refreshPromise) return this.refreshPromise
    const generation = useAuthStore.getState().generation
    const task = this.scheduleSession(async () => {
      try {
        return await this.coordinator.run(true, async (attempt) => {
          if (this.disposed || generation !== useAuthStore.getState().generation) throw new ApiRequestError('Session ended', 401, 'SESSION_CHANGED')
          try {
            const data = this.acceptSession(await this.post<BrowserSession>('/auth/refresh', undefined, { attachAuth: false, retryAuth: false }))
            if (this.disposed || generation !== useAuthStore.getState().generation) throw new ApiRequestError('Session ended', 401, 'SESSION_CHANGED')
            useSessionActionStore.getState().setSessionNotice(null)
            useAuthStore.getState().setSession(data)
            this.coordinator.succeeded(attempt)
            return data
          } catch (error) {
            // An older operation cannot announce invalidation after logout or
            // account change. Uncoordinated 401s can be lost races, not global
            // proof of invalidation. Origin denial and network/5xx stay local.
            if (!this.disposed && generation === useAuthStore.getState().generation) {
              this.coordinator.failed(attempt, classifyRefreshFailure(this.normalize(error), attempt.exclusive))
            }
            throw error
          }
        })
      } catch (error) { if (!this.disposed) this.failSession(error, generation); throw error }
    })
    this.refreshPromise = task
    const finished = () => { if (this.refreshPromise === task) this.refreshPromise = null }
    void task.then(finished, finished)
    return task
  }
  bootstrap(): Promise<void> {
    if (this.bootstrapPromise) return this.bootstrapPromise
    clearLegacyAuthStorage()
    if (useAuthStore.getState().status !== 'idle') return Promise.resolve()
    useAuthStore.getState().setStatus('bootstrapping')
    this.bootstrapPromise = this.refreshSession().then(() => undefined, () => undefined)
    return this.bootstrapPromise
  }
  retryBootstrap(): Promise<void> {
    if (useAuthStore.getState().status !== 'error') return this.bootstrapPromise ?? Promise.resolve()
    this.bootstrapPromise = null
    useAuthStore.getState().setStatus('idle')
    return this.bootstrap()
  }
  authenticate(path: '/auth/login' | '/auth/register', input: unknown): Promise<BrowserSession> {
    this.authChange++
    useSessionActionStore.getState().setSessionNotice(null)
    const generation = useAuthStore.getState().clear('bootstrapping')
    // Credentials may change accounts. Peers discard old presentation state;
    // they stay unauthenticated until deliberate server bootstrap/login.
    this.coordinator.announce('auth-state-changed')
    return this.scheduleSession(async () => {
      try {
        return await this.coordinator.run(false, async () => {
          if (this.disposed || generation !== useAuthStore.getState().generation) throw new ApiRequestError('Session ended', 401, 'SESSION_CHANGED')
          const session = this.acceptSession(await this.post<BrowserSession>(path, input, { attachAuth: false, retryAuth: false }))
          if (this.disposed || generation !== useAuthStore.getState().generation) throw new ApiRequestError('Session ended', 401, 'SESSION_CHANGED')
          useAuthStore.getState().setSession(session)
          this.coordinator.announce('auth-state-changed')
          return session
        })
      } catch (error) { if (!this.disposed) this.failSession(error, generation); throw error }
    })
  }
  logout(): Promise<void> {
    // Invalidate state immediately, even when the network is unavailable.
    const authChange = this.authChange
    useSessionActionStore.getState().setSessionNotice(null)
    useAuthStore.getState().clear()
    this.coordinator.announce('logout')
    return this.scheduleSession(() => this.coordinator.run(false, async () => {
      // A newer login intent supersedes an old queued logout. Otherwise wait
      // for pending replies, then revoke the cookie they actually installed.
      // Other logout/invalidation notifications must NOT cancel revocation:
      // simultaneous logout otherwise lets every document skip the server.
      if (this.disposed || authChange !== this.authChange) return
      try { await this.post<void>('/auth/logout', undefined, { attachAuth: false, retryAuth: false }) }
      finally {
        if (!this.disposed && authChange === this.authChange) this.coordinator.announce('logout')
      }
    }))
  }
  private normalize(error: unknown): ApiRequestError { return normalizeApiError(error) }
}

export const apiClient = new ApiClient(undefined, sessionCoordinator)
if (import.meta.hot) import.meta.hot.dispose(() => apiClient.dispose())
