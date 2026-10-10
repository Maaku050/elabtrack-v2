// datetime-local is deliberately interpreted in Manila, never the device zone.
export function manilaInstant(value: string): string | undefined {
 if (!/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$/.test(value)) return undefined
 const date = new Date(value + ':00+08:00')
 if (!Number.isFinite(date.getTime())) return undefined
 const roundtrip = new Date(date.getTime() + 8 * 60 * 60 * 1000).toISOString().slice(0, 16)
 return roundtrip === value ? date.toISOString() : undefined
}
export function manilaTime(value: string | null | undefined): string {
 return value ? new Intl.DateTimeFormat('en-PH', { timeZone: 'Asia/Manila', dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value)) + ' · Manila' : '—'
}
// Local context is shown once by the containing page/table.
export function manilaTimeCompact(value: string | null | undefined): string { return manilaTime(value).replace(' · Manila', '') }
