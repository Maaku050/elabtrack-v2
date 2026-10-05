import { create } from 'zustand'
import type { AuthUser, BrowserSession } from '@/types/common'
import { clearLegacyAuthStorage } from '@/lib/storage'
import { queryClient } from '@/app/query-client'

export type SessionStatus = 'idle' | 'bootstrapping' | 'authenticated' | 'unauthenticated' | 'error'
interface AuthState {
  user: AuthUser | null
  accessToken: string | null
  status: SessionStatus
  generation: number
  isAuthenticated: boolean
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
  setSession: (session) => {
    clearLegacyAuthStorage()
    if (get().user?.id !== session.user.id) queryClient.clear()
    set({ user: session.user, accessToken: session.access_token, status: 'authenticated', isAuthenticated: true })
  },
  setUser: (user) => {
    if (get().isAuthenticated && get().user?.id === user.id) set({ user })
  },
  setStatus: (status) => set({ status }),
  clear: (status = 'unauthenticated') => {
    clearLegacyAuthStorage()
    const generation = get().generation + 1
    set({ user: null, accessToken: null, status, isAuthenticated: false, generation })
    queryClient.clear()
    return generation
  },
}))
