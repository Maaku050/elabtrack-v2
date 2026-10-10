import { useMutation, useQueryClient } from '@tanstack/react-query'
import { inventoryApi } from '@/features/inventory/api/inventory.api'
import type { Equipment } from '@/features/inventory/types'
import { ApiRequestError } from '@/lib/api-error'
import { useAuthStore } from '@/stores/auth-store'
import { borrowingApi } from '../api/borrowing.api'
import { checkoutSchema, selectionSchema } from '../schemas/borrowing.schema'
import { validQuantity } from '../stores/cart-store'
import { manilaInstant } from '../lib/time'
export type IssuanceLine = { equipment: Equipment; quantity: string }
export function useBorrowerEligibility() {
 return useMutation({ retry: false, mutationFn: async (id: string) => {
  const generation = useAuthStore.getState().generation
  const result = await borrowingApi.eligibility(id)
  if (generation !== useAuthStore.getState().generation) throw new ApiRequestError('Session changed.',401,'SESSION_CHANGED')
  return result
 } })
}
/** A draft preview. The final direct command always rechecks eligibility and custody atomically. */
export function useDirectReview() {
 const client = useQueryClient()
 return useMutation({ retry: false, mutationFn: async ({ borrower, lines, due, handover }: { borrower: string; lines: IssuanceLine[]; due: string; handover: boolean }) => {
  const actor = useAuthStore.getState()
  const checkout = checkoutSchema.safeParse({ due_local: due, handover })
  if (!checkout.success) throw new ApiRequestError(checkout.error.issues[0]?.message || 'Check the due date and physical handover confirmation.',400,'VALIDATION_ERROR')
  const selection = selectionSchema.parse({ items: lines.map(v => ({ equipment_id:v.equipment.id, quantity:Number(v.quantity) })) })
  await borrowingApi.eligibility(borrower)
  const records: Equipment[] = []
  for (let offset=0; offset<lines.length; offset+=10) records.push(...await Promise.all(lines.slice(offset,offset+10).map(v=>client.fetchQuery({ queryKey:['inventory',actor.user?.id,'detail',v.equipment.id], queryFn:({signal})=>inventoryApi.detail(v.equipment.id,signal), staleTime:0, meta:{authenticated:true} }))))
  const next = useAuthStore.getState()
  if (actor.generation!==next.generation || actor.user?.id!==next.user?.id || !['STAFF','ADMIN'].includes(next.user?.role||'')) throw new ApiRequestError('Session changed.',401,'SESSION_CHANGED')
  return { records, valid: records.every(v=>v.status==='ACTIVE'&&validQuantity(lines.find(i=>i.equipment.id===v.id)!.quantity,v.stock.available)), review: { input:{ items:selection.items, borrower_id:borrower, due_at:manilaInstant(checkout.data.due_local), physical_handover_confirmed:checkout.data.handover, confirm:true }, key:crypto.randomUUID() } }
 } })
}
