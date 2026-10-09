import { useRef, useState } from 'react'
import { inventoryApi } from '../api/inventory.api'
import { useInventoryCommand } from './use-inventory'
import type { Equipment } from '../types'
export async function validateCatalogImage(file: File) {
 if (file.size === 0 || file.size > 512 * 1024) throw Error('Choose a nonempty image no larger than 512 KiB.')
 if (!['image/png', 'image/jpeg'].includes(file.type)) throw Error('Choose a PNG or JPEG equipment photograph.')
 const bytes = new Uint8Array(await file.slice(0, 8).arrayBuffer())
 const png = bytes.length >= 8 && [137,80,78,71,13,10,26,10].every((n, i) => bytes[i] === n), jpeg = bytes[0] === 255 && bytes[1] === 216 && bytes[2] === 255
 if ((file.type === 'image/png' && !png) || (file.type === 'image/jpeg' && !jpeg)) throw Error('The file contents must match its PNG or JPEG format.')
 let bitmap: ImageBitmap
 try { bitmap = await createImageBitmap(file) } catch { throw Error('This image could not be decoded. Choose a valid PNG or JPEG.') }
 const valid = bitmap.width > 0 && bitmap.height > 0 && bitmap.width <= 2048 && bitmap.height <= 2048
 bitmap.close(); if (!valid) throw Error('Images must be at most 2048 pixels per dimension.')
}
export function useCatalogImage() {
 const [file, setFile] = useState<File | null>(null), [error, setError] = useState(''), [validating, setValidating] = useState(false)
 const selection = useRef(0), attempt = useRef<{ record: Equipment; file: File; key: string } | null>(null)
 const command = useInventoryCommand<{ record: Equipment; file: File; key: string }, Equipment>(v => inventoryApi.saveImage(v.record, v.file, v.key))
 async function select(next: File | null) { const generation = ++selection.current; setFile(null); setError(''); command.reset(); attempt.current = null; if (!next) { setValidating(false); return } setValidating(true); try { await validateCatalogImage(next); if (generation === selection.current) setFile(next) } catch (e) { if (generation === selection.current) setError(e instanceof Error ? e.message : 'Choose a valid catalog image.') } finally { if (generation === selection.current) setValidating(false) } }
 async function save(record: Equipment) { if (!file) return record; if (!attempt.current || attempt.current.file !== file || attempt.current.record.id !== record.id) attempt.current = { record: { ...record }, file, key: crypto.randomUUID() }; const result = await command.mutateAsync(attempt.current); await select(null); return result }
 return { file, error, setError, select, save, validating, command }
}
