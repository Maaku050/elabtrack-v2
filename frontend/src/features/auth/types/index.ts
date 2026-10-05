import type { AuthUser, TokenPair } from '@/types/common'

export type { AuthUser, TokenPair }

export interface LoginInput {
  email: string
  password: string
}

export interface RegisterInput {
  email: string
  name: string
  password: string
}

export interface RefreshInput {
  refresh_token: string
}
