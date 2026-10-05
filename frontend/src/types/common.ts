/** Common shared types. */

export type ID = string

export type ISODateString = string

export type UserRole = 'user' | 'admin'

export interface User {
  id: ID
  email: string
  name: string
  role: UserRole
  is_active: boolean
  created_at: ISODateString
  updated_at: ISODateString
}

export interface AuthUser {
  id: ID
  email: string
  name: string
  role: UserRole
}

export interface TokenPair {
  access_token: string
  refresh_token: string
  expires_at: ISODateString
  token_type: string
}

export type RequestState =
  | 'idle'
  | 'loading'
  | 'success'
  | 'empty'
  | 'error'
  | 'refreshing'
  | 'submitting'
  | 'optimistic'
  | 'disabled'
