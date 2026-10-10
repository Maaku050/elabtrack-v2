import { create } from 'zustand'
import { useAuthStore } from '@/stores/auth-store'
import type { Equipment } from '@/features/inventory/types'
import type { RequestInput } from '../types'

export interface CartLine { equipment: Equipment; quantity: string }
export interface CartReview { input: RequestInput; key: string }
interface CartState {
 lines: CartLine[]
 review: CartReview | null
 add: (equipment: Equipment, quantity: string) => void
 quantity: (id: string, quantity: string) => void
 remove: (id: string) => void
 refresh: (records: Equipment[]) => void
 setReview: (review: CartReview | null) => void
 clear: () => void
}
// Private, ephemeral intent. Never persisted; stock and eligibility remain server state.
export const useCartStore = create<CartState>((set) => ({
 lines: [], review: null,
 add: (equipment, quantity) => set(s => s.review || s.lines.length >= 100 || s.lines.some(v => v.equipment.id === equipment.id) ? s : { lines: [...s.lines, { equipment, quantity }] }),
 quantity: (id, quantity) => set(s => s.review ? s : { lines: s.lines.map(v => v.equipment.id === id ? { ...v, quantity } : v) }),
 remove: id => set(s => s.review ? s : { lines: s.lines.filter(v => v.equipment.id !== id) }),
 refresh: records => set(s => {
  let changed = false
  const lines = s.lines.map(v => {
   const equipment = records.find(e => e.id === v.equipment.id)
   if (!equipment || JSON.stringify(equipment) === JSON.stringify(v.equipment)) return v
   changed = true
   return { ...v, equipment }
  })
  return changed ? { lines } : s
 }),
 setReview: review => set({ review }),
 clear: () => set({ lines: [], review: null }),
}))
useAuthStore.subscribe((next, previous) => {
 if (next.user?.id !== previous.user?.id || next.generation !== previous.generation || (previous.isAuthenticated && !next.isAuthenticated)) useCartStore.getState().clear()
})
export function validQuantity(value: string, available: number): boolean {
 return /^\d+$/.test(value) && Number.isSafeInteger(Number(value)) && Number(value) >= 1 && Number(value) <= Math.min(available, 2147483647)
}
