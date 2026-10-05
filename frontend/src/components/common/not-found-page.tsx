import { Link } from 'react-router-dom'
import { buttonVariants } from '@/components/ui/button'
import { Empty, EmptyHeader, EmptyTitle, EmptyDescription, EmptyContent } from '@/components/ui/empty'

export function NotFoundPage() {
  return <main className="flex min-h-svh items-center justify-center p-6"><Empty><EmptyHeader><EmptyTitle><h1>404 — Page not found</h1></EmptyTitle><EmptyDescription>The page you're looking for doesn't exist.</EmptyDescription></EmptyHeader><EmptyContent><Link to="/" className={buttonVariants()}>Back home</Link></EmptyContent></Empty></main>
}
