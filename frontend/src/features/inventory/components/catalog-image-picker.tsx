import { useEffect, useRef } from 'react'
import { FileDropzone } from '@/components/application/file-dropzone'
import { EquipmentImage } from '../pages/equipment-image'
import type { Equipment } from '../types'
export function CatalogImagePicker({ record, file, error, busy, onSelect, onError }: { record?: Equipment; file: File | null; error: string; busy: boolean; onSelect: (file: File | null) => Promise<void>; onError: (error: string) => void }) {
 const preview = useRef<HTMLImageElement>(null)
 useEffect(() => { if (!file || !preview.current) return; const next = URL.createObjectURL(file); preview.current.src = next; return () => URL.revokeObjectURL(next) }, [file])
 return <section className="catalog-image-picker form-section" aria-label="Equipment catalog image"><h3>Catalog image {record ? '' : '(optional)'}</h3><FileDropzone file={file} accept="image/png,image/jpeg" label="Drop an equipment photo here" help="PNG or JPEG · up to 512 KiB · at most 2048 × 2048 pixels and 4 million pixels" busy={busy} error={error} browseLabel={!file && record?.image_id ? 'Replace image' : undefined} preview={file ? <img className="catalog-image-preview" ref={preview} alt="Selected equipment catalog image preview" /> : record?.image_id ? <EquipmentImage record={record} /> : undefined} onError={onError} onSelect={next => { void onSelect(next) }} />{record?.image_id && <p className="management-note">The current image is kept unless a replacement is saved. Removing an unsaved selection keeps it; stored image deletion is not supported.</p>}</section>
}
