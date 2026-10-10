import { Link } from 'react-router-dom'
import { Package, ChevronRight } from 'lucide-react'
import { AppButton, BorrowingStatusBadge, ToneBadge } from '@/components/application/visual'
import { RecordReference } from '@/components/application/record-reference'
import { useEquipmentDetail } from '@/features/inventory/hooks/use-inventory'
import { EquipmentImage } from '@/features/inventory/pages/equipment-image'
import { manilaTimeCompact } from '../lib/time'
import type { Borrowing } from '../types'
export function BorrowingThumbnail({ id }: { id?: string }) { const q = useEquipmentDetail(id ?? ''); return q.data ? <EquipmentImage record={q.data} compact /> : <span className="borrowing-thumbnail-fallback"><Package aria-hidden="true" /></span> }
export function BorrowingCard({ record, compact = false }: { record: Borrowing; compact?: boolean }) {
 const v=record,href=`/borrower/borrowings/${v.id}`, units=v.items.reduce((n,i)=>n+i.quantity,0)
 return <article className={`borrowing-summary${compact?' is-compact':''}`}><div className="borrowing-summary-top"><BorrowingThumbnail id={v.items[0]?.equipment_id}/><div><p>{v.items[0]?.name ?? 'Borrowing'}{v.items.length>1?` + ${v.items.length-1} more`:''}</p><Link to={href} className="record-link"><RecordReference value={v.reference}/></Link><small>{units} {units===1?'unit':'units'} · {v.items.length} equipment {v.items.length===1?'type':'types'}</small></div><BorrowingStatusBadge status={v.status}/></div>{v.is_overdue&&<ToneBadge tone="danger">Overdue</ToneBadge>}<div className="borrowing-summary-dates"><span>Submitted <time dateTime={v.created_at}>{manilaTimeCompact(v.created_at)}</time></span>{(v.status==='PENDING'?v.expires_at:v.due_at)&&<span>{v.status==='PENDING'?'Reservation expires':'Due'} <time>{manilaTimeCompact(v.status==='PENDING'?v.expires_at:v.due_at)}</time></span>}</div>{!compact&&<AppButton variant="ghost" role="link" nativeButton={false} render={<Link to={href} aria-label={`Open ${v.reference}`}/>}>View Details<ChevronRight aria-hidden="true"/></AppButton>}</article>
}
