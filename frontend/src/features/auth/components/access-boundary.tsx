import { Navigate, Outlet, useLocation, Link } from 'react-router-dom'
import { useAuthStore } from '@/stores/auth-store'
import { useCurrentUser } from '../hooks/use-auth'
import { homeForRole } from '../navigation'
import { apiErrorMessage } from '@/lib/api-error'
import { AppBrand, AppButton, ThemeControl, SurfaceCard } from '@/components/application/visual'
import { LoadingState } from '@/components/feedback/states'
import type { UserRole } from '@/types/common'

export function AccessDenied({ embedded = false }: { embedded?: boolean }) {
  const role = useAuthStore(s => s.user?.role)
  const Frame = embedded ? 'section' : 'main'
  return <Frame className={embedded ? 'workspace-access-state' : 'auth-page'}>{!embedded && <div className="auth-top"><AppBrand /><ThemeControl /></div>}<SurfaceCard className="auth-card">
    <h1>Access denied</h1><p>Your account does not have access to this workspace.</p>
    <Link className="app-button auth-link" to={role ? homeForRole(role) : '/login'}>Return to your workspace</Link>
  </SurfaceCard></Frame>
}

export function ProtectedRoute({ roles, embedded = false }: { roles: readonly UserRole[]; embedded?: boolean }) {
  const location = useLocation()
  const status = useAuthStore(s => s.status)
  const current = useAuthStore(s => s.user)
  const account = useCurrentUser(location.pathname)
  const withinWorkspace = embedded && !!current && ['ADMIN', 'STAFF'].includes(current.role)
  const Frame = withinWorkspace ? 'section' : 'main'
  if (status === 'idle' || status === 'bootstrapping') return <LoadingState label="Restoring your session…" />
  if (status === 'error') return <LoadingState label="Reconnect to restore your session using the session recovery above." />
  if (status !== 'authenticated') return <Navigate to="/login" replace state={{ from: location.pathname }} />
  // Do not render even cached privileged content while authority is unresolved.
  if (account.isPending || account.isFetching) return <LoadingState label="Confirming account access…" />
  if (account.isError) return <Frame className={withinWorkspace ? 'workspace-access-state' : 'auth-page'}>{!withinWorkspace && <div className="auth-top"><AppBrand /><ThemeControl /></div>}<SurfaceCard className="auth-card">
    <h1>Unable to confirm access</h1><p role="alert">{apiErrorMessage(account.error)}</p>
    <AppButton onClick={() => { void account.refetch() }}>Try again</AppButton>
  </SurfaceCard></Frame>
  if (!account.data?.is_active || account.data.id !== current?.id || account.data.role !== current?.role || !roles.includes(account.data.role)) return <AccessDenied embedded={withinWorkspace} />
  return <Outlet />
}

export function EntryRoute() {
  const status = useAuthStore(s => s.status)
  const role = useAuthStore(s => s.user?.role)
  return <Navigate to={status === 'authenticated' && role ? homeForRole(role) : '/login'} replace />
}

export function RouteFailure() {
  return <main className="auth-page"><header className="auth-top"><AppBrand /><ThemeControl /></header><SurfaceCard className="auth-card">
    <h1>Unable to open this page</h1><p role="alert">Check your connection and reload to try again.</p>
    <AppButton onClick={() => window.location.reload()}>Reload page</AppButton>
  </SurfaceCard></main>
}
