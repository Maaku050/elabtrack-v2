import { Outlet, useLocation } from 'react-router-dom'
import { OperationalFrame } from './operational-frame'
import { useAuthStore } from '@/stores/auth-store'
import { useLogout } from '@/features/auth/hooks/use-auth'

const operationalDestinations = {
  Dashboard: '/staff/dashboard',
  'Requests & Borrowings': '/staff/requests',
  Inventory: '/staff/inventory',
  Borrowers: '/staff/borrowers',
  Reports: '/admin/reports',
  Administration: '/admin/administration',
  Account: '/staff/account',
  Notifications: '/staff/notifications',
}

/** Persistent presentation only. Nested ProtectedRoute retains server authority. */
export function OperationalLayout() {
  const location = useLocation()
  const user = useAuthStore(s => s.user)
  const logout = useLogout()
  if (!user || !['ADMIN', 'STAFF'].includes(user.role)) return <Outlet />

  const active = location.pathname.includes('/inventory') ? 'Inventory'
    : location.pathname.includes('/borrowers') ? 'Borrowers'
    : Object.entries(operationalDestinations).find(([, href]) => location.pathname.startsWith(href))?.[0] ?? 'Workspace'
  const signOut = () => { void logout() }

  return <OperationalFrame identity={user} active={active} onSignOut={signOut}><Outlet /></OperationalFrame>
}
