import { create } from 'zustand'

// Client-only lifecycle feedback across route unmounts; never credentials.
export const useSessionActionStore = create<{
  signingIn: boolean
  setSigningIn: (value: boolean) => void
  logout: 'idle' | 'pending' | 'failed'
  message?: string
  setLogout: (logout: 'idle' | 'pending' | 'failed', message?: string) => void
}>((set) => ({ signingIn: false, setSigningIn: (signingIn) => set({ signingIn }), logout: 'idle', setLogout: (logout, message) => set({ logout, message }) }))
