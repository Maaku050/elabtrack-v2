import { useRef, useState } from 'react'
import { useForm, useFieldArray } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useEquipment, useEquipmentDetail } from '@/features/inventory/hooks/use-inventory'
import { ManagementCard, FormField, FormActions, Notice, WorkflowSteps } from '@/components/application/management'
import { DirectorySearch } from '@/components/application/directory'
import { ServerPagination } from '@/components/application/server-pagination'
import { AppButton } from '@/components/application/visual'
import { Input } from '@/components/ui/input'
import { ErrorState, LoadingState } from '@/components/feedback/states'
import { apiErrorMessage } from '@/lib/api-error'
import { useBorrowingCommand } from '../hooks/use-borrowing'
import { borrowingApi } from '../api/borrowing.api'
import { manilaInstant, manilaTime } from '../lib/time'
import { StaffBorrowerPicker } from './staff-borrower-picker'
import { checkoutSchema, selectionSchema } from '../schemas/borrowing.schema'
import type { Borrowing, RequestInput } from '../types'
import type { Equipment } from '@/features/inventory/types'
import { BorrowingFrame } from './borrowing-frame'

type Selection = { items: { equipment_id: string; quantity: number }[] }
export function BorrowerRequest() { return <RequestCompose /> }
export function StaffDirectIssuance() { return <RequestCompose direct /> }
export function RequestCompose({ direct = false }: { direct?: boolean }) {
 const [params] = useSearchParams(), initial = params.get('equipment_id') ?? '', selected = useEquipmentDetail(initial)
 if (initial && selected.isPending) return <BorrowingFrame title="New Borrowing Request" description="Loading selected equipment."><LoadingState /></BorrowingFrame>
 if (initial && selected.isError) return <BorrowingFrame title="New Borrowing Request" description="Selected equipment could not be loaded."><ErrorState description={apiErrorMessage(selected.error)} onRetry={() => { void selected.refetch() }} /></BorrowingFrame>
 const prefill = selected.data?.status === 'ACTIVE' && selected.data.stock.available > 0 ? selected.data : undefined
 return <RequestComposer key={`${direct}:${initial}`} direct={direct} prefill={prefill} />
}
function RequestComposer({ direct, prefill }: { direct: boolean; prefill?: Equipment }) {
 const [search, setSearch] = useState(''), [page, setPage] = useState(1), [names, setNames] = useState<Record<string, string>>(() => prefill ? { [prefill.id]: prefill.name } : {}), [review, setReview] = useState<{ input: RequestInput; key: string } | null>(null)
 const [borrower, setBorrower] = useState({ id: '', name: '' })
 const due = useForm({ resolver: zodResolver(checkoutSchema), defaultValues: { due_local: '', handover: false } })
 const navigate = useNavigate(), busy = useRef(false)
 const form = useForm<Selection>({ resolver: zodResolver(selectionSchema), defaultValues: { items: prefill ? [{ equipment_id: prefill.id, quantity: 1 }] : [] } }), fields = useFieldArray({ control: form.control, name: 'items' })
 const query = useEquipment({ page, per_page: 25, search, status: 'ACTIVE', available_only: true, sort: 'name' })
 const mutation = useBorrowingCommand<{ input: RequestInput; key: string }, Borrowing>(v => direct ? borrowingApi.direct(v.input, v.key) : borrowingApi.submit(v.input, v.key))
 function add(v: Equipment) { if (form.getValues('items').some(i => i.equipment_id === v.id)) return; setNames(n => ({ ...n, [v.id]: v.name })); fields.append({ equipment_id: v.id, quantity: 1 }) }
 async function confirm() { if (!review || busy.current) return; busy.current = true; try { const data = await mutation.mutateAsync(review); navigate(`${direct ? '/staff/requests' : '/borrower/borrowings'}/${data.id}`, { replace: true, state: { borrowingNotice: direct ? 'Physical issuance recorded. Equipment is now in checked-out custody.' : 'Request submitted. Stock is reserved for up to 24 hours.' } }) } catch { /* Stable reviewed key remains available for a safe retry. */ } finally { busy.current = false } }
 return <BorrowingFrame title={direct ? "Direct Equipment Issuance" : "New Borrowing Request"} description="Select equipment and quantities, then review your request." back={{ to: direct ? '/staff/requests' : '/borrower/equipment', label: direct ? 'Back to borrowings' : 'Back to catalog' }}><WorkflowSteps steps={['Select equipment', 'Review', direct ? 'Physical issuance' : 'Pending reservation']} current={review ? 1 : 0} />{mutation.error && <Notice tone="danger">{apiErrorMessage(mutation.error)} Refresh availability before changing a conflicting request.</Notice>}
 {review ? <ManagementCard title={direct ? "Review direct issuance" : "Review borrowing request"} description={direct ? "Confirm the selected borrower and actual physical handover." : "Submitting reserves these quantities immediately. Staff must still confirm physical handover."}>{direct && <Notice>Borrower: {borrower.name}. Due {manilaTime(review.input.due_at)}. Physical handover confirmed.</Notice>}<ul className="space-y-3">{review.input.items.map(i => <li key={i.equipment_id}><strong>{names[i.equipment_id] ?? 'Equipment'}</strong> · {i.quantity} units</li>)}</ul><Notice>{direct ? 'Available quantities move directly to checked-out custody. No pending reservation is created.' : 'Pending requests expire 24 hours after submission. Cancellation, denial or expiration releases reserved stock.'}</Notice><FormActions><AppButton variant="outline" disabled={mutation.isPending} onClick={() => { mutation.reset(); setReview(null) }}>Edit Request</AppButton><AppButton disabled={mutation.isPending} onClick={() => { void confirm() }}>{mutation.isPending ? 'Recording…' : direct ? 'Confirm Direct Issuance' : 'Confirm Request'}</AppButton></FormActions></ManagementCard> : <>
 {direct && <StaffBorrowerPicker selected={borrower.id} onSelect={(id, name) => setBorrower({ id, name })} />}<ManagementCard title="Choose equipment" description="Only active equipment with available units is listed."><DirectorySearch label="Search equipment to request" placeholder="Equipment name" value={search} onChange={e => { setSearch(e.target.value); setPage(1) }} />{query.isPending ? <LoadingState /> : query.isError ? <ErrorState description={apiErrorMessage(query.error)} onRetry={() => { void query.refetch() }} /> : <><div className="divide-y">{query.data.items.map(v => <div key={v.id} className="flex flex-wrap items-center justify-between gap-3 py-3"><div><strong>{v.name}</strong><p className="management-note">{v.stock.available} available · {v.category_name || 'Uncategorized'}</p></div><AppButton variant="outline" disabled={fields.fields.some(i => i.equipment_id === v.id)} onClick={() => add(v)} aria-label={`Select ${v.name}`}>{fields.fields.some(i => i.equipment_id === v.id) ? 'Selected' : 'Select'}</AppButton></div>)}{query.data.items.length === 0 && <p>No equipment matches.</p>}</div><ServerPagination label="Request equipment" page={page} perPage={query.data.per_page} total={query.data.total} totalPages={Math.ceil(query.data.total / query.data.per_page)} count={query.data.items.length} pending={query.isFetching} onPage={setPage} /></>}</ManagementCard>
 <ManagementCard title="Your selection" description="Selection is local intent. Availability is confirmed by the server when you submit."><form onSubmit={form.handleSubmit(async v => { if (direct && (!borrower.id || !await due.trigger())) return; mutation.reset(); setReview({ input: { items: v.items, confirm: true, ...(direct ? { borrower_id: borrower.id, due_at: manilaInstant(due.getValues('due_local')), physical_handover_confirmed: due.getValues('handover') } : {}) }, key: crypto.randomUUID() }) })} noValidate>{fields.fields.map((i, index) => <div key={i.id} className="flex flex-wrap items-end gap-4 py-3"><FormField id={`request-quantity-${index}`} label={`${names[i.equipment_id] ?? 'Equipment'} quantity`} required error={form.formState.errors.items?.[index]?.quantity?.message}><Input type="number" min={1} max={2147483647} {...form.register(`items.${index}.quantity`, { valueAsNumber: true })} /></FormField><AppButton variant="ghost" onClick={() => fields.remove(index)} aria-label={`Remove ${names[i.equipment_id] ?? 'equipment'}`}>Remove</AppButton></div>)}{!fields.fields.length && <p>Select at least one equipment type above.</p>}{form.formState.errors.items && <Notice tone="danger">Select distinct equipment and positive whole quantities.</Notice>}<>{direct && <><FormField id="direct-due" label="Due date and time (Asia/Manila)" required hint="Select a future timestamp in Manila time." error={due.formState.errors.due_local?.message}><Input type="datetime-local" {...due.register('due_local')} /></FormField><label className="management-checkbox"><input type="checkbox" {...due.register('handover')} />I confirm this equipment is physically being handed to the borrower.</label>{due.formState.errors.handover && <Notice tone="danger">{due.formState.errors.handover.message}</Notice>}{!borrower.id && <Notice>Select an existing eligible borrower.</Notice>}</>}</><FormActions><AppButton type="submit" disabled={!fields.fields.length || (direct && !borrower.id)}>{direct ? 'Review Direct Issuance' : 'Review Request'}</AppButton></FormActions></form></ManagementCard></>}
 </BorrowingFrame>
}
