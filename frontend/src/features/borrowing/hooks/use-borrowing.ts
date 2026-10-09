import { useMutation, useQuery } from '@tanstack/react-query'
import { useAuthStore } from '@/stores/auth-store'
import { queryClient } from '@/app/query-client'
import { ApiRequestError } from '@/lib/api-error'
import { borrowingApi } from '../api/borrowing.api'
import type { BorrowingFilter } from '../types'
export function useBorrowings(filter: BorrowingFilter) {
 const actor = useAuthStore(s => s.user?.id)
 return useQuery({ queryKey: ['borrowings', actor, 'list', filter], queryFn: ({ signal }) => borrowingApi.list(filter, signal), enabled: !!actor, meta: { authenticated: true } })
}
export function useBorrowing(id: string) {
 const actor = useAuthStore(s => s.user?.id)
 return useQuery({ queryKey: ['borrowings', actor, 'detail', id], queryFn: ({ signal }) => borrowingApi.detail(id, signal), enabled: !!actor && !!id, meta: { authenticated: true }, refetchInterval: query => query.state.data?.status === 'PENDING' ? 30_000 : false })
}
async function invalidate() { await Promise.all(['borrowings', 'inventory', 'accounts'].map(key => queryClient.invalidateQueries({ queryKey: [key] }))) }
export function useBorrowingCommand<I, O>(operation: (input: I) => Promise<O>) {
 return useMutation({ mutationFn: async (input: I) => { const generation = useAuthStore.getState().generation; const data = await operation(input); if (generation !== useAuthStore.getState().generation) throw new ApiRequestError('Session changed.', 401, 'SESSION_CHANGED'); return data }, onSuccess: invalidate, onError: async error => { if (error instanceof ApiRequestError && error.status === 409) await invalidate() }, retry: false })
}
