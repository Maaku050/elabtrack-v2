import { create } from 'zustand'
import type { AuthUser, BrowserSession } from '@/types/common'
import { clearLegacyAuthStorage } from '@/lib/storage'
import { clearAuthenticatedQueries } from '@/app/query-client'

export type SessionStatus = 'idle' | 'bootstrapping' | 'authenticated' | 'unauthenticated' | 'error'
interface AuthState {
  user: AuthUser | null
  accessToken: string | null
  status: SessionStatus
  generation: number
  isAuthenticated: boolean
  sessionEnded: boolean
  setSession: (session: BrowserSession) => void
  setUser: (user: AuthUser) => void
  setStatus: (status: SessionStatus) => void
  clear: (status?: SessionStatus) => number
}

// No persistence middleware. User/role are presentation metadata; every API
// authorizes against PostgreSQL. Generation prevents stale async resurrection.
export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  accessToken: null,
  status: 'idle',
  generation: 0,
  isAuthenticated: false,
  sessionEnded: false,
  setSession: (session) => {
    clearLegacyAuthStorage()
    if (get().user?.id !== session.user.id) clearAuthenticatedQueries()
    const roleChanged = !!get().user && get().user?.role !== session.user.role
    if (roleChanged) clearAuthenticatedQueries(true)
    set({ user: session.user, accessToken: session.access_token, status: 'authenticated', isAuthenticated: true, sessionEnded: false, generation: get().generation + (roleChanged ? 1 : 0) })
  },
  setUser: (user) => {
    if (get().isAuthenticated && get().user?.id === user.id) {
      const roleChanged = get().user?.role !== user.role
      if (roleChanged) clearAuthenticatedQueries(true)
      set({ user, generation: get().generation + (roleChanged ? 1 : 0) })
    }
  },
  setStatus: (status) => set({ status }),
  clear: (status = 'unauthenticated') => {
    clearLegacyAuthStorage()
    const generation = get().generation + 1
    const sessionEnded = get().isAuthenticated || get().sessionEnded
    set({ user: null, accessToken: null, status, isAuthenticated: false, generation, sessionEnded })
    clearAuthenticatedQueries()
    return generation
  },
}))
