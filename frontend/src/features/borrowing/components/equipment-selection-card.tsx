import { useState } from 'react'
import { Link } from 'react-router-dom'
import { AppButton, SurfaceCard } from '@/components/application/visual'
import { EquipmentImage } from '@/features/inventory/pages/equipment-image'
import type { Equipment } from '@/features/inventory/types'
import { QuantityControl } from './quantity-control'
import { validQuantity } from '../stores/cart-store'
export function EquipmentCard({ equipment: v, selected, selectedQuantity, locked, onAdd, onCart, staff = false }: { equipment: Equipment; selected: boolean; selectedQuantity?: string; locked: boolean; onAdd: (v: Equipment, qty: string) => void; onCart: () => void; staff?: boolean }) {
 const [quantity, setQuantity] = useState('1'), eligible = v.status === 'ACTIVE' && v.stock.available > 0
 return <SurfaceCard className="selection-equipment" data-equipment-id={v.id}><EquipmentImage record={v} /><div className="selection-equipment-copy"><Link to={`${staff?'/staff/inventory':'/borrower/equipment'}/${v.id}`}><h3>{v.name}</h3></Link><p>{v.category_name || 'Uncategorized'}</p>{v.description&&<p className="selection-equipment-description" title={v.description}>{v.description}</p>}<strong className={eligible ? 'selection-available' : 'selection-unavailable'}>{v.stock.available} available</strong></div><div className="selection-equipment-action"><QuantityControl name={v.name} value={selected ? selectedQuantity ?? quantity : quantity} max={v.stock.available} disabled={!eligible || selected || locked} onChange={setQuantity} action={<AppButton disabled={locked || (!selected && (!eligible || !validQuantity(quantity,v.stock.available)))} aria-label={selected ? `Edit ${v.name} in cart` : `Add ${v.name} to cart`} onClick={() => selected ? onCart() : onAdd(v,quantity)}>{selected ? 'In cart' : 'Add'}</AppButton>} />{!eligible && <p className="management-note">Currently unavailable for borrowing.</p>}{eligible && !selected && !validQuantity(quantity,v.stock.available) && <p className="selection-error" role="status">Enter 1–{v.stock.available} whole units.</p>}</div></SurfaceCard>
}
