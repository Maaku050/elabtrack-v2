import { useRef, useState } from 'react'
import { inventoryApi } from '../api/inventory.api'
import { useInventoryCommand } from './use-inventory'
import type { Equipment } from '../types'
import { validateImageFile } from '@/lib/image-validation'
export const validateCatalogImage = (file: File) => validateImageFile(file,'equipment photograph')
export function useCatalogImage() {
 const [file, setFile] = useState<File | null>(null), [error, setError] = useState(''), [validating, setValidating] = useState(false)
 const selection = useRef(0), attempt = useRef<{ record: Equipment; file: File; key: string } | null>(null)
 const command = useInventoryCommand<{ record: Equipment; file: File; key: string }, Equipment>(v => inventoryApi.saveImage(v.record, v.file, v.key))
 async function select(next: File | null) { const generation = ++selection.current; setFile(null); setError(''); command.reset(); attempt.current = null; if (!next) { setValidating(false); return } setValidating(true); try { await validateCatalogImage(next); if (generation === selection.current) setFile(next) } catch (e) { if (generation === selection.current) setError(e instanceof Error ? e.message : 'Choose a valid catalog image.') } finally { if (generation === selection.current) setValidating(false) } }
 async function save(record: Equipment) { if (!file) return record; if (!attempt.current || attempt.current.file !== file || attempt.current.record.id !== record.id) attempt.current = { record: { ...record }, file, key: crypto.randomUUID() }; const result = await command.mutateAsync(attempt.current); await select(null); return result }
 return { file, error, setError, select, save, validating, command }
}
