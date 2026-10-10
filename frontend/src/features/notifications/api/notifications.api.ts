import type { Notification, NotificationFilter, NotificationPage } from '../types'
const transport = async () => (await import('@/lib/api-client')).apiClient
export const notificationsApi = {
 list: async (filter: NotificationFilter, signal?: AbortSignal): Promise<NotificationPage> => (await transport()).get('/notifications', { params: filter, signal }),
 count: async (signal?: AbortSignal): Promise<{ unread: number }> => (await transport()).get('/notifications/count', { signal }),
 markAll: async ():Promise<{updated:number}> => (await transport()).post('/notifications/read-all',{confirm:true}),
 mark: async (id: string, read: boolean): Promise<Notification> => (await transport()).patch(`/notifications/${id}`, { read }),
}
