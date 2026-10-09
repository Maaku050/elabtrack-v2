import { useRef, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useForm, useWatch } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { ArrowRight, ShieldCheck } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { AppButton, ToneBadge } from '@/components/application/visual'
import { ManagementPage, ManagementCard, FormSection, FormField, FormActions, Notice } from '@/components/application/management'
import { NativeSelect } from '@/components/ui/native-select'
import { Textarea } from '@/components/ui/textarea'
import { ErrorState, LoadingState } from '@/components/feedback/states'
import { apiErrorMessage } from '@/lib/api-error'
import { inventoryApi } from '../api/inventory.api'
import { useEquipmentDetail, useInventoryCommand } from '../hooks/use-inventory'
import { adjustmentSchema, type AdjustmentValues } from '../schemas/inventory.schema'
import { stockProjection, signedQuantity } from '../lib/stock-projection'
import type { Equipment, Adjustment, Stock } from '../types'
function AdjustmentForm({ record, reconcile, refresh }: { record: Equipment; reconcile: boolean; refresh: () => Promise<unknown> }) {
 const navigate = useNavigate(), key = useRef({ payload: '', key: crypto.randomUUID() }), submitting = useRef(false)
 const [review, setReview] = useState<{ values: AdjustmentValues; sequence: number; stock: Stock } | null>(null)
 const form = useForm<AdjustmentValues>({ resolver: zodResolver(adjustmentSchema), defaultValues: { kind: reconcile ? 'RECONCILE' : 'ADD', quantity: reconcile ? record.stock.available : 1, reason: '' } })
 const values = useWatch({ control: form.control }) as AdjustmentValues
 const command = useInventoryCommand<Adjustment, Equipment>(async value => { const payload = JSON.stringify(value); if (key.current.payload !== payload) key.current = { payload, key: crypto.randomUUID() }; return inventoryApi.adjust(record.id, value, key.current.key) })
 const basis = review?.stock ?? record.stock, proposed = review?.values ?? values, projection = stockProjection(basis, proposed)
 const stale = !!review && review.sequence !== record.stock_sequence
 async function confirm() {
  if (!review || !projection || stale || submitting.current || command.isPending) return
  submitting.current = true
  try { const saved = await command.mutateAsync({ ...review.values, expected_sequence: review.sequence, confirm: true }); navigate(`/staff/inventory/${record.id}`, { state: { stockNotice: `Stock change recorded. Available: ${saved.stock.available}. Total tracked: ${saved.stock.total_tracked}.` } }) }
  catch { await refresh() } finally { submitting.current = false }
 }
 return <div className="stock-workflow">
 <ManagementCard className="stock-identity"><div><p className="stock-eyebrow">{record.category_name || 'Uncategorized'}</p><h2>{record.name}</h2></div><ToneBadge tone={reconcile ? 'warning' : 'neutral'}>{reconcile ? 'Admin reconciliation' : 'Physical stock'}</ToneBadge><dl className="stock-metrics">{[['Available', record.stock.available], ['Reserved', record.stock.reserved], ['Checked out', record.stock.checked_out], ['Damaged held', record.stock.damaged_held], ['Total tracked', record.stock.total_tracked]].map(([label, value]) => <div key={label}><dt>{label}</dt><dd>{value}</dd></div>)}</dl></ManagementCard>
 <ManagementCard className={`stock-editor ${reconcile ? 'stock-exceptional' : ''}`}><div className="stock-editor-heading"><h2>{review ? 'Review stock change' : reconcile ? 'Verify available stock' : 'Adjust usable stock'}</h2>{reconcile && <ShieldCheck size={22} aria-hidden="true" />}</div><p className="management-note">{reconcile ? 'Exceptional reconciliation of the physically verified available count. Reserved, checked-out and damaged-held stock is preserved.' : 'Record an acquisition or remove units from available stock, with an explanation.'}</p>
 {command.error && <Notice tone="danger">{apiErrorMessage(command.error)}</Notice>}{stale && <Notice tone="warning">Stock changed after your review. Go back to edit and review the current quantities before confirming.</Notice>}
 <div className="stock-editor-layout"><div>{review ? <><h3>Change to confirm</h3><dl className="management-dl"><dt>Operation</dt><dd>{review.values.kind === 'RECONCILE' ? 'Inventory correction' : review.values.kind === 'REMOVE' ? 'Remove usable stock' : 'Add usable stock'}</dd><dt>{reconcile ? 'Verified available' : 'Quantity'}</dt><dd>{review.values.quantity}</dd><dt>Required explanation</dt><dd className="stock-explanation">{review.values.reason}</dd></dl></> : <form id="stock-change-form" onSubmit={form.handleSubmit(v => { if (!stockProjection(record.stock, v)) { form.setError('quantity', { message: 'The projected stock is outside the allowed range. Check the available quantity.' }); return } command.reset(); setReview({ values: { ...v }, sequence: record.stock_sequence, stock: { ...record.stock } }) })} noValidate><FormSection title={reconcile ? 'Verified available count' : 'Stock change'} disabled={command.isPending || record.status === 'ARCHIVED'}>{!reconcile && <FormField id="adjust-kind" label="Operation" required><NativeSelect {...form.register('kind')}><option value="ADD">Add usable stock</option><option value="REMOVE">Remove available stock</option></NativeSelect></FormField>}<FormField id="adjust-quantity" label={reconcile ? 'Verified available count' : 'Quantity'} required error={form.formState.errors.quantity?.message}><Input type="number" min={reconcile ? 0 : 1} max={2147483647} {...form.register('quantity', {valueAsNumber:true})} /></FormField><FormField id="adjust-reason" label="Required explanation" required error={form.formState.errors.reason?.message}><Textarea {...form.register('reason')} /></FormField></FormSection></form>}</div>
 <aside className="stock-projection" aria-label="Projected stock" aria-live="polite"><h3>{review ? 'Reviewed result' : 'Projected result'}</h3><dl><div><dt>Current available</dt><dd>{basis.available}</dd></div><div className="stock-effect"><dt>{reconcile ? 'Correction difference' : proposed.kind === 'REMOVE' ? 'Remove' : 'Add'}</dt><dd>{projection ? signedQuantity(projection.delta) : '—'}</dd></div><div className="stock-result"><dt>Projected available</dt><dd>{projection?.available ?? '—'}<ArrowRight size={20} aria-hidden="true" /></dd></div><div><dt>Projected total tracked</dt><dd>{projection?.total ?? '—'}</dd></div></dl>{projection && projection.delta < 0 && <p className="management-warning">This {reconcile ? 'correction' : 'removal'} decreases available physical stock by {Math.abs(projection.delta)} units.</p>}{!projection && <p className="field-error">Enter a valid quantity within available stock and the supported range.</p>}</aside></div>
 <p className="management-note">Every confirmed change records a reasoned ledger movement. It does not record a return or settle a fine or replacement obligation.</p><FormActions>{review ? <><AppButton variant="outline" disabled={command.isPending} onClick={() => { setReview(null); command.reset() }}>Back to Edit</AppButton><AppButton disabled={command.isPending || stale || !projection} onClick={() => { void confirm() }}>{command.isPending ? 'Recording…' : 'Confirm Stock Change'}</AppButton></> : <AppButton type="submit" form="stock-change-form" disabled={record.status === 'ARCHIVED'}>Review Stock Change</AppButton>}</FormActions></ManagementCard></div>
}
export function StockAdjustment({ reconcile = false }: { reconcile?: boolean }) { const { id = '' } = useParams(); const query = useEquipmentDetail(id); return <ManagementPage title={reconcile ? 'Inventory Correction' : 'Stock Adjustment'} domain="Equipment management" description={reconcile ? 'Admin-only verified-count reconciliation with an audit trail.' : 'Record verified acquisition or removal of usable stock.'} back={{to:`/staff/inventory/${id}`,label:'Back to equipment'}}>{query.isPending ? <LoadingState /> : query.isError ? <ErrorState description={apiErrorMessage(query.error)} /> : <AdjustmentForm key={id} record={query.data} reconcile={reconcile} refresh={query.refetch} />}</ManagementPage> }

export function InventoryReconciliation() { return <StockAdjustment reconcile /> }
