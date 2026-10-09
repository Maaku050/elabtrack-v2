import { create } from 'zustand'

// Client-only lifecycle feedback across route unmounts; never credentials.
export const useSessionActionStore = create<{
  signingIn: boolean
  setSigningIn: (value: boolean) => void
  logout: 'idle' | 'pending' | 'failed'
  // Presentation only: set on a terminal failure of a live session.
  sessionNotice: 'invalidated' | null
  setSessionNotice: (notice: 'invalidated' | null) => void
  message?: string
  setLogout: (logout: 'idle' | 'pending' | 'failed', message?: string) => void
}>((set) => ({ sessionNotice: null, setSessionNotice: (sessionNotice) => set({ sessionNotice }), signingIn: false, setSigningIn: (signingIn) => set({ signingIn }), logout: 'idle', setLogout: (logout, message) => set({ logout, message }) }))
