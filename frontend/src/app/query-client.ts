import { QueryClient } from '@tanstack/react-query'

/**
 * TanStack Query client configured for the app.
 *
 * - Retries are limited to avoid hammering the API on hard errors.
 * - `refetchOnWindowFocus` is off by default for SPA ergonomics; enable
 *   per-query if a screen needs always-fresh data.
 */
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      staleTime: 30_000,
      gcTime: 5 * 60_000,
      refetchOnWindowFocus: false,
    },
    mutations: {
      retry: 0,
    },
  },
})
