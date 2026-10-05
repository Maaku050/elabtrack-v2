import type { ReactNode } from 'react'
import { Inbox, AlertTriangle } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Spinner } from '@/components/ui/spinner'
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert'
import { Empty, EmptyHeader, EmptyMedia, EmptyTitle, EmptyDescription, EmptyContent } from '@/components/ui/empty'

export function EmptyState({ title = 'Nothing here yet', description, action, className }: { title?: string; description?: string; action?: ReactNode; className?: string }) {
  return <Empty className={cn('border border-dashed', className)}><EmptyHeader><EmptyMedia variant="icon"><Inbox aria-hidden="true" /></EmptyMedia><EmptyTitle>{title}</EmptyTitle>{description && <EmptyDescription>{description}</EmptyDescription>}</EmptyHeader>{action && <EmptyContent>{action}</EmptyContent>}</Empty>
}

export function ErrorState({ title = 'Something went wrong', description, onRetry, className }: { title?: string; description?: string; onRetry?: () => void; className?: string }) {
  return <Alert variant="destructive" className={className}><AlertTriangle aria-hidden="true" /><AlertTitle>{title}</AlertTitle><AlertDescription>{description}{onRetry && <Button variant="outline" onClick={onRetry} className="mt-3 w-fit">Try again</Button>}</AlertDescription></Alert>
}

export function LoadingState({ label = 'Loading…', className }: { label?: string; className?: string }) {
  return <div role="status" className={cn('flex items-center justify-center gap-2 p-10 text-sm text-muted-foreground', className)}><Spinner aria-hidden="true" />{label}</div>
}
