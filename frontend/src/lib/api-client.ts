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
import { clearLegacyAuthStorage } from '@/lib/storage'

const baseURL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api/v1'

export class ApiRequestError extends Error {
  status: number
  code?: string
  fields?: Record<string, string>
  constructor(message: string, status: number, code?: string, fields?: Record<string, string>) {
    super(message)
    this.name = 'ApiRequestError'
    this.status = status
    this.code = code
    this.fields = fields
  }
}
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
}

// One transport and one refresh flight per application instance. Session
// mutations are serialized so logout revokes the cookie from a pending reply.
export class ApiClient {
  private axios: AxiosInstance
  private refreshPromise: Promise<BrowserSession> | null = null
  private bootstrapPromise: Promise<void> | null = null
  private sessionTask: Promise<unknown> = Promise.resolve()

  constructor(adapter?: AxiosAdapter) {
    this.axios = axios.create({ baseURL, timeout: 15_000, headers: { 'Content-Type': 'application/json' }, withCredentials: true, adapter })
    this.axios.interceptors.request.use(this.attachToken)
    this.axios.interceptors.response.use((response) => response, (error: AxiosError<ApiResponse>) => this.onError(error))
  }
  get<T>(url: string, opts: RequestOptions = {}): Promise<T> { return this.request<T>('GET', url, undefined, opts) }
  post<T>(url: string, body?: unknown, opts: RequestOptions = {}): Promise<T> { return this.request<T>('POST', url, body, opts) }
  patch<T>(url: string, body?: unknown, opts: RequestOptions = {}): Promise<T> { return this.request<T>('PATCH', url, body, opts) }
  put<T>(url: string, body?: unknown, opts: RequestOptions = {}): Promise<T> { return this.request<T>('PUT', url, body, opts) }
  delete<T>(url: string, opts: RequestOptions = {}): Promise<T> { return this.request<T>('DELETE', url, undefined, opts) }

  private async request<T>(method: string, url: string, body: unknown, opts: RequestOptions): Promise<T> {
    try {
      const response = await this.axios.request<ApiResponse<T>>({ ...opts, method, url, data: body })
      if (response.status === 204) return undefined as T
      if (!response.data.success) throw new ApiRequestError(response.data.message || 'Request failed', response.status, response.data.error?.code, response.data.error?.fields)
      return response.data.data
    } catch (error) { throw this.normalize(error) }
  }
  private attachToken = (original: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
    const config = original as SessionRequest
    const state = useAuthStore.getState()
    if (config.sessionRetry && config.sessionGeneration !== state.generation) {
      throw new ApiRequestError('Session ended', 401)
    }
    config.sessionGeneration = state.generation
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
    const config = error.config as SessionRequest | undefined
    if (!config || error.response?.status !== 401 || config.attachAuth === false || config.retryAuth === false || this.isSessionEndpoint(config) || !config.sentAccessToken) throw error
    const state = useAuthStore.getState()
    if (config.sessionGeneration !== state.generation) throw error
    if (config.sessionRetry) {
      // A persistent denial ends this session; Query must not restart it.
      state.clear()
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
    if (!data || typeof data.access_token !== 'string' || !data.access_token || data.token_type !== 'Bearer' || typeof data.expires_at !== 'string' || !user || typeof user.id !== 'string' || !user.id || typeof user.email !== 'string' || typeof user.name !== 'string' || !['user', 'admin'].includes(user.role) || user.is_active !== true) {
      throw new ApiRequestError('Invalid session response', 502)
    }
    // Explicit safe projection; never decode JWTs or retain arbitrary fields.
    const account: AuthUser = { id: user.id, email: user.email, name: user.name, role: user.role, is_active: user.is_active }
    return { access_token: data.access_token, expires_at: data.expires_at, token_type: 'Bearer', user: account }
  }
  private failSession(error: unknown, generation: number): void {
    if (useAuthStore.getState().generation !== generation) return
    const status = this.normalize(error).status
    useAuthStore.getState().clear(status === 401 || status === 403 ? 'unauthenticated' : 'error')
  }
  refreshSession(): Promise<BrowserSession> {
    if (this.refreshPromise) return this.refreshPromise
    const generation = useAuthStore.getState().generation
    const task = this.scheduleSession(async () => {
      if (generation !== useAuthStore.getState().generation) throw new ApiRequestError('Session ended', 401)
      try {
        const data = this.acceptSession(await this.post<BrowserSession>('/auth/refresh', undefined, { attachAuth: false, retryAuth: false }))
        if (generation !== useAuthStore.getState().generation) throw new ApiRequestError('Session ended', 401)
        useAuthStore.getState().setSession(data)
        return data
      } catch (error) { this.failSession(error, generation); throw error }
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
    const generation = useAuthStore.getState().clear('bootstrapping')
    return this.scheduleSession(async () => {
      if (generation !== useAuthStore.getState().generation) throw new ApiRequestError('Session ended', 401)
      try {
        const session = this.acceptSession(await this.post<BrowserSession>(path, input, { attachAuth: false, retryAuth: false }))
        if (generation !== useAuthStore.getState().generation) throw new ApiRequestError('Session ended', 401)
        useAuthStore.getState().setSession(session)
        return session
      } catch (error) { this.failSession(error, generation); throw error }
    })
  }
  logout(): Promise<void> {
    // Invalidate state immediately, even when the network is unavailable.
    useAuthStore.getState().clear()
    return this.scheduleSession(() => this.post<void>('/auth/logout', undefined, { attachAuth: false, retryAuth: false }))
  }
  private normalize(error: unknown): ApiRequestError {
    if (error instanceof ApiRequestError) return error
    if (axios.isAxiosError<ApiResponse>(error)) {
      const payload = error.response?.data
      return new ApiRequestError(payload?.message || 'Request failed', error.response?.status ?? 0, payload?.error?.code, payload?.error?.fields)
    }
    return new ApiRequestError('Request failed', 0)
  }
}

export const apiClient = new ApiClient()
