import { apiErrorMessage } from '@/lib/api-error'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { usersApi } from '../api/users.api'
import { queryKeys } from '@/lib/query-keys'
import { useToast } from '@/components/feedback/toast'
import { useAuthStore } from '@/stores/auth-store'
import type { UpdateProfileInput, User } from '../types'

/**
 * Current-user profile query (full User record, not just auth metadata).
 */
export function useProfile() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  return useQuery({
    queryKey: queryKeys.users.me(),
    queryFn: ({ signal }) => usersApi.me(signal),
    enabled: isAuthenticated,
  })
}

/**
 * Optimistic profile update.
 *
 * This is a display-name-only mutation, so we use
 * the optimistic UI pattern:
 *   1. cancel related queries
 *   2. snapshot previous cache
 *   3. optimistically update cache
 *   4. send mutation
 *   5. on error: rollback + toast
 *   6. on settle: invalidate to revalidate
 */
export function useUpdateProfile() {
  const queryClient = useQueryClient()
  const toast = useToast()
  const setUser = useAuthStore((s) => s.setUser)

  return useMutation({
    mutationFn: (input: UpdateProfileInput) => usersApi.updateMe(input),

    onMutate: async (input) => {
      // 1. Cancel any outgoing refetches so they don't overwrite our optimistic update.
      await queryClient.cancelQueries({ queryKey: queryKeys.users.me() })

      // 2. Snapshot the previous value for rollback.
      const previous = queryClient.getQueryData<User>(queryKeys.users.me())

      // 3. Optimistically update the cache.
      if (previous) {
        const next: User = {
          ...previous,
          name: input.name ?? previous.name,
          updated_at: new Date().toISOString(),
        }
        queryClient.setQueryData<User>(queryKeys.users.me(), next)
        // Keep the auth store's lightweight metadata in sync too.
        setUser({ id: next.id, name: next.name, email: next.email, role: next.role, is_active: next.is_active })
      }

      return { previous }
    },

    onError: (err, _input, context) => {
      // 5. Rollback to the snapshot on failure.
      if (context?.previous) {
        queryClient.setQueryData(queryKeys.users.me(), context.previous)
        setUser({
          id: context.previous.id,
          name: context.previous.name,
          email: context.previous.email,
          role: context.previous.role,
          is_active: context.previous.is_active,
        })
      }
      toast.error('Update failed', apiErrorMessage(err))
    },

    onSuccess: (user) => {
      setUser({ id: user.id, name: user.name, email: user.email, role: user.role, is_active: user.is_active })
      toast.success('Profile updated', 'Your changes have been saved.')
    },

    onSettled: () => {
      // 6. Revalidate to make sure our optimistic state matches the server.
      void queryClient.invalidateQueries({ queryKey: queryKeys.users.me() })
      void queryClient.invalidateQueries({ queryKey: queryKeys.auth.me() })
    },
  })
}
