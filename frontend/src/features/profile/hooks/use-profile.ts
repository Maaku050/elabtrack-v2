import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useAuthStore } from '@/stores/auth-store'
import { ApiRequestError } from '@/lib/api-error'
import { profileApi } from '../api/profile.api'
export function useProfile(id: string) {
 const actor = useAuthStore(s=>s.user)
 return useQuery({queryKey:['profile',actor?.id,actor?.role,id],queryFn:({signal})=>profileApi.metadata(id,signal),enabled:!!actor&&!!id,staleTime:60_000,retry:false,meta:{authenticated:true}})
}
export function useProfileCommand<I,O>(operation:(input:I)=>Promise<O>){
 const client=useQueryClient()
 return useMutation({retry:false,mutationFn:async(input:I)=>{const before=useAuthStore.getState();const result=await operation(input);const next=useAuthStore.getState();if(before.generation!==next.generation||before.user?.id!==next.user?.id)throw new ApiRequestError('Session changed.',401,'SESSION_CHANGED');return result},onSuccess:async()=>{await Promise.all(['profile','accounts','reporting'].map(key=>client.invalidateQueries({queryKey:[key]})))},onError:async error=>{if(error instanceof ApiRequestError&&error.status===409)await client.invalidateQueries({queryKey:['profile']})}})
}

export function useOwnAccount(){const user=useAuthStore(s=>s.user);return useQuery({queryKey:['accounts',user?.id,user?.role,'me'],queryFn:({signal})=>profileApi.own(signal),enabled:!!user,staleTime:30000,retry:false,meta:{authenticated:true}})}
