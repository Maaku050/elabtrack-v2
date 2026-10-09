import { useEffect, useRef, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Package } from 'lucide-react'
import { useAuthStore } from '@/stores/auth-store'
import { inventoryApi } from '../api/inventory.api'
import type { Equipment } from '../types'
export function EquipmentImage({ record, compact = false }: { record: Equipment; compact?: boolean }) {
 const actor = useAuthStore(s => s.user?.id), host = useRef<HTMLDivElement>(null), img = useRef<HTMLImageElement>(null)
 const [visible, setVisible] = useState(() => typeof IntersectionObserver === 'undefined'), [failedImage, setFailedImage] = useState<string | null>(null)
 const query = useQuery({ queryKey: ['inventory', actor, 'image', record.id, record.image_id], queryFn: ({ signal }) => inventoryApi.image(record, signal), enabled: !!actor && !!record.image_id && visible, staleTime: 5 * 60_000, retry: false, meta: { authenticated: true } })
 const failed = failedImage === record.image_id && !!record.image_id
 useEffect(() => { const target = host.current; if (!target) return; if (!('IntersectionObserver' in window)) return; const observer = new IntersectionObserver(entries => { if (entries.some(e => e.isIntersecting)) { setVisible(true); observer.disconnect() } }, { rootMargin: '100px' }); observer.observe(target); return () => observer.disconnect() }, [])
 useEffect(() => { if (!query.data || !img.current) return; const url = URL.createObjectURL(query.data); img.current.src = url; return () => URL.revokeObjectURL(url) }, [query.data, record.image_id])
 return <div ref={host} className={`equipment-thumbnail ${compact ? 'equipment-thumbnail-compact' : ''}`}>{query.data && !query.isError && !failed ? <img ref={img} loading="lazy" alt={record.name} onError={() => setFailedImage(record.image_id)} /> : <Package size={compact ? 22 : 28} aria-label={`${record.name}: ${record.image_id ? 'image unavailable' : 'no catalog image'}`} />}</div>
}
