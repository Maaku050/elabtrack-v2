import { describe, it, expect, beforeEach } from 'vitest'
import { useAuthStore } from '@/stores/auth-store'
import type { AuthUser, BrowserSession } from '@/types/common'
import { queryClient } from '@/app/query-client'

const user: AuthUser = {
  id: '00000000-0000-0000-0000-000000000001',
  email: 'test@example.com',
  name: 'Test User',
  role: 'BORROWER',
  is_active: true,
}

const tokens: BrowserSession = {
  access_token: 'access-123',
  expires_at: '2025-01-01T00:00:00Z',
  token_type: 'Bearer',
  user,
}

describe('auth store', () => {
  beforeEach(() => {
    useAuthStore.setState({
      user: null,
      accessToken: null,
      isAuthenticated: false,
      status: 'idle',
      generation: 0,
    })
    localStorage.clear()
    queryClient.clear()
  })

  it('sets a session', () => {
    useAuthStore.getState().setSession(tokens)
    const s = useAuthStore.getState()
    expect(s.isAuthenticated).toBe(true)
    expect(s.user).toEqual(user)
    expect(s.accessToken).toBe('access-123')
    expect(localStorage.getItem('elabtrack_v2.access_token')).toBeNull()
    expect(localStorage.getItem('elabtrack_v2.refresh_token')).toBeNull()
    expect(s).not.toHaveProperty('refreshToken')
  })

  it('clears the session', () => {
    useAuthStore.getState().setSession(tokens)
    useAuthStore.getState().clear()
    const s = useAuthStore.getState()
    expect(s.isAuthenticated).toBe(false)
    expect(s.user).toBeNull()
    expect(localStorage.getItem('elabtrack_v2.access_token')).toBeNull()
  })

  it('clears private queries on logout', () => {
    useAuthStore.getState().setSession(tokens)
    queryClient.setQueryData(['auth', 'private'], { owner: user.id })
    queryClient.setQueryData(['foundation', 'health'], { status: 'ok' })
    useAuthStore.getState().clear()
    expect(queryClient.getQueryData(['auth', 'private'])).toBeUndefined()
    expect(queryClient.getQueryData(['foundation', 'health'])).toEqual({ status: 'ok' })
  })

  it('clears private queries when the account changes', () => {
    useAuthStore.getState().setSession(tokens)
    queryClient.setQueryData(['auth', 'private'], { owner: user.id })
    queryClient.setQueryData(['foundation', 'health'], { status: 'ok' })
    useAuthStore.getState().setSession({ ...tokens, user: { ...user, id: 'other-user' } })
    expect(queryClient.getQueryData(['auth', 'private'])).toBeUndefined()
    expect(queryClient.getQueryData(['foundation', 'health'])).toEqual({ status: 'ok' })
  })

  it('preserves queries when renewing the same account', () => {
    useAuthStore.getState().setSession(tokens)
    queryClient.setQueryData(['auth', 'private'], { owner: user.id })
    useAuthStore.getState().setSession({ ...tokens, access_token: 'renewed' })
    expect(queryClient.getQueryData(['auth', 'private'])).toEqual({ owner: user.id })
  })

  it('updates the user without touching tokens', () => {
    useAuthStore.getState().setSession(tokens)
    useAuthStore.getState().setUser({ ...user, name: 'Renamed' })
    expect(useAuthStore.getState().user?.name).toBe('Renamed')
    expect(useAuthStore.getState().accessToken).toBe('access-123')
  })
  it.each(['current-account', 'rotation'])('fences stale work and removes private cache when a role changes through %s', source => {
    useAuthStore.getState().setSession({ ...tokens, user: { ...user, role: 'ADMIN' } })
    queryClient.setQueryData(['users', 'directory'], { confidential: true })
    queryClient.setQueryData(['auth', 'me', '/admin/reports'], user)
    queryClient.setQueryData(['foundation', 'health'], { status: 'ok' })
    const generation = useAuthStore.getState().generation
    if (source === 'current-account') useAuthStore.getState().setUser(user)
    else useAuthStore.getState().setSession(tokens)
    expect(useAuthStore.getState().generation).toBe(generation + 1)
    expect(queryClient.getQueryData(['users', 'directory'])).toBeUndefined()
    expect(queryClient.getQueryData(['auth', 'me', '/admin/reports'])).toEqual(user)
    expect(queryClient.getQueryData(['foundation', 'health'])).toEqual({ status: 'ok' })
  })

})
