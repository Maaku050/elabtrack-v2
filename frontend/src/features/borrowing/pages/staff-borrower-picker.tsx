import { useState } from 'react'
import { useAccounts } from '@/features/accounts/hooks/use-accounts'
import { DirectorySearch } from '@/components/application/directory'
import { ServerPagination } from '@/components/application/server-pagination'
import { AppButton } from '@/components/application/visual'
import { ManagementCard } from '@/components/application/management'
import { ErrorState, LoadingState } from '@/components/feedback/states'
import { apiErrorMessage } from '@/lib/api-error'
export function StaffBorrowerPicker({ selected, onSelect }: { selected: string; onSelect: (id: string, name: string) => void }) {
 const [search, setSearch] = useState(''), [page, setPage] = useState(1)
 const query = useAccounts({ page, per_page: 25, search, status: 'ACTIVE', borrower_type: '' }, false)
 return <ManagementCard title="Select an existing Borrower" description="Student or Faculty. The borrower must be active, activated and have accepted current terms."><DirectorySearch label="Search issuance borrower" placeholder="Name, email or Student ID" value={search} onChange={e => { setSearch(e.target.value); setPage(1) }} />{query.isPending ? <LoadingState /> : query.isError ? <ErrorState description={apiErrorMessage(query.error)} onRetry={() => { void query.refetch() }} /> : <><div className="divide-y">{query.data.items.map(v => <div key={v.id} className="flex flex-wrap items-center justify-between gap-3 py-3"><div><strong>{v.name}</strong><p className="management-note">{v.borrower_type}{v.student_id && ` · ${v.student_id}`} · {v.activation_required ? 'Activation pending' : 'Activation completed'}</p></div><AppButton variant="outline" disabled={v.activation_required || !v.is_active || v.id === selected} onClick={() => onSelect(v.id, v.name)} aria-label={`Select borrower ${v.name}`}>{v.id === selected ? 'Selected' : 'Select Borrower'}</AppButton></div>)}{query.data.items.length === 0 && <p>No active Borrowers match.</p>}</div><ServerPagination label="Issuance borrowers" page={page} perPage={query.data.per_page} total={query.data.total} totalPages={Math.ceil(query.data.total / query.data.per_page)} count={query.data.items.length} pending={query.isFetching} onPage={setPage} /></>}</ManagementCard>
}
