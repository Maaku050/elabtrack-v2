/** Retained success envelope; 204 has no body. */
export interface ApiSuccess<T> {
  success: true
  message: string
  data: T
  meta: PageMeta | null
}
export interface ApiFailure {
  success: false
  message: string
  data: null
  meta: null
  error: ApiError
}
export type ApiResponse<T = unknown> = ApiSuccess<T> | ApiFailure

export interface ApiError {
  code: string
  message: string
  requestId: string
  fields?: Record<string, string>
}
export interface PageMeta {
  page: number
  per_page: number
  total: number
  last_page: number
}
export interface Paginated<T> {
  items: T[]
  meta: PageMeta
}
