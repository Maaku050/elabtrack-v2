import { describe, it, expect } from 'vitest'
import { signedQuantity, stockProjection } from './stock-projection'
const stock = { available: 15, reserved: 2, checked_out: 3, damaged_held: 4, total_tracked: 24 }
describe('stock display projections', () => {
 it.each([['ADD', 5, 20, 29, 5], ['REMOVE', 5, 10, 19, -5], ['RECONCILE', 20, 20, 29, 5], ['RECONCILE', 0, 0, 9, -15]] as const)('projects %s while preserving unavailable stock', (kind, quantity, available, total, delta) => { expect(stockProjection(stock, { kind, quantity })).toEqual({ available, total, delta }); expect(stock.reserved + stock.checked_out + stock.damaged_held).toBe(9) })
 it.each([['REMOVE', 16], ['ADD', 2147483647], ['ADD', 0], ['RECONCILE', -1], ['ADD', 1.5], ['RECONCILE', NaN]] as const)('does not display a valid result for %s %s', (kind, quantity) => { expect(stockProjection(stock, { kind, quantity })).toBeNull() })
 it('shows the quantity direction', () => { expect(signedQuantity(5)).toBe('+5'); expect(signedQuantity(-5)).toBe('-5'); expect(signedQuantity(0)).toBe('0') })
})
