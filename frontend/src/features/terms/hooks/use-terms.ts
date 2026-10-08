import { useMutation, useQuery } from '@tanstack/react-query'
import { queryClient } from '@/app/query-client'
import { queryKeys } from '@/lib/query-keys'
import { useAuthStore } from '@/stores/auth-store'
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
