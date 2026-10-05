/** Common shared types. */

export type ID = string

export type ISODateString = string

// Temporary generic infrastructure roles, not the final eLabTrack taxonomy.
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
  is_active: boolean
}

export interface BrowserSession {
  access_token: string
  expires_at: ISODateString
  token_type: 'Bearer'
  user: AuthUser
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
