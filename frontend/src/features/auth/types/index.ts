import type { AuthUser, BrowserSession } from '@/types/common'

export type { AuthUser, BrowserSession }

export interface LoginInput {
  email: string
  password: string
}

export interface RegisterInput {
  email: string
  name: string
  password: string
}
