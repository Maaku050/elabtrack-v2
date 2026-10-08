import { Link } from 'react-router-dom'
import { Users, ShieldCheck, FileText } from 'lucide-react'
import { OperationalShell } from '@/components/application/operational-shell'
import { AppButton, PageHeading, SurfaceCard } from '@/components/application/visual'
import { useCurrentTerms } from '@/features/terms/hooks/use-terms'
import { apiErrorMessage } from '@/lib/api-error'
import { EmbeddedStaffDirectory } from './account-directory'
export function AdministrationPage() {
 const terms = useCurrentTerms()
 return <OperationalShell active="Administration"><PageHeading title="Administration" description="Account authority, current terms and account audit history." />
 <div className="administration-grid"><SurfaceCard className="management-card"><h2>Staff / Admin accounts</h2><div className="administration-icon"><ShieldCheck size={26} aria-hidden="true" /></div><p>Manage named Staff accounts. Existing Admin accounts are read-only; new Admin provisioning needs an approved security policy.</p><AppButton variant="outline" nativeButton={false} render={<Link to="/admin/administration/accounts" />}>Manage Accounts</AppButton></SurfaceCard>
 <SurfaceCard className="management-card"><h2>Terms & policy</h2><div className="administration-icon"><FileText size={26} aria-hidden="true" /></div>{terms.isPending ? <p role="status">Loading current publication…</p> : terms.isError ? <p>{terms.error.code === 'TERMS_NOT_PUBLISHED' ? 'Official FSMO terms await institutional approval and publication after presentation.' : apiErrorMessage(terms.error)}</p> : <><p>Current published version: {terms.data.version}</p><p>{terms.data.title}</p></>}<p className="management-note">Publication and acceptance are separate from authentication. No terms are automatically accepted.</p></SurfaceCard>
 <SurfaceCard className="management-card"><h2>Borrower accounts & audit</h2><div className="administration-icon"><Users size={26} aria-hidden="true" /></div><p>Create Student and Faculty Borrowers. Review immutable actor, time and action history in each account’s details.</p><div className="management-actions"><AppButton variant="outline" nativeButton={false} render={<Link to="/staff/borrowers" />}>Borrower Directory</AppButton><AppButton variant="outline" nativeButton={false} render={<Link to="/admin/borrowers/bulk" />}>Student Bulk Management</AppButton></div></SurfaceCard></div>
 <SurfaceCard className="management-card"><h2>Staff / Admin account directory</h2><EmbeddedStaffDirectory /></SurfaceCard></OperationalShell>
}
