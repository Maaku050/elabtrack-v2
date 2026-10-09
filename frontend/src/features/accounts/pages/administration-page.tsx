import { Link } from 'react-router-dom'
import { Users, ShieldCheck, FileText } from 'lucide-react'
import { AppButton } from '@/components/application/visual'
import { ManagementPage, ManagementCard, Notice } from '@/components/application/management'
import { useCurrentTerms } from '@/features/terms/hooks/use-terms'
import { apiErrorMessage } from '@/lib/api-error'
import { EmbeddedStaffDirectory } from './account-directory'
export function AdministrationPage() {
 const terms=useCurrentTerms()
 return <ManagementPage title="Administration" domain="Administration" description="Account authority, current terms and account audit history."><div className="management-overview"><ManagementCard title="Staff / Admin accounts"><ShieldCheck className="management-overview-icon" aria-hidden="true" /><p>Manage named Staff accounts. Existing Admin accounts are read-only; new Admin provisioning needs an approved security policy.</p><AppButton variant="outline" nativeButton={false} render={<Link to="/admin/administration/accounts" />}>Manage Accounts</AppButton></ManagementCard>
 <ManagementCard title="Terms & policy"><FileText className="management-overview-icon" aria-hidden="true" />{terms.isPending?<p role="status">Loading current publication…</p>:terms.isError?<Notice tone="danger">{apiErrorMessage(terms.error)}<AppButton variant="outline" onClick={()=>{void terms.refetch()}}>Retry publication status</AppButton></Notice>:terms.data===null?<Notice tone="warning">Official FSMO terms await institutional approval and publication after presentation.</Notice>:<><p>Current published version: {terms.data.version}</p><p>{terms.data.title}</p></>}<p className="management-note">Publication and acceptance are separate from authentication. No terms are automatically accepted.</p></ManagementCard>
 <ManagementCard title="Borrower accounts & audit"><Users className="management-overview-icon" aria-hidden="true" /><p>Create Student and Faculty Borrowers. Review actor, time and action history in each account’s details.</p><div className="management-actions"><AppButton variant="outline" nativeButton={false} render={<Link to="/staff/borrowers" />}>Borrower Directory</AppButton><AppButton variant="outline" nativeButton={false} render={<Link to="/admin/borrowers/bulk" />}>Student Bulk Management</AppButton></div></ManagementCard></div><EmbeddedStaffDirectory /></ManagementPage>
}
