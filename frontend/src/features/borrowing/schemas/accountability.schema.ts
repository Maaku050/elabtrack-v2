import { z } from 'zod'
const quantity = z.number().int().min(0).max(2147483647)
const reason = z.string().trim().min(1, 'Record the in-person verification explanation.').max(1000)
export const returnSchema = z.object({ reason, lines: z.array(z.object({ item_id: z.uuid(), good: quantity, damaged: quantity, lost: quantity })).min(1).max(100) }).refine(v => v.lines.some(l => l.good + l.damaged + l.lost > 0), 'Record at least one unit.').refine(v => new Set(v.lines.map(l => l.item_id)).size === v.lines.length, 'Record each item once.')
export const replacementSchema = z.object({ obligation_id: z.uuid(), quantity: quantity.min(1), reason, equivalent_confirmed: z.boolean().refine(v => v, 'Confirm that the replacement was physically received and is equivalent.') })
export const fineSchema = z.object({ method: z.enum(['PAID', 'WAIVED', 'OTHER_RESOLUTION']), note: z.string().trim().max(1000) })
