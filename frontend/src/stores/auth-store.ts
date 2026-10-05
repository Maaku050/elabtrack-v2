import { create } from 'zustand'
import type { AuthUser, TokenPair } from '@/types/common'
import { storage, storageKeys } from '@/lib/storage'
import { queryClient } from '@/app/query-client'

interface AuthState {
  user: AuthUser | null
  accessToken: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  isHydrated: boolean
  setSession: (user: AuthUser, tokens: TokenPair) => void
  setUser: (user: AuthUser) => void
  hydrate: () => void
  clear: () => void
}

/**
 * Auth store — holds only session *metadata* and tokens.
 *
 * Server data (e.g. full user profile, lists of users) belongs in TanStack
 * Query. This store only keeps what's needed to gate the UI: who is logged
 * in and the tokens to attach to requests.
 */
export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  accessToken: null,
  refreshToken: null,
  isAuthenticated: false,
  isHydrated: false,

  setSession: (user, tokens) => {
    if (get().user?.id !== user.id) queryClient.clear()
    storage.set(storageKeys.accessToken, tokens.access_token)
    storage.set(storageKeys.refreshToken, tokens.refresh_token)
    set({
      user,
      accessToken: tokens.access_token,
      refreshToken: tokens.refresh_token,
      isAuthenticated: true,
    })
  },

  setUser: (user) => set({ user }),

  hydrate: () => {
    const accessToken = storage.get(storageKeys.accessToken)
    const refreshToken = storage.get(storageKeys.refreshToken)
    set({
      accessToken,
      refreshToken,
      isAuthenticated: !!accessToken,
      isHydrated: true,
    })
  },

  clear: () => {
    queryClient.clear()
    storage.remove(storageKeys.accessToken)
    storage.remove(storageKeys.refreshToken)
    set({ user: null, accessToken: null, refreshToken: null, isAuthenticated: false })
  },
}))
