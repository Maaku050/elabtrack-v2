import { useId, useState } from 'react'
import { CircleCheck, CircleAlert, Ban, Minus, Plus, Package } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { AppButton, SurfaceCard, ToneBadge } from './visual'

export interface EquipmentPresentation {
  id: string
  name: string
  category: string
  unitLabel: string
  tags: readonly string[]
  image?: string
  available: number
  /** Explicit supplied view state; this component invents no low-stock policy. */
  availability: 'available' | 'attention' | 'unavailable'
}

export function EquipmentThumbnail({ src, name }: { src?: string; name: string }) {
  const [failed, setFailed] = useState(false)
  return <div className="equipment-thumbnail">{src && !failed
    ? <img src={src} alt={name} loading="lazy" onError={() => setFailed(true)} />
    : <Package aria-label={`${name}: image unavailable`} size={28} />}</div>
}

export function AvailabilityBadge({ available, state }: { available: number; state: EquipmentPresentation['availability'] }) {
  return <ToneBadge tone={state === 'unavailable' ? 'danger' : state === 'attention' ? 'warning' : 'success'} icon={state === 'unavailable' ? Ban : state === 'attention' ? CircleAlert : CircleCheck}>
    {state === 'unavailable' ? 'Unavailable' : `${state === 'attention' ? 'Low stock' : 'Available'}: ${available}`}
  </ToneBadge>
}

export function QuantitySelector({ name, value, max, onChange }: { name: string; value: number; max: number; onChange: (value: number) => void }) {
  const id = useId()
  const [error, setError] = useState('')
  function change(next: number) {
    if (!Number.isSafeInteger(next) || next < 0 || next > max) {
      setError(`Enter a whole number from 0 to ${max}.`)
      return
    }
    setError('')
    onChange(next)
  }
  return <div className="quantity-field"><div className="quantity-selector">
    <AppButton variant="ghost" aria-label={`Remove one ${name}`} disabled={value === 0} onClick={() => change(value - 1)}><Minus size={18} aria-hidden="true" /></AppButton>
    <label htmlFor={id} className="sr-only">Quantity for {name}</label>
    <Input id={id} type="number" inputMode="numeric" min={0} max={max} step={1} value={value} aria-invalid={!!error} aria-describedby={`${id}-limit${error ? ` ${id}-error` : ''}`}
      onChange={e => change(e.target.value === '' ? Number.NaN : Number(e.target.value))} />
    <AppButton variant="ghost" aria-label={`Add one ${name}`} disabled={value === max} onClick={() => change(value + 1)}><Plus size={18} aria-hidden="true" /></AppButton>
  </div><span id={`${id}-limit`} className="sr-only">Up to {max} available in this preview. Zero removes this selection.</span>
    {error && <p className="field-error" id={`${id}-error`} role="alert">{error}</p>}
  </div>
}

export function EquipmentCard({ item, selected, onSelect }: { item: EquipmentPresentation; selected: number; onSelect: (quantity: number) => void }) {
  return <SurfaceCard className="equipment-card">
    <EquipmentThumbnail src={item.image} name={item.name} />
    <div className="equipment-copy"><h2>{item.name}</h2><p>{item.category} · {item.unitLabel}</p><div className="equipment-tags">{item.tags.map(tag => <span key={tag}>{tag}</span>)}</div></div>
    <div className="equipment-controls"><AvailabilityBadge available={item.available} state={item.availability} />
      {selected > 0 ? <QuantitySelector name={item.name} value={selected} max={item.available} onChange={onSelect} />
        : <AppButton disabled={item.availability === 'unavailable'} onClick={() => onSelect(1)}>
          {item.availability !== 'unavailable' && <Plus size={18} aria-hidden="true" />}{item.availability === 'unavailable' ? 'Out of stock' : 'Add to Request'}
        </AppButton>}
    </div>
  </SurfaceCard>
}
