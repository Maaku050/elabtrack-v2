import { QueryClient, type Query } from '@tanstack/react-query'
import { ApiRequestError } from '@/lib/api-error'

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
      retry: (failureCount, error) => {
        if (error instanceof ApiRequestError && error.status >= 400 && error.status < 500) return false
        return failureCount < 1
      },
      staleTime: 30_000,
      gcTime: 5 * 60_000,
      refetchOnWindowFocus: false,
    },
    mutations: {
      retry: 0,
    },
  },
})

/** Retained private roots plus explicit opt-in for future approved features.
 * Unclassified new server queries must use meta.authenticated=true if private. */
export function isAuthenticatedQuery(query: Query): boolean {
  return query.meta?.authenticated === true || ['auth', 'users'].includes(String(query.queryKey[0]))
}

export function clearAuthenticatedQueries(): void {
  // removeQueries cancels active work synchronously (including AbortSignal)
  // before destroying the cache. Public health/status data survives logout.
  queryClient.removeQueries({ predicate: isAuthenticatedQuery })
  // Mutations may contain credentials or private inputs. Removal cannot undo
  // an HTTP mutation; application callbacks still fence by session generation.
  queryClient.getMutationCache().clear()
}
