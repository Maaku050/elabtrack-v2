import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useAuthStore } from '@/stores/auth-store'
import { inventoryApi } from '@/features/inventory/api/inventory.api'
import { ApiRequestError } from '@/lib/api-error'
import { useCartStore, validQuantity } from '../stores/cart-store'
import { selectionSchema } from '../schemas/borrowing.schema'
export function useCartReview() {
 const client = useQueryClient()
 return useMutation({ retry: false, mutationFn: async () => {
  const actor = useAuthStore.getState(), lines = useCartStore.getState().lines
  const input = selectionSchema.parse({ items: lines.map(v => ({ equipment_id: v.equipment.id, quantity: Number(v.quantity) })) })
  const records = []
  // Explicit review refreshes selected equipment, including selections from other server pages.
  // Bound parallel requests rather than refetching every catalog item on navigation.
  for (let offset = 0; offset < lines.length; offset += 10) {
   records.push(...await Promise.all(lines.slice(offset, offset + 10).map(v => client.fetchQuery({ queryKey: ['inventory', actor.user?.id, 'detail', v.equipment.id], queryFn: ({ signal }) => inventoryApi.detail(v.equipment.id, signal), staleTime: 0, meta: { authenticated: true } }))))
  }
  const next = useAuthStore.getState()
  if (actor.generation !== next.generation || actor.user?.id !== next.user?.id || next.user?.role !== 'BORROWER') throw new ApiRequestError('Session changed.', 401, 'SESSION_CHANGED')
  if (JSON.stringify(lines.map(v=>[v.equipment.id,v.quantity])) !== JSON.stringify(useCartStore.getState().lines.map(v=>[v.equipment.id,v.quantity]))) throw new ApiRequestError('Cart changed. Review your current selection.',409,'CART_CHANGED')
  useCartStore.getState().refresh(records)
  if (records.some(e => e.status !== 'ACTIVE' || !validQuantity(lines.find(v => v.equipment.id === e.id)!.quantity, e.stock.available))) throw new ApiRequestError('Availability changed. Update the highlighted quantities or remove unavailable equipment before review.', 409, 'CART_AVAILABILITY_CHANGED')
  const review = { input: { items: input.items, confirm: true }, key: crypto.randomUUID() }
  useCartStore.getState().setReview(review)
  return review
 } })
}
