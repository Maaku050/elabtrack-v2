import { Link } from 'react-router-dom'
import { Bell } from 'lucide-react'
import { AppButton } from '@/components/application/visual'
import { useAuthStore } from '@/stores/auth-store'
import { useNotificationCount } from '../hooks/use-notifications'
function CountedBell({ href, className }: { href: string; className?: string }) {
 const query = useNotificationCount(), count = query.data?.unread
 return <AppButton className={`notification-bell ${className??''}`} variant="ghost" aria-label="Notifications" aria-description={query.isError ? 'Unread count unavailable' : count === undefined ? 'Loading unread count' : `${count} unread`} role="link" nativeButton={false} render={<Link to={href} />}><Bell aria-hidden="true" />{count !== undefined && count > 0 && <span className="rounded-full bg-primary px-1.5 text-xs font-semibold text-primary-foreground" aria-hidden="true">{count > 99 ? '99+' : count}</span>}</AppButton>
}
export function NotificationBell({ className }: { className?: string }) {
 const user = useAuthStore(s => s.user), authenticated = useAuthStore(s => s.isAuthenticated), href = user?.role === 'BORROWER' ? '/borrower/notifications' : '/staff/notifications'
 return authenticated && user ? <CountedBell href={href} className={`notification-bell ${className??''}`} /> : <AppButton className={`notification-bell ${className??''}`} variant="ghost" aria-label="Notifications" role="link" nativeButton={false} render={<Link to={href} />}><Bell aria-hidden="true" /></AppButton>
}
