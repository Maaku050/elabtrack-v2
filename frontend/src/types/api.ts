/** Standard API response envelope returned by the Go backend. */
export interface ApiResponse<T = unknown> {
  success: boolean
  message: string
  data: T
  meta: PageMeta | null
  error: ApiError | null
}

export interface ApiError {
  message: string
  code?: string
  fields?: Record<string, string>
}

export interface PageMeta {
  page: number
  per_page: number
  total: number
  last_page: number
}

/** Paginated response convenience type. */
export interface Paginated<T> {
  items: T[]
  meta: PageMeta
}
