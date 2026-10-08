import { useEffect, useRef } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Package } from 'lucide-react'
import { useAuthStore } from '@/stores/auth-store'
import { inventoryApi } from '../api/inventory.api'
import type { Equipment } from '../types'
export function EquipmentImage({ record }: { record: Equipment }) {
 const actor = useAuthStore(s => s.user?.id); const img = useRef<HTMLImageElement>(null)
 const query = useQuery({ queryKey: ['inventory', actor, 'image', record.id, record.image_id], queryFn: ({ signal }) => inventoryApi.image(record, signal), enabled: !!actor && !!record.image_id, meta: { authenticated: true } })
 useEffect(() => { if (!query.data || !img.current) return; const url = URL.createObjectURL(query.data); img.current.src = url; return () => URL.revokeObjectURL(url) }, [query.data])
 return <div className="equipment-thumbnail">{query.data && !query.isError ? <img ref={img} alt={record.name} /> : <Package size={28} aria-label={`${record.name}: image unavailable`} />}</div>
}
