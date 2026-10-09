import type { ReactNode, ComponentProps } from 'react'
import { Search, SearchX, CircleAlert, type LucideIcon } from 'lucide-react'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { NativeSelect } from '@/components/ui/native-select'
import { Empty, EmptyHeader, EmptyTitle, EmptyDescription, EmptyMedia, EmptyContent } from '@/components/ui/empty'
import { Skeleton } from '@/components/ui/skeleton'
import { AppButton } from './visual'

export function DirectoryHeader({ title, description, actions, eyebrow }: { title: string; description: string; actions?: ReactNode; eyebrow?: string }) {
  return <header className="directory-header"><div><p className="directory-eyebrow">FSMO / {eyebrow ?? (title === 'Inventory' ? 'Equipment management' : 'Account management')}</p><h1>{title}</h1><p>{description}</p></div>{actions && <div className="directory-actions">{actions}</div>}</header>
}
export function DirectorySearch({ label, ...props }: { label: string } & ComponentProps<typeof Input>) {
  return <label className="directory-filter directory-search"><span>{label}</span><div><Search aria-hidden="true" /><Input type="search" {...props} aria-label={label} /></div></label>
}
export function DirectorySelect({ label, children, ...props }: { label: string; children: ReactNode } & ComponentProps<typeof NativeSelect>) {
  return <label className="directory-filter"><span>{label}</span><NativeSelect {...props} aria-label={label}>{children}</NativeSelect></label>
}
export function DirectoryPanel({ title, summary, toolbar, footer, pending, error, empty, onRetry, children }: { title: string; summary: string; toolbar: ReactNode; footer?: ReactNode; pending: boolean; error?: string; empty?: boolean; onRetry: () => void; children: ReactNode }) {
  return <Card className="directory-panel"><div className="directory-panel-heading"><h2>{title}</h2><p>{summary}</p></div><section className="directory-filters" aria-label={`${title} filters`}>{toolbar}</section><div className="directory-results" aria-busy={pending}>
    {pending ? <div className="directory-loading" role="status"><span className="sr-only">Loading {title.toLowerCase()}…</span>{Array.from({ length: 5 }, (_, i) => <div className="directory-skeleton-row" key={i}><Skeleton className="h-10 w-10" /><Skeleton className="h-5 flex-1" /><Skeleton className="h-5 w-24" /></div>)}</div> : error ? <Empty role="alert"><EmptyHeader><EmptyMedia variant="icon"><CircleAlert aria-hidden="true" /></EmptyMedia><EmptyTitle>Could not load {title.toLowerCase()}</EmptyTitle><EmptyDescription>{error}</EmptyDescription></EmptyHeader><EmptyContent><AppButton variant="outline" onClick={onRetry}>Retry</AppButton></EmptyContent></Empty> : empty ? <Empty><EmptyHeader><EmptyMedia variant="icon"><SearchX aria-hidden="true" /></EmptyMedia><EmptyTitle>No matching records</EmptyTitle><EmptyDescription>Try another search or change the filters.</EmptyDescription></EmptyHeader></Empty> : children}
  </div>{!pending && !error && footer}</Card>
}
export function DirectoryTableRegion({ label, children }: { label: string; children: ReactNode }) { return <div className="directory-table-region" tabIndex={0} role="region" aria-label={label}>{children}</div> }
export function DirectoryStat({ label, value, icon: Icon, tone = 'primary' }: { label: string; value?: number; icon: LucideIcon; tone?: string }) {
  return <Card className="directory-stat" data-tone={tone} data-large={value !== undefined && value >= 1000000}><span className="directory-stat-icon"><Icon aria-hidden="true" /></span><div><p>{label}</p><strong>{value === undefined ? '—' : value.toLocaleString()}</strong></div></Card>
}
