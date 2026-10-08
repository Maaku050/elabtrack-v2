import type { ReactNode } from 'react'
import { useNavigate } from 'react-router-dom'
import { LogOut } from 'lucide-react'
import { StaffShell } from './shells'
import { AppButton } from './visual'
import { useAuthStore } from '@/stores/auth-store'
import { useLogout } from '@/features/auth/hooks/use-auth'
const operationalDestinations = { Dashboard: '/staff/dashboard', 'Requests & Borrowings': '/staff/requests', Inventory: '/staff/inventory', Borrowers: '/staff/borrowers', Reports: '/admin/reports', Administration: '/admin/administration', Account: '/staff/account', Notifications: '/staff/notifications' }
export function OperationalShell({ active, children }: { active: string; children: ReactNode }) {
 const user = useAuthStore(s => s.user); const navigate = useNavigate(); const logout = useLogout()
 if (!user) return null
 return <StaffShell active={active} audience={user.role === 'ADMIN' ? 'admin' : 'staff'} identity={user} destinationHrefs={operationalDestinations} dashboardHref="/staff/dashboard" requestsHref="/staff/requests" onUnavailable={label => navigate(operationalDestinations[label as keyof typeof operationalDestinations] ?? '/staff/account')} footerAction={<AppButton variant="outline" className="workspace-signout" onClick={() => { void logout() }}><LogOut size={18} aria-hidden="true" />Sign Out</AppButton>}>{children}</StaffShell>
}
