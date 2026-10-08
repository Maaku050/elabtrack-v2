import type { AuthUser, LoginInput } from '../types'

// Lazy load the same centralized singleton; no competing session/HTTP layer.
const transport = async () => (await import('@/lib/api-client')).apiClient
export const authApi = {
  login: async (input: LoginInput): Promise<{ user: AuthUser }> => {
    const session = await (await transport()).authenticate('/auth/login', input)
    return { user: session.user }
  },
  refresh: async (): Promise<void> => { await (await transport()).refreshSession() },
  logout: async (): Promise<void> => { await (await transport()).logout() },
  me: async (signal?: AbortSignal): Promise<AuthUser> => (await transport()).get<AuthUser>('/auth/me', { signal }),
}
