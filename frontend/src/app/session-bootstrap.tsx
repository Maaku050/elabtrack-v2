import { useEffect, type ReactNode } from 'react'
import { useAuthStore } from '@/stores/auth-store'
import { AppBrand, AppButton, ThemeControl, SurfaceCard } from '@/components/application/visual'
import { LoadingState } from '@/components/feedback/states'
import { useSessionActionStore } from '@/stores/session-action-store'

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
  const signingIn = useSessionActionStore(state => state.signingIn)
  useEffect(() => { void restoreSession() }, [])
  if ((status === 'idle' || status === 'bootstrapping') && !signingIn) {
    return <main className="auth-page"><header className="auth-top"><AppBrand /><ThemeControl /></header><SurfaceCard className="auth-card"><LoadingState label="Restoring session…" /></SurfaceCard></main>
  }
  return <>
    {status === 'error' && <div role="status" className="session-recovery">
      <span>Unable to restore your session.</span>
      <AppButton variant="outline" onClick={() => { void restoreSession(true) }}>Retry session</AppButton>
    </div>}
    {children}
  </>
}
