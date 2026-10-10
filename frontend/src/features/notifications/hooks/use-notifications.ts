import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useAuthStore } from '@/stores/auth-store'
import { ApiRequestError } from '@/lib/api-error'
import { notificationsApi } from '../api/notifications.api'
import type { NotificationFilter } from '../types'
export function useNotifications(filter: NotificationFilter) {
 const user = useAuthStore(s => s.user)
 return useQuery({ queryKey: ['notifications', user?.id, user?.role, 'list', filter], queryFn: ({ signal }) => notificationsApi.list(filter, signal), enabled: !!user, meta: { authenticated: true }, staleTime: 10_000, refetchInterval: 30_000 })
}
export function useNotificationCount() {
 const user = useAuthStore(s => s.user)
 return useQuery({ queryKey: ['notifications', user?.id, user?.role, 'count'], queryFn: ({ signal }) => notificationsApi.count(signal), enabled: !!user, meta: { authenticated: true }, staleTime: 10_000, refetchInterval: 30_000 })
}
export function useNotificationReadState() {
 const client = useQueryClient()
 return useMutation({ mutationFn: async ({ id, read }: { id: string; read: boolean }) => { const generation = useAuthStore.getState().generation; const value = await notificationsApi.mark(id, read); if (generation !== useAuthStore.getState().generation) throw new ApiRequestError('Session changed.', 401, 'SESSION_CHANGED'); return value }, onSuccess: () => client.invalidateQueries({ queryKey: ['notifications'] }), retry: false })
}

export function useNotificationsMarkAll(){const client=useQueryClient();return useMutation({retry:false,mutationFn:async()=>{const generation=useAuthStore.getState().generation;const out=await notificationsApi.markAll();if(generation!==useAuthStore.getState().generation)throw new ApiRequestError('Session changed.',401,'SESSION_CHANGED');return out},onSuccess:()=>client.invalidateQueries({queryKey:['notifications']})})}
