import { z } from 'zod'
import { manilaInstant } from '../lib/time'
export const selectionSchema = z.object({ items: z.array(z.object({ equipment_id: z.uuid(), quantity: z.number().int().min(1).max(2147483647) })).min(1).max(100) }).refine(v => new Set(v.items.map(i => i.equipment_id)).size === v.items.length, 'Select each equipment type once.')
export const checkoutSchema = z.object({ due_local: z.string().refine(v => { const instant = manilaInstant(v); return !!instant && new Date(instant).getTime() > Date.now() }, 'Select a future due date and time in Manila.'), handover: z.boolean().refine(v => v, 'Confirm physical equipment handover.') })
export const denialSchema = z.object({ reason: z.string().trim().min(1, 'Explain the denial.').max(1000) })
