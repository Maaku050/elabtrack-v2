import type { ReactNode } from 'react'
import { Minus, Plus } from 'lucide-react'
import { AppButton } from '@/components/application/visual'
import { Input } from '@/components/ui/input'
import { validQuantity } from '../stores/cart-store'
export function QuantityControl({ name, value, max, disabled, onChange, action }: { name: string; value: string; max: number; disabled?: boolean; onChange: (value: string) => void; action?: ReactNode }) {
 const number = Number(value), valid = validQuantity(value, max)
 return <div className="selection-action-row"><AppButton variant="outline" aria-label={`Decrease ${name} quantity`} disabled={disabled || !valid || number <= 1} onClick={() => onChange(String(number - 1))}><Minus size={16} aria-hidden="true" /></AppButton><Input aria-label={`${name} quantity`} type="number" inputMode="numeric" min={1} max={Math.min(max,2147483647)} step={1} value={value} aria-invalid={!valid} disabled={disabled} onChange={e => onChange(e.target.value)} /><AppButton variant="outline" aria-label={`Increase ${name} quantity`} disabled={disabled || !valid || number >= max} onClick={() => onChange(String(number + 1))}><Plus size={16} aria-hidden="true" /></AppButton>{action}</div>
}
