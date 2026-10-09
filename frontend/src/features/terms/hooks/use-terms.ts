import { useMutation, useQuery } from '@tanstack/react-query'
import { queryClient } from '@/app/query-client'
import { queryKeys } from '@/lib/query-keys'
import { useAuthStore } from '@/stores/auth-store'
import { ApiRequestError } from '@/lib/api-error'
import { termsApi } from '../api/terms.api'

export function useTermsStatus() {
 const user = useAuthStore(s => s.user)
 return useQuery({ queryKey: queryKeys.terms.status(user?.id),
  queryFn: ({ signal }) => termsApi.status(signal),
  enabled: user?.role === 'BORROWER' && user.is_active,
  meta: { authenticated: true }, staleTime: 0, retry: false,
  refetchOnMount: 'always', refetchOnWindowFocus: true, refetchOnReconnect: true,
 })
}
export function useAcceptTerms() {
 return useMutation({ mutationFn: (versionId: string) => termsApi.accept(versionId), retry: false,
  onSuccess: async () => { await queryClient.invalidateQueries({ queryKey: queryKeys.terms.all }) },
  onError: async error => {
   if (error.code === 'TERMS_VERSION_CHANGED' || error.code === 'TERMS_NOT_PUBLISHED') {
    await queryClient.invalidateQueries({ queryKey: queryKeys.terms.all })
   }
  },
 })
}

// Read-only current publication for authorized administrative information.
export function useCurrentTerms() {
 const actor = useAuthStore(s => s.user)
 return useQuery({
  queryKey: ['terms', actor?.id, 'current-publication'],
  queryFn: async ({ signal }) => {
   try { return await termsApi.current(signal) }
   catch (error) {
    // Absence is an expected policy state, not a failed/no-data query that
    // TanStack reloads on every mount. Other failures retain error semantics.
    if (error instanceof ApiRequestError && error.status === 503 && error.code === 'TERMS_NOT_PUBLISHED') return null
    throw error
   }
  },
  meta: { authenticated: true }, enabled: actor?.role === 'ADMIN',
  staleTime: 30_000, retry: false,
 })
}
