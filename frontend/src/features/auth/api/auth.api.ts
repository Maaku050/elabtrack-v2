import { apiClient } from '@/lib/api-client'
import { storage, storageKeys } from '@/lib/storage'
import type { AuthUser, LoginInput, RegisterInput, RefreshInput, TokenPair } from '../types'

/**
 * Auth API layer. Components never call this directly — they go through
 * the hooks in `features/auth/hooks/`.
 */
export const authApi = {
  login(input: LoginInput): Promise<TokenPair & { user: AuthUser }> {
    // The backend returns a TokenPair at /auth/login. We persist the tokens
    // FIRST so the apiClient attaches them to the subsequent /auth/me call,
    // then fetch the user so the store has session metadata.
    return apiClient
      .post<TokenPair>('/auth/login', input)
      .then(async (tokens) => {
        storage.set(storageKeys.accessToken, tokens.access_token)
        storage.set(storageKeys.refreshToken, tokens.refresh_token)
        const user = await apiClient.get<AuthUser>('/auth/me')
        return { ...tokens, user }
      })
  },

  // Local development/test convenience only; production does not mount this route.
  register(input: RegisterInput): Promise<TokenPair & { user: AuthUser }> {
    return apiClient
      .post<TokenPair>('/auth/register', input)
      .then(async (tokens) => {
        storage.set(storageKeys.accessToken, tokens.access_token)
        storage.set(storageKeys.refreshToken, tokens.refresh_token)
        const user = await apiClient.get<AuthUser>('/auth/me')
        return { ...tokens, user }
      })
  },

  refresh(input: RefreshInput): Promise<TokenPair> {
    return apiClient.post<TokenPair>('/auth/refresh', input)
  },

  logout(input: RefreshInput): Promise<void> {
    return apiClient.post<void>('/auth/logout', input).finally(() => {
      // Clear tokens locally regardless of the server response so the
      // client never keeps stale credentials after a logout attempt.
      storage.remove(storageKeys.accessToken)
      storage.remove(storageKeys.refreshToken)
    })
  },

  me(signal?: AbortSignal): Promise<AuthUser> {
    return apiClient.get<AuthUser>('/auth/me', { signal })
  },
}
