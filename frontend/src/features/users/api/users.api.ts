import { apiClient } from '@/lib/api-client'
import type { UpdateProfileInput, User } from '../types'

/**
 * Users API layer. Components go through hooks, not this file directly.
 */
export const usersApi = {
  me(signal?: AbortSignal): Promise<User> {
    return apiClient.get<User>('/users/me', { signal })
  },

  /** Display-name-only self update; identity/security fields are deferred. */
  updateMe(input: UpdateProfileInput): Promise<User> {
    return apiClient.patch<User>('/users/me', { name: input.name })
  },
}
