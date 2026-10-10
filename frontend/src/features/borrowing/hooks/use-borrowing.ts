import { useMutation, useQueries, useQuery, useQueryClient } from '@tanstack/react-query'
import { useAuthStore } from '@/stores/auth-store'
import { ApiRequestError } from '@/lib/api-error'
import { borrowingApi } from '../api/borrowing.api'
import type { BorrowingFilter } from '../types'
export function useBorrowings(filter: BorrowingFilter) {
 const actor = useAuthStore(s => s.user?.id)
 return useQuery({ queryKey: ['borrowings', actor, 'list', filter], queryFn: ({ signal }) => borrowingApi.list(filter, signal), enabled: !!actor, meta: { authenticated: true } })
}
export function useBorrowing(id: string) {
 const actor = useAuthStore(s => s.user?.id)
 return useQuery({ queryKey: ['borrowings', actor, 'detail', id], queryFn: ({ signal }) => borrowingApi.detail(id, signal), enabled: !!actor && !!id, meta: { authenticated: true }, refetchInterval: query => ['PENDING', 'CHECKED_OUT'].includes(query.state.data?.status ?? '') ? 30_000 : false })
}
export function useBorrowingHistory<T=Record<string,unknown>>(id:string,kind:string,page:number,enabled=true,perPage=25){const actor=useAuthStore(s=>s.user?.id);return useQuery({queryKey:['borrowings',actor,'history',id,kind,page,perPage],queryFn:({signal})=>borrowingApi.history<T>(id,kind,page,perPage,signal),enabled:!!actor&&!!id&&enabled,meta:{authenticated:true}})}
export function useBorrowingCommand<I,O>(operation:(input:I)=>Promise<O>){
 const client=useQueryClient()
 const invalidate=async()=>{await Promise.all(['borrowings','inventory','accounts','reporting','notifications'].map(key=>client.invalidateQueries({queryKey:[key]})))}
 return useMutation({mutationFn:async(input:I)=>{const {generation,user}=useAuthStore.getState();const data=await operation(input);const now=useAuthStore.getState();if(generation!==now.generation||user?.id!==now.user?.id)throw new ApiRequestError('Session changed.',401,'SESSION_CHANGED');return data},onSuccess:invalidate,onError:async error=>{if(error instanceof ApiRequestError&&error.status===409)await invalidate()},retry:false})
}

// Count each supported state through the authoritative bounded list endpoint.
// These caches share command invalidation and never grant eligibility or privileges.
export const borrowingQueueStates = [
 {status:'PENDING',label:'Pending'}, {status:'CHECKED_OUT',label:'Active'},
 {status:'COMPLETED',label:'Completed'}, {status:'DENIED',label:'Denied'},
 {status:'CANCELLED',label:'Cancelled'}, {status:'EXPIRED',label:'Expired'}, {status:'',label:'All'},
] as const
export function useBorrowingQueueCounts(search: string, borrowerId?: string) {
 const actor=useAuthStore(s=>s.user?.id), role=useAuthStore(s=>s.user?.role)
 return useQueries({queries:borrowingQueueStates.map(({status})=>{
  const filter={page:1,per_page:1,status,search,...(borrowerId?{borrower_id:borrowerId}:{})}
  return {queryKey:['borrowings',actor,'queue-count',filter],queryFn:({signal}:{signal:AbortSignal})=>borrowingApi.list(filter,signal),enabled:!!actor&&role!=='BORROWER',staleTime:30_000,meta:{authenticated:true}}
 })})
}
