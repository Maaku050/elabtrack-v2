import { useRef, useState } from 'react'
import { Link } from 'react-router-dom'
import { Bell, CheckCheck, ChevronRight } from 'lucide-react'
import { AppButton, ToneBadge } from '@/components/application/visual'
import { ManagementPage, Notice } from '@/components/application/management'
import { BorrowerShell } from '@/components/application/shells'
import { borrowerDestinations } from '@/components/application/borrower-navigation'
import { DirectoryPanel, DirectorySelect } from '@/components/application/directory'
import { ServerPagination } from '@/components/application/server-pagination'
import { useDirectoryState } from '@/hooks/use-directory-state'
import { useAuthStore } from '@/stores/auth-store'
import { apiErrorMessage } from '@/lib/api-error'
import { manilaTimeCompact } from '@/features/borrowing/lib/time'
import { useNotifications, useNotificationReadState, useNotificationsMarkAll } from '../hooks/use-notifications'
export function NotificationCenter(){
 const role=useAuthStore(s=>s.user?.role),view=useDirectoryState(),read=view.value('read'),filter=read==='read'||read==='unread'?read:'',query=useNotifications({page:view.page(),per_page:25,read:filter}),mutation=useNotificationReadState(),bulk=useNotificationsMarkAll(),data=query.data,[notice,setNotice]=useState(''),busy=useRef(false)
 async function markAll(){if(busy.current||!data?.unread)return;busy.current=true;setNotice('');try{const out=await bulk.mutateAsync();setNotice(`${out.updated} notifications marked read.`)}catch{/* Visible error below. */}finally{busy.current=false}}
 const body=<ManagementPage title="Notifications" description="Your FSMO updates · timestamps in Asia/Manila." domain="In-app notifications" actions={<AppButton variant="outline" disabled={!data?.unread||bulk.isPending||mutation.isPending} onClick={()=>{void markAll()}}><CheckCheck aria-hidden="true"/>{bulk.isPending?'Marking read…':'Mark All Read'}</AppButton>}>{notice&&<Notice tone="success">{notice}</Notice>}{(mutation.error||bulk.error)&&<Notice tone="danger">{apiErrorMessage(mutation.error??bulk.error)}</Notice>}<DirectoryPanel title="Notification center" summary={data?`${data.unread} unread · ${data.total} matching notifications`:'Your authorized updates'} pending={query.isPending} error={query.isError?apiErrorMessage(query.error):undefined} onRetry={()=>{void query.refetch()}} empty={!data?.items.length} toolbar={<DirectorySelect label="Read state" value={filter} onChange={e=>view.update({read:e.target.value})}><option value="">All notifications</option><option value="unread">Unread</option><option value="read">Read</option></DirectorySelect>} footer={data&&<ServerPagination label="Notifications" page={data.page} perPage={data.per_page} total={data.total} totalPages={Math.ceil(data.total/data.per_page)} count={data.items.length} pending={query.isFetching} onPage={p=>view.update({page:String(p)},false)}/>}><ul className="notification-rows">{data?.items.map(n=><li key={n.id} data-unread={!n.read_at}><span className="notification-event-icon"><Bell aria-hidden="true"/></span><div className="notification-copy"><div><h2>{n.title}</h2>{!n.read_at&&<ToneBadge tone="primary">Unread</ToneBadge>}</div><p>{n.body}</p><time dateTime={n.created_at}>{manilaTimeCompact(n.created_at)}{n.read_at?' · Read':''}</time></div><div className="notification-actions"><AppButton variant="ghost" role="link" nativeButton={false} render={<Link to={`${role==='BORROWER'?'/borrower/borrowings':'/staff/requests'}/${n.borrowing_id}`}/>}>Open Borrowing<ChevronRight aria-hidden="true"/></AppButton><AppButton variant="ghost" disabled={mutation.isPending||bulk.isPending} aria-label={`${n.read_at?'Mark unread':'Mark read'}: ${n.title}`} onClick={()=>mutation.mutate({id:n.id,read:!n.read_at})}>{n.read_at?'Mark Unread':'Mark Read'}</AppButton></div></li>)}</ul></DirectoryPanel></ManagementPage>
 return role==='BORROWER'?<BorrowerShell active="" destinations={borrowerDestinations}>{body}</BorrowerShell>:body
}
