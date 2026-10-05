import { Link } from 'react-router-dom'
import { Button, buttonVariants } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useHealth } from '@/features/foundation/hooks/use-health'

export function StatusPage() {
  const health = useHealth()
  return (
    <main className="mx-auto flex min-h-svh max-w-2xl flex-col justify-center gap-6 p-6">
      <h1 className="text-2xl font-semibold">eLabTrack V2 service connection</h1>
      <Card>
        <CardHeader><CardTitle>Service status</CardTitle></CardHeader>
        <CardContent className="space-y-4">
          <Button onClick={() => void health.refetch()} disabled={health.isFetching}>Check connection</Button>
          <div role="status" aria-live="polite">
            {health.isFetching ? <p>Checking connection…</p> : health.isError ? <p>Unable to connect. Check that the API and database are running, then try again.</p> : health.data ? (
              <div><p>{health.data.status === 'ok' ? 'API connected.' : 'Service is degraded.'}</p>
                <ul>{Object.entries(health.data.services).map(([service, status]) => <li key={service}>{service}: {status}</li>)}</ul>
              </div>
            ) : <p>Connection has not been checked yet.</p>}
          </div>
        </CardContent>
      </Card>
      <Link to="/" className={buttonVariants({ variant: 'outline' })}>Back to eLabTrack V2</Link>
    </main>
  )
}
