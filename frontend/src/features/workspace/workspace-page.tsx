import { useLocation, useNavigate } from 'react-router-dom'
import { Bell, LogOut } from 'lucide-react'
import { borrowerDestinations } from '@/components/application/borrower-navigation'
import { AccountTerms } from '@/features/terms/components/account-terms'
import { BorrowerShell, StaffShell } from '@/components/application/shells'
import { AppButton, PageHeading, SurfaceCard, ThemeControl } from '@/components/application/visual'
import { EmptyState } from '@/components/feedback/states'
import { useAuthStore } from '@/stores/auth-store'
import { useLogout } from '@/features/auth/hooks/use-auth'

const labels: Record<string, string> = {
  '/borrower/home': 'Home', '/borrower/equipment': 'Equipment Catalog', '/borrower/borrowings': 'My Borrowings',
  '/borrower/account': 'Account', '/borrower/notifications': 'Notifications',
  '/staff/dashboard': 'Dashboard', '/staff/requests': 'Requests & Borrowings', '/staff/inventory': 'Inventory',
  '/staff/borrowers': 'Borrowers', '/staff/account': 'Account', '/staff/notifications': 'Notifications',
  '/admin/reports': 'Reports', '/admin/administration': 'Administration',
}
const operationalDestinations = { Dashboard: '/staff/dashboard', 'Requests & Borrowings': '/staff/requests', Inventory: '/staff/inventory', Borrowers: '/staff/borrowers', Reports: '/admin/reports', Administration: '/admin/administration', Account: '/staff/account', Notifications: '/staff/notifications' }

export function WorkspacePage() {
  const user = useAuthStore(s => s.user)
  const location = useLocation()
  const navigate = useNavigate()
  const logout = useLogout()
  if (!user) return null
  const title = labels[location.pathname] ?? 'Workspace'
  const borrower = user.role === 'BORROWER'
  const account = location.pathname.endsWith('/account')
  const signOut = <AppButton variant="outline" className="workspace-signout" aria-label="Sign Out" onClick={() => { void logout() }}><LogOut size={18} aria-hidden="true" /><span>Sign Out</span></AppButton>
  const body = <>
    <PageHeading title={title} description={account ? 'Your FSMO account.' : `Welcome, ${user.name}.`} action={borrower ? <div className="workspace-actions"><ThemeControl />{signOut}</div> : undefined} />
    {account ? <SurfaceCard className="workspace-placeholder workspace-account"><h2>Account details</h2><dl><dt>Name</dt><dd>{user.name}</dd><dt>Email</dt><dd>{user.email}</dd><dt>Role</dt><dd>{user.role === 'BORROWER' ? 'Borrower' : user.role === 'STAFF' ? 'Staff' : 'Admin'}</dd></dl></SurfaceCard>
      : <SurfaceCard className="workspace-placeholder"><EmptyState title="This feature is not available yet" description="Your account is connected. Contact FSMO for assistance while this workspace is being prepared." /></SurfaceCard>}
    {account && borrower && <AccountTerms />}
  </>
  return borrower ? <BorrowerShell active={title === 'Equipment Catalog' ? 'Equipment' : title} destinations={borrowerDestinations}
    action={<AppButton variant="ghost" aria-label="Notifications" onClick={() => navigate('/borrower/notifications')}><Bell size={20} aria-hidden="true" /></AppButton>}>{body}</BorrowerShell>
    : <StaffShell active={title} audience={user.role === 'ADMIN' ? 'admin' : 'staff'} identity={user} destinationHrefs={operationalDestinations} footerAction={signOut}
      dashboardHref="/staff/dashboard" requestsHref="/staff/requests" onUnavailable={label => navigate(operationalDestinations[label as keyof typeof operationalDestinations] ?? '/staff/account')}>{body}</StaffShell>
}
