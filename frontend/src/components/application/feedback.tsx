import type { ReactNode } from 'react'
import { CircleAlert } from 'lucide-react'
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert'
import { Skeleton } from '@/components/ui/skeleton'
import { Dialog, DialogContent, DialogTitle, DialogDescription, DialogClose } from '@/components/ui/dialog'
import { AppButton, SurfaceCard } from './visual'

export function ConflictAlert({ title = 'This information has changed', children, onRefresh }: { title?: string; children: ReactNode; onRefresh?: () => void }) {
  return <Alert className="conflict-alert"><CircleAlert aria-hidden="true" /><AlertTitle>{title}</AlertTitle><AlertDescription>{children}
    {onRefresh && <AppButton variant="outline" onClick={onRefresh}>Refresh and review</AppButton>}
  </AlertDescription></Alert>
}

export function LoadingRegion({ label }: { label: string }) {
  return <SurfaceCard className="loading-region" role="status" aria-busy="true" aria-label={label}>
    <Skeleton className="h-5 w-1/2" /><Skeleton className="h-10 w-full" /><Skeleton className="h-5 w-3/4" />
    <span className="sr-only">{label}</span>
  </SurfaceCard>
}

/** Read-only preview disclosure, not a confirmation for a business mutation. */
export function PreviewDialog({ title, children, open, onOpenChange }: { title: string; children: ReactNode; open: boolean; onOpenChange: (open: boolean) => void }) {
  return <Dialog open={open} onOpenChange={onOpenChange}><DialogContent className="app-dialog">
    <DialogTitle>{title}</DialogTitle><DialogDescription>Development visual preview · synthetic data. No business action is recorded.</DialogDescription>
    <div className="dialog-body">{children}</div><DialogClose render={<AppButton variant="outline" />}>Close preview</DialogClose>
  </DialogContent></Dialog>
}
