import axios, {
  AxiosError,
  type AxiosInstance,
  type AxiosRequestConfig,
  type InternalAxiosRequestConfig,
} from 'axios'
import type { ApiResponse, ApiError } from '@/types/api'
import { storage, storageKeys } from '@/lib/storage'

const baseURL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api/v1'

/** Typed error thrown by the API client. */
export class ApiRequestError extends Error {
  code?: string
  fields?: Record<string, string>
  status: number

  constructor(message: string, status: number, code?: string, fields?: Record<string, string>) {
    super(message)
    this.name = 'ApiRequestError'
    this.status = status
    this.code = code
    this.fields = fields
  }
}

/** Options for the apiClient.request method. */
export interface RequestOptions extends Omit<AxiosRequestConfig, 'auth'> {
  /** Whether to attach the access token (default true). */
  attachAuth?: boolean
  /** AbortSignal passthrough. */
  signal?: AbortSignal
}

/**
 * Centralized API client.
 *
 * Responsibilities:
 *   - base URL
 *   - JSON serialization
 *   - Authorization header injection
 *   - 401 handling + single-flight token refresh
 *   - typed errors
 *   - AbortSignal support
 */
class ApiClient {
  private axios: AxiosInstance
  private refreshPromise: Promise<string> | null = null
  private onUnauthorized?: () => void

  constructor() {
    this.axios = axios.create({
      baseURL,
      timeout: 15_000,
      headers: { 'Content-Type': 'application/json' },
      withCredentials: true,
    })

    this.axios.interceptors.request.use(this.attachToken)
    this.axios.interceptors.response.use(
      (r) => r,
      (err: AxiosError<ApiResponse>) => this.onError(err),
    )
  }

  /** Register a callback invoked when the session is no longer valid. */
  setOnUnauthorized(cb: () => void): void {
    this.onUnauthorized = cb
  }

  /** GET wrapper returning the unwrapped `data` field. */
  async get<T>(url: string, opts: RequestOptions = {}): Promise<T> {
    return this.request<T>('GET', url, undefined, opts)
  }

  /** POST wrapper. */
  async post<T>(url: string, body?: unknown, opts: RequestOptions = {}): Promise<T> {
    return this.request<T>('POST', url, body, opts)
  }

  /** PATCH wrapper. */
  async patch<T>(url: string, body?: unknown, opts: RequestOptions = {}): Promise<T> {
    return this.request<T>('PATCH', url, body, opts)
  }

  /** PUT wrapper. */
  async put<T>(url: string, body?: unknown, opts: RequestOptions = {}): Promise<T> {
    return this.request<T>('PUT', url, body, opts)
  }

  /** DELETE wrapper. */
  async delete<T>(url: string, opts: RequestOptions = {}): Promise<T> {
    return this.request<T>('DELETE', url, undefined, opts)
  }

  /** Raw request method used by the verb wrappers. */
  private async request<T>(
    method: string,
    url: string,
    body: unknown,
    opts: RequestOptions,
  ): Promise<T> {
    try {
      const res = await this.axios.request<ApiResponse<T>>({
        method,
        url,
        data: body,
        ...opts,
        headers: { ...opts.headers },
      })
      // 204 No Content has no body.
      if (res.status === 204) return undefined as T
      const payload = res.data
      if (!payload.success) {
        throw this.toRequestError(payload, res.status)
      }
      return payload.data
    } catch (err) {
      throw this.normalize(err)
    }
  }

  private attachToken = (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
    const token = storage.get(storageKeys.accessToken)
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  }

  private async onError(err: AxiosError<ApiResponse>): Promise<unknown> {
    if (err.response?.status === 401 && !this.isRefreshCall(err.config)) {
      const refreshed = await this.tryRefresh()
      if (refreshed && err.config) {
        // Retry the original request once with the new token.
        const token = storage.get(storageKeys.accessToken)
        err.config.headers = err.config.headers ?? {}
        err.config.headers.Authorization = `Bearer ${token}`
        return this.axios.request(err.config)
      }
      // Refresh failed -> session is invalid.
      this.onUnauthorized?.()
    }
    return Promise.reject(err)
  }

  private isRefreshCall(config?: InternalAxiosRequestConfig): boolean {
    return !!config?.url?.includes('/auth/refresh')
  }

  /** Single-flight refresh: concurrent 401s share one refresh request. */
  private async tryRefresh(): Promise<string | null> {
    if (this.refreshPromise) return this.refreshPromise
    const refreshToken = storage.get(storageKeys.refreshToken)
    if (!refreshToken) return null

    this.refreshPromise = (async () => {
      try {
        const res = await this.axios.post<ApiResponse<{ access_token: string; refresh_token: string }>>(
          '/auth/refresh',
          { refresh_token: refreshToken },
        )
        const pair = res.data.data
        storage.set(storageKeys.accessToken, pair.access_token)
        storage.set(storageKeys.refreshToken, pair.refresh_token)
        return pair.access_token
      } catch {
        storage.remove(storageKeys.accessToken)
        storage.remove(storageKeys.refreshToken)
        return ''
      } finally {
        this.refreshPromise = null
      }
    })()

    return this.refreshPromise
  }

  private normalize(err: unknown): ApiRequestError {
    if (err instanceof ApiRequestError) return err
    if (axios.isAxiosError<ApiResponse>(err)) {
      const payload = err.response?.data
      if (payload && !payload.success && payload.error) {
        return new ApiRequestError(
          payload.message || payload.error.message || 'Request failed',
          err.response?.status ?? 0,
          payload.error.code,
          payload.error.fields,
        )
      }
      return new ApiRequestError(err.message, err.response?.status ?? 0)
    }
    if (err instanceof Error) return new ApiRequestError(err.message, 0)
    return new ApiRequestError('Unknown error', 0)
  }

  private toRequestError(payload: ApiError & { message?: string }, status: number): ApiRequestError {
    return new ApiRequestError(
      payload.message ?? 'Request failed',
      status,
      payload.code,
      payload.fields,
    )
  }
}

/** Singleton API client used across the app. */
export const apiClient = new ApiClient()
