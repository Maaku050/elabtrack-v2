/**
 * Centralized query key factories.
 *
 * Never scatter raw string keys like ["users"] across the codebase.
 * Always import from here so keys are consistent and refactor-safe.
 */
export const queryKeys = {
  auth: {
    all: ['auth'] as const,
    me: () => [...queryKeys.auth.all, 'me'] as const,
  },
  users: {
    all: ['users'] as const,
    lists: () => [...queryKeys.users.all, 'list'] as const,
    list: (filters: Record<string, unknown> = {}) =>
      [...queryKeys.users.lists(), filters] as const,
    detail: (id: string) => [...queryKeys.users.all, 'detail', id] as const,
    me: () => [...queryKeys.users.all, 'me'] as const,
  },
  // Add new feature query keys here, e.g.:
  // products: { all: ['products'] as const, ... },
} as const
