import { useDeferredValue, useState } from 'react'
import { Link } from 'react-router-dom'
import { Plus, Upload, ChevronRight } from 'lucide-react'
import { OperationalShell } from '@/components/application/operational-shell'
import { AppButton, PageHeading, ToneBadge, SurfaceCard } from '@/components/application/visual'
import { DataTableShell, SearchField } from '@/components/application/data-table'
import { LoadingState, EmptyState } from '@/components/feedback/states'
import { useAuthStore } from '@/stores/auth-store'
import { apiErrorMessage } from '@/lib/api-error'
import { useAccounts } from '../hooks/use-accounts'
export function AccountDirectory({ staff = false, embedded = false }: { staff?: boolean; embedded?: boolean }) {
 const admin = useAuthStore(s => s.user?.role === 'ADMIN'); const [search, setSearch] = useState(''); const deferred = useDeferredValue(search)
 const [page, setPage] = useState(1); const [type, setType] = useState(''); const [status, setStatus] = useState('')
 const query = useAccounts({ page, per_page: 25, search: deferred, borrower_type: type, status }, staff)
 const base = staff ? '/admin/administration/accounts' : '/staff/borrowers'
 const content = <>{!embedded && <PageHeading title={staff ? 'Staff & Admin Accounts' : 'Borrowers'} description={staff ? 'Restricted account administration. New Admin provisioning requires a separate approved security policy.' : 'Manage borrower accounts and operational information.'} action={admin && <div className="management-actions"><AppButton  nativeButton={false} render={<Link to={`${base}/new`} />}> <Plus size={18} />{staff ? 'Add Staff' : 'Add Borrower'}</AppButton>{!staff && <AppButton  variant="outline" nativeButton={false} render={<Link to="/admin/borrowers/bulk" />}> <Upload size={18} />Student Bulk Management</AppButton>}</div>} />}
 <div className="management-toolbar"><SearchField label="Search accounts" placeholder="Search name, email or Student ID…" value={search} onChange={value => { setSearch(value); setPage(1) }} />{!staff && <label>Borrower type<select value={type} onChange={e => { setType(e.target.value); setPage(1) }}><option value="">All types</option><option>STUDENT</option><option>FACULTY</option></select></label>}<label>Status<select value={status} onChange={e => { setStatus(e.target.value); setPage(1) }}><option value="">All statuses</option><option value="ACTIVE">Active</option><option value="INACTIVE">Inactive</option><option value="PENDING">Activation pending</option></select></label></div>
 {query.isPending ? <LoadingState label="Loading accounts…" /> : query.isError ? <SurfaceCard className="management-card"><p role="alert">{apiErrorMessage(query.error)}</p><AppButton variant="outline" onClick={() => { void query.refetch() }}>Retry</AppButton></SurfaceCard> : <DataTableShell label="Account directory" footer={<div className="pagination-controls"><p aria-live="polite">{query.data.total} accounts · Page {page} of {Math.max(1, Math.ceil(query.data.total / 25))}</p><nav aria-label="Account pagination"><AppButton variant="outline" disabled={page === 1} onClick={() => setPage(page - 1)}>Previous</AppButton><AppButton variant="outline" disabled={page * 25 >= query.data.total} onClick={() => setPage(page + 1)}>Next</AppButton></nav></div>}>
 {query.data.items.length === 0 ? <EmptyState title="No accounts found" description={search || type || status ? 'Try adjusting your search or filters.' : 'Provisioned accounts will appear here.'} /> : <table className="management-table"><thead><tr><th>{staff ? 'Account' : 'Borrower'}</th><th>{staff ? 'Role' : 'Type / program'}</th><th>Status</th><th>Activation</th>{!staff && <th>Accountability</th>}<th><span className="sr-only">Details</span></th></tr></thead><tbody>{query.data.items.map(account => <tr key={account.id}><td><strong>{account.name}</strong><small>{account.email}</small>{account.student_id && <small>Student ID: {account.student_id}</small>}</td><td>{staff ? account.role : account.borrower_type || 'Unclassified'}{account.course && <small>{account.course}</small>}</td><td><ToneBadge tone={account.is_active ? 'success' : 'neutral'}>{account.is_active ? 'Active' : 'Inactive'}</ToneBadge></td><td>{account.activation_required ? 'Pending' : 'Established'}</td>{!staff && <td><ToneBadge tone="neutral">{account.obligations.availability === 'UNAVAILABLE' ? 'Unavailable' : 'Review details'}</ToneBadge></td>}<td><AppButton  variant="outline" nativeButton={false} render={<Link to={`${base}/${account.id}`} aria-label={`Open ${account.name}`} />}> Open<ChevronRight size={16} /></AppButton></td></tr>)}</tbody></table>}
 </DataTableShell>}
 </>
 return embedded ? content : <OperationalShell active={staff ? 'Administration' : 'Borrowers'}>{content}</OperationalShell>
}
export function StaffDirectory() { return <AccountDirectory staff /> }

export function EmbeddedStaffDirectory() { return <AccountDirectory staff embedded /> }
