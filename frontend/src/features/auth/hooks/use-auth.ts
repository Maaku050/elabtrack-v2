import { useMutation, useQuery } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { authApi } from '../api/auth.api'
import { useAuthStore } from '@/stores/auth-store'
import { useToast } from '@/components/feedback/toast'
import { queryKeys } from '@/lib/query-keys'
import type { LoginInput, RegisterInput } from '../types'

/**
 * Login mutation. The centralized client sets the memory session before navigation.
 */
export function useLogin() {
  const navigate = useNavigate()
  const toast = useToast()

  return useMutation({
    mutationFn: (input: LoginInput) => authApi.login(input),
    onSuccess: (data) => {
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
 * Local registration mutation; the centralized client owns session state.
 */
export function useRegister() {
  const navigate = useNavigate()
  const toast = useToast()

  return useMutation({
    mutationFn: (input: RegisterInput) => authApi.register(input),
    onSuccess: (data) => {
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
  const navigate = useNavigate()

  return () => {
    void authApi.logout().catch(() => {})
    navigate('/', { replace: true })
  }
}

/**
 * Current-user query refreshes safe presentation metadata from the server.
 * Bootstrap already accepts current account state with its access response.
 */
export function useCurrentUser() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  const setUser = useAuthStore((s) => s.setUser)

  return useQuery({
    queryKey: queryKeys.auth.me(),
    queryFn: async ({ signal }) => {
      const generation = useAuthStore.getState().generation
      const user = await authApi.me(signal)
      if (!signal.aborted && generation === useAuthStore.getState().generation) setUser(user)
      return user
    },
    enabled: isAuthenticated,
    staleTime: 60_000,
  })
}
