import { beforeEach, afterEach, describe, expect, it, vi } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import type { PropsWithChildren } from 'react'
import { apiClient } from '@/lib/api-client'
import { queryKeys } from '@/lib/query-keys'
import { useAuthStore } from '@/stores/auth-store'
import { useCurrentUser } from './hooks/use-auth'
import { authApi } from './api/auth.api'
import { usersApi } from '@/features/users/api/users.api'
import type { UpdateProfileInput } from '@/features/users/types'
import type { AuthUser } from '@/types/common'

vi.mock('@/lib/api-client', () => ({
  apiClient: { get: vi.fn(), post: vi.fn(), patch: vi.fn() },
}))

const currentAccount: AuthUser = {
  id: '00000000-0000-0000-0000-000000000001',
  email: 'synthetic@example.invalid',
  name: 'Current User',
  role: 'user',
  is_active: true,
}

beforeEach(() => {
  vi.clearAllMocks()
  localStorage.clear()
  useAuthStore.setState({ user: null, isAuthenticated: false, accessToken: null, refreshToken: null })
})
afterEach(() => localStorage.clear())

describe('Phase 1B trusted account adapters', () => {
  it('reads /auth/me with cancellation and preserves server role/status', async () => {
    const controller = new AbortController()
    vi.mocked(apiClient.get).mockResolvedValue(currentAccount)
    expect(await authApi.me(controller.signal)).toEqual(currentAccount)
    expect(apiClient.get).toHaveBeenCalledWith('/auth/me', { signal: controller.signal })
  })

  it('hydrates metadata from the current-account endpoint without replacing profile cache or tokens', async () => {
    vi.mocked(apiClient.get).mockResolvedValue(currentAccount)
    useAuthStore.setState({
      isAuthenticated: true,
      user: { ...currentAccount, role: 'admin' },
      accessToken: 'existing-access',
      refreshToken: 'existing-refresh',
    })
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } })
    const profile = { ...currentAccount, created_at: '2026-01-01T00:00:00Z' }
    client.setQueryData(queryKeys.users.me(), profile)
    const wrapper = ({ children }: PropsWithChildren) => (
      <QueryClientProvider client={client}>{children}</QueryClientProvider>
    )
    const { result, unmount } = renderHook(() => useCurrentUser(), { wrapper })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(apiClient.get).toHaveBeenCalledWith('/auth/me', { signal: expect.any(AbortSignal) })
    expect(useAuthStore.getState().user).toEqual(currentAccount)
    expect(client.getQueryData(queryKeys.users.me())).toEqual(profile)
    expect(client.getQueryData(queryKeys.auth.me())).toEqual(currentAccount)
    expect(useAuthStore.getState().refreshToken).toBe('existing-refresh')
    expect(useAuthStore.getState().accessToken).toBe('existing-access')
    unmount()
    client.clear()
  })

  it('does not request current-account metadata without a session', () => {
    const client = new QueryClient()
    const wrapper = ({ children }: PropsWithChildren) => (
      <QueryClientProvider client={client}>{children}</QueryClientProvider>
    )
    const { unmount } = renderHook(() => useCurrentUser(), { wrapper })
    expect(apiClient.get).not.toHaveBeenCalled()
    expect(apiClient.post).not.toHaveBeenCalled()
    unmount()
    client.clear()
  })

  it('sends only a display name even if a runtime caller supplies privileged fields', async () => {
    vi.mocked(apiClient.patch).mockResolvedValue(currentAccount)
    const injected = { name: 'New Name', email: 'other@example.invalid', role: 'admin', is_active: true, id: 'other' }
    await usersApi.updateMe(injected as UpdateProfileInput)
    expect(apiClient.patch).toHaveBeenCalledWith('/users/me', { name: 'New Name' })
  })
})
