import { apiClient } from '@/lib/api-client'
import type { AuthUser, LoginInput, RegisterInput } from '../types'

export const authApi = {
  login: async (input: LoginInput): Promise<{ user: AuthUser }> => {
    const session = await apiClient.authenticate('/auth/login', input)
    return { user: session.user }
  },
  // Local development/test only; production registration remains absent.
  register: async (input: RegisterInput): Promise<{ user: AuthUser }> => {
    const session = await apiClient.authenticate('/auth/register', input)
    return { user: session.user }
  },
  refresh: async (): Promise<void> => { await apiClient.refreshSession() },
  logout: () => apiClient.logout(),
  me: (signal?: AbortSignal): Promise<AuthUser> => apiClient.get<AuthUser>('/auth/me', { signal }),
}
