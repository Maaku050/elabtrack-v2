import type { AxiosError } from 'axios'

/** One safe failure representation; never retains Axios config/headers/body. */
export class ApiRequestError extends Error {
  status: number
  code: string
  requestId?: string
  fields?: Record<string, string>
  constructor(message: string, status: number, code = 'REQUEST_FAILED', fields?: Record<string, string>, requestId?: string) {
    super(message)
    this.name = 'ApiRequestError'
    this.status = status
    this.code = code
    this.fields = fields
    this.requestId = requestId
  }
}

declare module '@tanstack/react-query' {
  interface Register { defaultError: ApiRequestError }
}

const publicCodes = new Set(['VALIDATION_ERROR', 'BAD_REQUEST', 'UNAUTHORIZED', 'INVALID_CREDENTIALS', 'TOKEN_INVALID', 'FORBIDDEN', 'NOT_FOUND', 'METHOD_NOT_ALLOWED', 'CONFLICT', 'PAYLOAD_TOO_LARGE', 'UNSUPPORTED_MEDIA_TYPE', 'RATE_LIMITED', 'INTERNAL_ERROR', 'SERVICE_UNAVAILABLE', 'TERMS_NOT_PUBLISHED', 'TERMS_VERSION_NOT_FOUND', 'TERMS_VERSION_CHANGED', 'TERMS_ACCEPTANCE_REQUIRED', 'TERMS_VERSION_EXISTS', 'TERMS_PUBLICATION_CHANGED'])
const record = (value: unknown): value is Record<string, unknown> => typeof value === 'object' && value !== null && !Array.isArray(value)
const safeText = (value: unknown): value is string => typeof value === 'string' && value.length > 0 && value.length <= 256 && [...value].every((character) => character.charCodeAt(0) >= 32 && character.charCodeAt(0) !== 127)
const safeID = (value: unknown): value is string => typeof value === 'string' && /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/.test(value)

export function normalizeApiResponseError(status: number, payload: unknown, headerID?: unknown): ApiRequestError {
  const error = record(payload) && payload.success === false && record(payload.error) ? payload.error : undefined
  const valid = error && typeof error.code === 'string' && publicCodes.has(error.code) && safeText(error.message) && safeID(error.requestId)
  const requestId = safeID(headerID) ? headerID : error && safeID(error.requestId) ? error.requestId : undefined
  // A reverse proxy/nonstandard error payload is never displayed as raw text.
  const code = valid ? error.code as string : status >= 500 ? 'INTERNAL_ERROR' : 'REQUEST_FAILED'
  const message = status >= 500 ? 'Something went wrong. Please try again later.' : valid ? error.message as string : 'Request failed. Please try again.'
  let fields: Record<string, string> | undefined
  if (valid && code === 'VALIDATION_ERROR' && record(error.fields)) {
    fields = Object.fromEntries(Object.entries(error.fields).filter(([key, value]) => /^[a-z][a-z0-9_]{0,63}$/.test(key) && key !== '__proto__' && safeText(value))) as Record<string, string>
  }
  return new ApiRequestError(message, status, code, fields, requestId)
}

export function normalizeApiError(error: unknown): ApiRequestError {
  if (error instanceof ApiRequestError) return error
  if (record(error) && error.isAxiosError === true) {
    const failure = error as unknown as AxiosError<unknown>
    if (!failure.response) return new ApiRequestError('Unable to connect. Please try again.', 0, 'NETWORK_ERROR')
    const headers = failure.response.headers
    const id = typeof headers.get === 'function' ? headers.get('x-request-id') : headers['x-request-id'] ?? headers['X-Request-ID']
    return normalizeApiResponseError(failure.response.status, failure.response.data, id)
  }
  return new ApiRequestError('Request failed. Please try again.', 0)
}

/** Optional support reference for unexpected failures, not routine field errors. */
export function apiErrorMessage(error: unknown): string {
  const safe = normalizeApiError(error)
  return safe.status >= 500 && safe.requestId ? `${safe.message} Reference: ${safe.requestId}` : safe.message
}
