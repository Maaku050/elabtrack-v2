import { describe, it, expect, beforeEach } from 'vitest'
import { useAuthStore } from '@/stores/auth-store'
import type { AuthUser, TokenPair } from '@/types/common'
import { queryClient } from '@/app/query-client'

const user: AuthUser = {
  id: '00000000-0000-0000-0000-000000000001',
  email: 'test@example.com',
  name: 'Test User',
  role: 'user',
  is_active: true,
}

const tokens: TokenPair = {
  access_token: 'access-123',
  refresh_token: 'refresh-456',
  expires_at: '2025-01-01T00:00:00Z',
  token_type: 'Bearer',
}

describe('auth store', () => {
  beforeEach(() => {
    useAuthStore.setState({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
      isHydrated: false,
    })
    localStorage.clear()
    queryClient.clear()
  })

  it('sets a session', () => {
    useAuthStore.getState().setSession(user, tokens)
    const s = useAuthStore.getState()
    expect(s.isAuthenticated).toBe(true)
    expect(s.user).toEqual(user)
    expect(s.accessToken).toBe('access-123')
    expect(localStorage.getItem('elabtrack_v2.access_token')).toBe('access-123')
  })

  it('clears the session', () => {
    useAuthStore.getState().setSession(user, tokens)
    useAuthStore.getState().clear()
    const s = useAuthStore.getState()
    expect(s.isAuthenticated).toBe(false)
    expect(s.user).toBeNull()
    expect(localStorage.getItem('elabtrack_v2.access_token')).toBeNull()
  })

  it('clears private queries on logout', () => {
    useAuthStore.getState().setSession(user, tokens)
    queryClient.setQueryData(['private'], { owner: user.id })
    useAuthStore.getState().clear()
    expect(queryClient.getQueryData(['private'])).toBeUndefined()
  })

  it('clears private queries when the account changes', () => {
    useAuthStore.getState().setSession(user, tokens)
    queryClient.setQueryData(['private'], { owner: user.id })
    useAuthStore.getState().setSession({ ...user, id: 'other-user' }, tokens)
    expect(queryClient.getQueryData(['private'])).toBeUndefined()
  })

  it('preserves queries when renewing the same account', () => {
    useAuthStore.getState().setSession(user, tokens)
    queryClient.setQueryData(['private'], { owner: user.id })
    useAuthStore.getState().setSession(user, { ...tokens, access_token: 'renewed' })
    expect(queryClient.getQueryData(['private'])).toEqual({ owner: user.id })
  })

  it('updates the user without touching tokens', () => {
    useAuthStore.getState().setSession(user, tokens)
    useAuthStore.getState().setUser({ ...user, name: 'Renamed' })
    expect(useAuthStore.getState().user?.name).toBe('Renamed')
    expect(useAuthStore.getState().accessToken).toBe('access-123')
  })
})
