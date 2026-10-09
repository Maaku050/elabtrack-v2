import type { Stock } from '../types'
import type { AdjustmentValues } from '../schemas/inventory.schema'
// Display only. The backend validates and commits the authoritative stock vector.
export function stockProjection(stock: Stock, values: Pick<AdjustmentValues, 'kind' | 'quantity'>) {
 const q = values.quantity
 if (!Number.isInteger(q) || q < 0 || q > 2147483647 || (values.kind !== 'RECONCILE' && q === 0)) return null
 const delta = values.kind === 'RECONCILE' ? q - stock.available : values.kind === 'REMOVE' ? -q : q
 const available = stock.available + delta, total = stock.total_tracked + delta
 if (available < 0 || total > 2147483647) return null
 return { delta, available, total }
}
export function signedQuantity(value: number) { return `${value > 0 ? '+' : ''}${value}` }
