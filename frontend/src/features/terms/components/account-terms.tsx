import { Link } from 'react-router-dom'
import { SurfaceCard } from '@/components/application/visual'
import { useTermsStatus } from '../hooks/use-terms'
export function AccountTerms() {
 const status = useTermsStatus()
 return <SurfaceCard className="terms-card"><h2>Borrowing terms</h2>
  <p>{status.isPending || status.isFetching ? 'Checking acceptance…' : status.isError ? 'Acceptance status is unavailable.' : status.data?.state === 'accepted' ? `Accepted Version ${status.data.current_terms?.version}.` : status.data?.state === 'unpublished' ? 'Official terms have not been published.' : status.data?.state === 'updated' ? 'Updated terms require your acceptance before new borrowing requests.' : 'Review and accept the current terms before new borrowing requests.'}</p>
  <Link className="app-button" to="/borrower/terms">View borrowing terms</Link>
 </SurfaceCard>
}
