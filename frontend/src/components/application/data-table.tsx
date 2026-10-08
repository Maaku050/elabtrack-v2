import type { ReactNode } from 'react'
import { ChevronLeft, ChevronRight, Search } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { AppButton, SurfaceCard } from './visual'

export function DataTableShell({ label, children, footer }: { label: string; children: ReactNode; footer?: ReactNode }) {
  return <SurfaceCard className="data-table-shell"><div className="table-scroll" role="region" aria-label={label} tabIndex={0}>{children}</div>{footer}</SurfaceCard>
}

export function SearchField({ label, placeholder, value, onChange }: { label: string; placeholder: string; value: string; onChange: (value: string) => void }) {
  return <label className="search-field"><Search size={19} aria-hidden="true" /><span className="sr-only">{label}</span>
    <Input type="search" placeholder={placeholder} value={value} onChange={e => onChange(e.target.value)} />
  </label>
}

export function PaginationControls({ page, pages, summary, onPageChange }: { page: number; pages: number; summary: string; onPageChange: (page: number) => void }) {
  return <div className="pagination-controls"><p aria-live="polite">{summary}</p><nav aria-label="Table pagination">
    <AppButton variant="outline" aria-label="Previous page" disabled={page <= 1} onClick={() => onPageChange(page - 1)}><ChevronLeft size={17} aria-hidden="true" /></AppButton>
    {Array.from({ length: pages }, (_, i) => i + 1).map(p => <AppButton key={p} variant={page === p ? 'default' : 'outline'} aria-label={`Page ${p}`} aria-current={page === p ? 'page' : undefined} onClick={() => onPageChange(p)}>{p}</AppButton>)}
    <AppButton variant="outline" aria-label="Next page" disabled={page >= pages} onClick={() => onPageChange(page + 1)}><ChevronRight size={17} aria-hidden="true" /></AppButton>
  </nav></div>
}
