import { useRef, useState } from 'react'
import { useLocation, useParams } from 'react-router-dom'
import { useAuthStore } from '@/stores/auth-store'
import { AppButton, BorrowingStatusBadge, ToneBadge } from '@/components/application/visual'
import { ManagementCard, Notice } from '@/components/application/management'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { DirectoryTableRegion } from '@/components/application/directory'
import { AlertDialog, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog'
import { ErrorState, LoadingState } from '@/components/feedback/states'
import { apiErrorMessage } from '@/lib/api-error'
import { borrowingApi } from '../api/borrowing.api'
import { useBorrowing, useBorrowingCommand } from '../hooks/use-borrowing'
import type { Borrowing } from '../types'
import { manilaTime } from '../lib/time'
import { StaffActions } from './staff-actions'
import { BorrowingFrame } from './borrowing-frame'
export function BorrowingDetail() {
 const { id = '' } = useParams(), query = useBorrowing(id), borrower = useAuthStore(s => s.user?.role === 'BORROWER'), location = useLocation(), [review, setReview] = useState<string | null>(null), busy = useRef(false)
 const [notice,setNotice] = useState('')
 const cancel = useBorrowingCommand<string, Borrowing>(key => borrowingApi.decide(id, 'cancel', { confirm: true }, key)), v = query.data, base = borrower ? '/borrower/borrowings' : '/staff/requests'
 async function confirm() { if (!review || busy.current) return; busy.current = true; try { await cancel.mutateAsync(review); setReview(null) } catch { /* Keep key and error for retry; 409 revalidates authoritative state. */ } finally { busy.current = false } }
 return <BorrowingFrame title="Borrowing Details" description="Reservation, handover and retained request history." back={{ to: base, label: 'Back to borrowings' }}>{(location.state as { borrowingNotice?: string } | null)?.borrowingNotice && <Notice tone="success">{(location.state as { borrowingNotice: string }).borrowingNotice}</Notice>}{notice && <Notice tone="success">{notice}</Notice>}{query.isPending ? <LoadingState /> : query.isError ? <ErrorState description={apiErrorMessage(query.error)} onRetry={() => { void query.refetch() }} /> : v && <>
 <ManagementCard title={v.reference} actions={<BorrowingStatusBadge status={v.status} />}><dl className="management-dl"><dt>Borrower</dt><dd>{v.borrower_name} · {v.borrower_type}{v.student_id && ` · ${v.student_id}`}</dd><dt>Entry</dt><dd>{v.entry_path === 'DIRECT' ? 'Direct physical issuance' : 'Borrower request'}</dd><dt>Submitted / recorded</dt><dd>{manilaTime(v.created_at)}</dd><dt>Reservation expires</dt><dd>{manilaTime(v.expires_at)}</dd><dt>Physical handover</dt><dd>{manilaTime(v.checked_out_at)}</dd><dt>Due</dt><dd>{manilaTime(v.due_at)} {v.is_overdue && <ToneBadge tone="danger">Overdue</ToneBadge>}</dd></dl>{v.status === 'PENDING' && <Notice>Stock remains reserved until checkout, cancellation, denial or expiration is recorded.</Notice>}{v.status === 'DENIED' && <Notice tone="warning" title="Denial explanation">{v.denial_reason}</Notice>}{['EXPIRED', 'CANCELLED'].includes(v.status) && <Notice>Reserved quantities were released. This record is retained in your history.</Notice>}{borrower && v.status === 'PENDING' && <AppButton variant="outline" onClick={() => { cancel.reset(); setReview(crypto.randomUUID()) }}>Cancel Pending Request</AppButton>}</ManagementCard>
 <ManagementCard title="Equipment quantities"><DirectoryTableRegion label="Borrowing equipment"><Table><TableHeader><TableRow>{['Equipment', 'Requested', 'Reserved', 'Issued'].map(h => <TableHead key={h} scope="col">{h}</TableHead>)}</TableRow></TableHeader><TableBody>{v.items.map(i => <TableRow key={i.id}><TableCell>{i.name}</TableCell><TableCell>{i.quantity}</TableCell><TableCell>{i.reserved_quantity}</TableCell><TableCell>{i.issued_quantity}</TableCell></TableRow>)}</TableBody></Table></DirectoryTableRegion>{v.status === 'CHECKED_OUT' && <Notice>Issued equipment remains in checked-out custody. Contact FSMO for face-to-face return processing.</Notice>}</ManagementCard>
 {!borrower && v.status==='PENDING' && <StaffActions record={v} onSuccess={setNotice} />}<ManagementCard title="Request & issuance history"><ol className="space-y-4">{v.events.map(e => <li key={e.id}><strong>{e.kind.replaceAll('_', ' ')}</strong><p className="management-note">{manilaTime(e.occurred_at)} · {e.actor_id === null ? 'Automatic expiration' : 'Authenticated account action'}</p>{e.reason && <p>{e.reason}</p>}</li>)}</ol></ManagementCard>
 </>}
 <AlertDialog open={!!review} onOpenChange={open => { if (!open && !cancel.isPending) setReview(null) }}><AlertDialogContent><AlertDialogHeader><AlertDialogTitle>Cancel pending request?</AlertDialogTitle><AlertDialogDescription>Reserved units will become available again. Your request history will be retained.</AlertDialogDescription></AlertDialogHeader>{cancel.error && <Notice tone="danger">{apiErrorMessage(cancel.error)}</Notice>}<AlertDialogFooter><AppButton variant="outline" disabled={cancel.isPending} onClick={() => setReview(null)}>Keep Request</AppButton><AppButton disabled={cancel.isPending || v?.status !== 'PENDING'} onClick={() => { void confirm() }}>{cancel.isPending ? 'Cancelling…' : 'Confirm Cancellation'}</AppButton></AlertDialogFooter></AlertDialogContent></AlertDialog>
 </BorrowingFrame>
}
