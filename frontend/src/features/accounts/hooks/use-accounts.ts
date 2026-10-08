import { useMutation, useQuery } from '@tanstack/react-query'
import { useAuthStore } from '@/stores/auth-store'
import { queryClient } from '@/app/query-client'
import { ApiRequestError } from '@/lib/api-error'
import { accountsApi } from '../api/accounts.api'
import type { AccountFilter } from '../types'
export function useAccounts(filter: AccountFilter, staff: boolean) {
 const actor = useAuthStore(s => s.user?.id)
 return useQuery({ queryKey: ['accounts', actor, staff, filter], queryFn: ({ signal }) => accountsApi.list(filter, staff, signal), meta: { authenticated: true }, enabled: !!actor })
}
export function useAccount(id: string, staff: boolean) {
 const actor = useAuthStore(s => s.user?.id)
 return useQuery({ queryKey: ['accounts', actor, 'detail', id, staff], queryFn: ({ signal }) => accountsApi.detail(id, staff, signal), meta: { authenticated: true }, enabled: !!actor && !!id })
}
export function useProvisioningPolicy() {
 const actor = useAuthStore(s => s.user)
 return useQuery({ queryKey: ['accounts', actor?.id, 'policy'], queryFn: ({ signal }) => accountsApi.policy(signal), meta: { authenticated: true }, enabled: actor?.role === 'ADMIN' })
}
export function useAccountAudit(id: string, staff: boolean, page: number) {
 const actor = useAuthStore(s => s.user)
 return useQuery({ queryKey: ['accounts', actor?.id, 'audit', id, staff, page], queryFn: ({ signal }) => accountsApi.audit(id, staff, page, signal), meta: { authenticated: true }, enabled: actor?.role === 'ADMIN' && !!id })
}
// A late reply cannot populate another account's cache or drive navigation.
export function useAccountCommand<I, O>(operation: (input: I) => Promise<O>) {
 return useMutation({ mutationFn: async (input: I) => { const generation = useAuthStore.getState().generation; const data = await operation(input); if (generation !== useAuthStore.getState().generation) throw new ApiRequestError('Session changed.', 401, 'SESSION_CHANGED'); return data }, onSuccess: async () => { await queryClient.invalidateQueries({ queryKey: ['accounts'] }) }, retry: false })
}
