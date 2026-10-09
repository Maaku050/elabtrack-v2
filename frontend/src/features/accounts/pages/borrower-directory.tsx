import { useDeferredValue } from 'react'
import { Link } from 'react-router-dom'
import { Plus, Upload, ArrowUpRight, Check, Clock } from 'lucide-react'
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table'
import { DirectoryHeader, DirectoryPanel, DirectorySearch, DirectorySelect, DirectoryTableRegion } from '@/components/application/directory'
import { ServerPagination } from '@/components/application/server-pagination'
import { AppButton, ToneBadge } from '@/components/application/visual'
import { useDirectoryState } from '@/hooks/use-directory-state'
import { useAuthStore } from '@/stores/auth-store'
import { apiErrorMessage } from '@/lib/api-error'
import { useAccounts } from '../hooks/use-accounts'
import { activationLabel } from '../lib/account-status'

export function BorrowerDirectory() {
  const admin = useAuthStore(s => s.user?.role === 'ADMIN')
  const { value, page, update } = useDirectoryState()
  const currentPage = page(), search = value('search'), deferred = useDeferredValue(search)
  const type = ['STUDENT', 'FACULTY'].includes(value('borrower_type')) ? value('borrower_type') : ''
  const status = ['ACTIVE', 'INACTIVE', 'PENDING'].includes(value('status')) ? value('status') : ''
  const query = useAccounts({ page: currentPage, per_page: 25, search: deferred, borrower_type: type, status }, false)
  const data = query.data
  return <div className="directory-page"><DirectoryHeader title="Borrowers" description="Student and Faculty accounts, activation and account status." actions={admin && <><AppButton variant="outline" nativeButton={false} render={<Link to="/admin/borrowers/bulk" />}><Upload aria-hidden="true" />Student Bulk Management</AppButton><AppButton nativeButton={false} render={<Link to="/staff/borrowers/new" />}><Plus aria-hidden="true" />Add Borrower</AppButton></>} />
    <DirectoryPanel title="Borrower directory" summary={data && !query.isError ? `${data.total.toLocaleString()} matching accounts` : 'Browse Student and Faculty borrowers'} pending={query.isPending} error={query.isError ? apiErrorMessage(query.error) : undefined} empty={!data?.items.length} onRetry={() => { void query.refetch() }} toolbar={<><DirectorySearch label="Search borrowers" placeholder="Name, email or Student ID" value={search} onChange={e => update({ search: e.target.value }, true, true)} /><DirectorySelect label="Borrower type" value={type} onChange={e => update({ borrower_type: e.target.value })}><option value="">All types</option><option value="STUDENT">Student</option><option value="FACULTY">Faculty</option></DirectorySelect><DirectorySelect label="Status" value={status} onChange={e => update({ status: e.target.value })}><option value="">All statuses</option><option value="ACTIVE">Active</option><option value="INACTIVE">Inactive</option><option value="PENDING">Activation pending</option></DirectorySelect>{(search || type || status) && <AppButton variant="ghost" onClick={() => update({ search: '', borrower_type: '', status: '' })}>Clear filters</AppButton>}</>} footer={data && <ServerPagination label="Borrowers" page={data.page} perPage={data.per_page} total={data.total} count={data.items.length} totalPages={Math.ceil(data.total / data.per_page)} pending={query.isFetching} onPage={p => update({ page: String(p) }, false)} />}>
      <DirectoryTableRegion label="Borrower directory table"><Table><TableHeader><TableRow><TableHead scope="col">Borrower</TableHead><TableHead scope="col">Type / program</TableHead><TableHead scope="col">Status</TableHead><TableHead scope="col">Activation</TableHead><TableHead scope="col">Accountability</TableHead><TableHead scope="col"><span className="sr-only">Details</span></TableHead></TableRow></TableHeader><TableBody>{data?.items.map(account => <TableRow key={account.id}><TableCell><div className="directory-person"><span className="directory-person-avatar" aria-hidden="true">{account.name.trim().split(/\s+/).slice(0, 2).map(n => n[0]).join('').toUpperCase()}</span><div><Link className="directory-record-name" to={`/staff/borrowers/${account.id}`}>{account.name}</Link><span title={account.email}>{account.email}</span>{account.student_id && <span>Student ID: {account.student_id}</span>}</div></div></TableCell><TableCell><span className="directory-type" data-type={account.borrower_type}>{account.borrower_type === 'STUDENT' ? 'Student' : account.borrower_type === 'FACULTY' ? 'Faculty' : 'Unclassified'}</span>{account.course && <span className="directory-secondary">{account.course}</span>}</TableCell><TableCell><ToneBadge tone={account.is_active ? 'success' : 'neutral'} icon={Check}>{account.is_active ? 'Active' : 'Inactive'}</ToneBadge></TableCell><TableCell><ToneBadge tone={account.activation_required ? 'warning' : 'neutral'} icon={account.activation_required ? Clock : Check}>{activationLabel(account.activation_required)}</ToneBadge></TableCell><TableCell><span className="directory-secondary">{account.obligations.availability === 'UNAVAILABLE' ? 'Unavailable' : 'Review details'}</span></TableCell><TableCell><AppButton variant="ghost" nativeButton={false} render={<Link to={`/staff/borrowers/${account.id}`} aria-label={`Open ${account.name}`} />}>Open<ArrowUpRight aria-hidden="true" /></AppButton></TableCell></TableRow>)}</TableBody></Table></DirectoryTableRegion>
    </DirectoryPanel>
  </div>
}
