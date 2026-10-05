import { useEffect, type ReactNode } from 'react'
import { useAuthStore } from '@/stores/auth-store'
import { Button } from '@/components/ui/button'

// Load the transport on demand; the shared module still owns one bootstrap
// promise across StrictMode effects. Chunk load failure is recoverable UI state.
async function restoreSession(retry = false): Promise<void> {
  try {
    const { apiClient } = await import('@/lib/api-client')
    await (retry ? apiClient.retryBootstrap() : apiClient.bootstrap())
  } catch { useAuthStore.getState().clear('error') }
}

export function SessionBootstrap({ children }: { children: ReactNode }) {
  const status = useAuthStore((state) => state.status)
  useEffect(() => { void restoreSession() }, [])
  if (status === 'idle' || status === 'bootstrapping') {
    return <div role="status" className="p-4">Restoring session…</div>
  }
  return <>
    {status === 'error' && <div role="status" className="flex items-center gap-3 p-4">
      <span>Unable to restore your session.</span>
      <Button onClick={() => { void restoreSession(true) }}>Retry session</Button>
    </div>}
    {children}
  </>
}
