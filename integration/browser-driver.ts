// Built only by build-browser.mjs into ignored output and copied to the
// disposable nginx container. Never included by the production entrypoint.
import { apiClient } from '../frontend/src/lib/api-client'
import { useAuthStore } from '../frontend/src/stores/auth-store'
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
} })
