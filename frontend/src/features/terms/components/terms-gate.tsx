import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { LoadingState } from '@/components/feedback/states'
import { SurfaceCard } from '@/components/application/visual'
import { useTermsStatus } from '../hooks/use-terms'
import { TermsFrame, TermsReadError } from './terms-frame'

// Presentation gate only. Future mutations must also use server policy inside
// their transaction; read-only account/existing obligations remain reachable.
export function TermsGate() {
 const status = useTermsStatus()
 const location = useLocation()
 if (status.isPending || status.isFetching) return <TermsFrame title="Borrowing terms"><SurfaceCard className="terms-card"><LoadingState label="Checking borrowing terms…" /></SurfaceCard></TermsFrame>
 if (status.isError) return <TermsReadError error={status.error} retry={() => { void status.refetch() }} />
 if (status.data?.state !== 'accepted' || !status.data.can_initiate_borrowing || !status.data.current_terms || status.data.acceptance?.terms_version_id !== status.data.current_terms.id) return <Navigate to="/borrower/terms" replace state={{ from: location.pathname }} />
 return <Outlet />
}
