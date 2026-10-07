// Built only by build-browser.mjs into ignored output and copied to the
// disposable nginx container. Never included by the production entrypoint.
import { apiClient } from '../frontend/src/lib/api-client'
import { useAuthStore } from '../frontend/src/stores/auth-store'
import { queryClient } from '../frontend/src/app/query-client'
import { sessionCoordinator } from '../frontend/src/lib/session-coordinator'
import '../frontend/src/main'

Object.assign(window, { phase1g: {
  login: (email: string, password: string) => apiClient.authenticate('/auth/login', { email, password }).then(() => undefined),
  refresh: () => apiClient.refreshSession().then(() => undefined),
  logout: () => apiClient.logout(),
  me: () => apiClient.get('/auth/me'),
  state: () => { const s = useAuthStore.getState(); return { status: s.status, hasAccess: !!s.accessToken, user: s.user, generation: s.generation } },
  // An invalid memory token stimulates the real 401/recovery boundary without
  // changing cookies, the API, or production security behavior.
  invalidateAccess: () => useAuthStore.setState({ accessToken: 'invalid-synthetic-access' }),
  initialMemoryEmpty: useAuthStore.getState().accessToken === null,
  coordination: () => ({ exclusive: sessionCoordinator.exclusive, broadcast: typeof BroadcastChannel !== 'undefined' }),
  cachePrivate: () => {
    queryClient.setQueryData(['auth', 'me'], useAuthStore.getState().user)
    queryClient.setQueryData(['foundation', 'health'], { status: 'ok' })
  },
  cacheState: () => ({ privatePresent: queryClient.getQueryData(['auth', 'me']) !== undefined, publicPresent: queryClient.getQueryData(['foundation', 'health']) !== undefined }),
} })
