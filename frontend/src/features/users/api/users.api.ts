import { apiClient } from '@/lib/api-client'
import type { UpdateProfileInput, User } from '../types'

/**
 * Users API layer. Components go through hooks, not this file directly.
 */
export const usersApi = {
  me(signal?: AbortSignal): Promise<User> {
    return apiClient.get<User>('/users/me', { signal })
  },

  /**
   * Partial update of the current user's profile. Only non-undefined fields
   * are sent; the backend treats null/absent fields as "no change".
   */
  updateMe(input: UpdateProfileInput): Promise<User> {
    const body: Record<string, unknown> = {}
    if (input.name !== undefined) body.name = input.name
    if (input.email !== undefined) body.email = input.email
    return apiClient.patch<User>('/users/me', body)
  },
}
