import { z } from 'zod'
const quantity = z.number().int('Enter a whole number.').min(0).max(2147483647)
export const equipmentSchema = z.object({ name: z.string().trim().min(1, 'Enter the equipment name.').max(160), description: z.string().max(2000), category_id: z.string(), opening_quantity: quantity, reason: z.string().trim().max(1000) }).refine(v => v.opening_quantity === 0 || v.reason.length > 0, { path: ['reason'], message: 'Explain the opening quantity.' })
export type EquipmentValues = z.infer<typeof equipmentSchema>
export const adjustmentSchema = z.object({ kind: z.enum(['ADD', 'REMOVE', 'RECONCILE']), quantity, reason: z.string().trim().min(1, 'Explain this change.').max(1000) }).refine(v => v.kind === 'RECONCILE' || v.quantity > 0, { path: ['quantity'], message: 'Enter a positive whole quantity.' })
export type AdjustmentValues = z.infer<typeof adjustmentSchema>
