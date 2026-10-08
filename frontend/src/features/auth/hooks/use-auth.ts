import { apiErrorMessage } from '@/lib/api-error'
import { useMutation, useQuery } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { authApi } from '../api/auth.api'
import { loginDestination } from '../navigation'
import { useAuthStore } from '@/stores/auth-store'
import { useSessionActionStore } from '@/stores/session-action-store'
import { queryKeys } from '@/lib/query-keys'
import type { LoginInput } from '../types'

export function useLogin(requested?: unknown) {
  const navigate = useNavigate()
  return useMutation({
    onMutate: () => { useSessionActionStore.getState().setSigningIn(true) },
    onSettled: () => { useSessionActionStore.getState().setSigningIn(false) },
    mutationFn: (input: LoginInput) => authApi.login(input),
    onSuccess: ({ user }) => {
      useSessionActionStore.getState().setLogout('idle')
      navigate(loginDestination(user.role, requested), { replace: true })
    },
  })
}

export function useLogout() {
  const navigate = useNavigate()
  return async () => {
    if (useSessionActionStore.getState().logout === 'pending') return
    useSessionActionStore.getState().setLogout('pending')
    // Existing transport immediately clears memory/queries and notifies peers.
    useAuthStore.getState().clear()
    const operation = authApi.logout()
    navigate('/login', { replace: true, state: { signedOut: true } })
    try {
      await operation
      useSessionActionStore.getState().setLogout('idle')
    } catch (error) {
      useSessionActionStore.getState().setLogout('failed', apiErrorMessage(error))
    }
  }
}

export function useCurrentUser(route?: string) {
  const isAuthenticated = useAuthStore(s => s.isAuthenticated)
  return useQuery({
    queryKey: route ? [...queryKeys.auth.me(), route] : queryKeys.auth.me(),
    queryFn: async ({ signal }) => {
      const generation = useAuthStore.getState().generation
      const user = await authApi.me(signal)
      if (!signal.aborted && generation === useAuthStore.getState().generation) useAuthStore.getState().setUser(user)
      return user
    },
    enabled: isAuthenticated,
    staleTime: 0,
    retry: false,
    refetchOnMount: 'always',
    refetchOnWindowFocus: true,
    refetchOnReconnect: true,
  })
}
