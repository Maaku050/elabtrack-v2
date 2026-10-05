import { useMutation, useQuery } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { authApi } from '../api/auth.api'
import { useAuthStore } from '@/stores/auth-store'
import { storage, storageKeys } from '@/lib/storage'
import { useToast } from '@/components/feedback/toast'
import { queryKeys } from '@/lib/query-keys'
import type { LoginInput, RegisterInput } from '../types'

/**
 * Login mutation. On success, stores tokens + user and navigates to the app.
 */
export function useLogin() {
  const setSession = useAuthStore((s) => s.setSession)
  const navigate = useNavigate()
  const toast = useToast()

  return useMutation({
    mutationFn: (input: LoginInput) => authApi.login(input),
    onSuccess: (data) => {
      setSession(data.user, data)
      toast.success('Welcome back!', `Signed in as ${data.user.email}`)
      navigate('/', { replace: true })
    },
    onError: (err) => {
      const message = err instanceof Error ? err.message : 'Login failed'
      toast.error('Login failed', message)
    },
  })
}

/**
 * Register mutation. On success, stores tokens + user and navigates to the app.
 */
export function useRegister() {
  const setSession = useAuthStore((s) => s.setSession)
  const navigate = useNavigate()
  const toast = useToast()

  return useMutation({
    mutationFn: (input: RegisterInput) => authApi.register(input),
    onSuccess: (data) => {
      setSession(data.user, data)
      toast.success('Account created', `Welcome, ${data.user.name}!`)
      navigate('/', { replace: true })
    },
    onError: (err) => {
      const message = err instanceof Error ? err.message : 'Registration failed'
      toast.error('Registration failed', message)
    },
  })
}

/**
 * Logout mutation. Clears the session regardless of the server response
 * (idempotent) and navigates to the placeholder.
 */
export function useLogout() {
  const clear = useAuthStore((s) => s.clear)
  const navigate = useNavigate()

  return () => {
    const refreshToken = storage.get(storageKeys.refreshToken)
    if (refreshToken) {
      // Fire-and-forget; we clear the session either way.
      authApi.logout({ refresh_token: refreshToken }).catch(() => {})
    }
    clear()
    navigate('/', { replace: true })
  }
}

/**
 * Current-user query. Used on app load to hydrate the auth store with the
 * actual user profile (the store only knows there *is* a session).
 */
export function useCurrentUser() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  const setUser = useAuthStore((s) => s.setUser)

  return useQuery({
    queryKey: queryKeys.auth.me(),
    queryFn: async ({ signal }) => {
      const user = await authApi.me(signal)
      if (!signal.aborted) setUser(user)
      return user
    },
    enabled: isAuthenticated,
    staleTime: 60_000,
  })
}
