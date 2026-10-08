import { describe, it, expect } from 'vitest'
import { adjustmentSchema, equipmentSchema } from './inventory.schema'
describe('inventory input boundaries', () => {
 it('requires a reason for opening stock, permits an empty pool', () => { const input = { name: 'Spoon', description: '', category_id: '', opening_quantity: 0, reason: '' }; expect(equipmentSchema.safeParse(input).success).toBe(true); expect(equipmentSchema.safeParse({ ...input, opening_quantity: 2 }).success).toBe(false) })
 it('rejects fractions, negatives and overflow', () => { for (const quantity of [-1, 1.2, 2147483648]) expect(adjustmentSchema.safeParse({ kind: 'ADD', quantity, reason: 'Acquisition' }).success).toBe(false) })
 it('zero is a valid observed count only in reconciliation', () => { expect(adjustmentSchema.safeParse({ kind: 'RECONCILE', quantity: 0, reason: 'Physical count' }).success).toBe(true); expect(adjustmentSchema.safeParse({ kind: 'REMOVE', quantity: 0, reason: 'Physical count' }).success).toBe(false) })
})
