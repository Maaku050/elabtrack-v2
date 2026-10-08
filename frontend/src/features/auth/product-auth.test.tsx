import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { createMemoryRouter, RouterProvider } from 'react-router-dom'
import { QueryClientProvider } from '@tanstack/react-query'
import { routes } from '@/app/router'
import { queryClient } from '@/app/query-client'
import { useAuthStore } from '@/stores/auth-store'
import { useSessionActionStore } from '@/stores/session-action-store'
import { useUIStore } from '@/stores/ui-store'
import { ApiRequestError } from '@/lib/api-error'
import { authApi } from './api/auth.api'
import { loginDestination } from './navigation'
import type { AuthUser, UserRole } from '@/types/common'

const account = (role: UserRole): AuthUser => ({ id: '00000000-0000-0000-0000-000000000001', email: 'synthetic@example.invalid', name: 'Synthetic Account', role, is_active: true })
function authenticate(role: UserRole) {
  useAuthStore.getState().setSession({ access_token: 'unit-memory-only', expires_at: '2030-01-01T00:00:00Z', token_type: 'Bearer', user: account(role) })
}
function open(path = '/login', from?: string) {
  const router = createMemoryRouter(routes, { initialEntries: [{ pathname: path, state: from ? { from } : undefined }] })
  render(<QueryClientProvider client={queryClient}><RouterProvider router={router} /></QueryClientProvider>)
  return router
}
async function credentials() {
  fireEvent.change(await screen.findByLabelText('Email'), { target: { value: 'synthetic@example.invalid' } })
  fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'unit-only-input' } })
}
function submit() { fireEvent.click(screen.getByRole('button', { name: 'Sign In' })) }
function deferred<T>() { let resolve!: (value: T) => void; const promise = new Promise<T>(r => { resolve = r }); return { promise, resolve } }

beforeEach(() => {
  vi.stubGlobal('matchMedia', vi.fn(() => ({ matches: false, addEventListener: vi.fn(), removeEventListener: vi.fn() })))
  queryClient.clear()
  useAuthStore.getState().clear()
  useSessionActionStore.setState({ signingIn: false, logout: 'idle', message: undefined })
  useUIStore.getState().setTheme('light')
  vi.spyOn(authApi, 'login')
  vi.spyOn(authApi, 'me').mockImplementation(async () => useAuthStore.getState().user ?? account('BORROWER'))
  vi.spyOn(authApi, 'logout').mockResolvedValue(undefined)
})
afterEach(() => { vi.restoreAllMocks(); vi.unstubAllGlobals(); queryClient.clear() })

describe('Phase 4A product authentication', () => {
  it('validates fields, exposes accessible errors, and never submits invalid credentials', async () => {
    open(); await screen.findByLabelText('Email'); submit()
    expect(await screen.findByText('Enter a valid email address')).toBeInTheDocument()
    expect(screen.getByLabelText('Email')).toHaveAttribute('aria-invalid', 'true')
    expect(screen.getByLabelText('Password')).toHaveAttribute('aria-describedby', 'password-error')
    expect(authApi.login).not.toHaveBeenCalled()
    expect(screen.queryByRole('link', { name: /register|forgot/i })).not.toBeInTheDocument()
  })
  it('prevents duplicate submissions and supports password visibility', async () => {
    const gate = deferred<{ user: AuthUser }>(); vi.mocked(authApi.login).mockReturnValue(gate.promise)
    open(); await credentials()
    fireEvent.click(screen.getByRole('button', { name: 'Show password' }))
    expect(screen.getByLabelText('Password')).toHaveAttribute('type', 'text')
    fireEvent.click(screen.getByRole('button', { name: 'Hide password' })); submit()
    const pending = await screen.findByRole('button', { name: 'Signing in…' })
    expect(pending).toBeDisabled(); fireEvent.click(pending)
    expect(authApi.login).toHaveBeenCalledTimes(1)
    expect(screen.getByLabelText('Email')).toBeDisabled()
    gate.resolve({ user: account('BORROWER') })
    await waitFor(() => expect(useSessionActionStore.getState().signingIn).toBe(false))
  })
  it.each([
    [401, 'Unable to sign in.'], [0, 'Unable to connect.'], [429, 'Too many sign-in attempts.'], [503, 'Reference: 00000000-0000-4000-8000-000000000001'],
  ])('handles failure %s safely with a retryable form', async (status, text) => {
    vi.mocked(authApi.login).mockRejectedValue(new ApiRequestError('Something went wrong. Please try again later.', Number(status), 'REQUEST_FAILED', undefined, '00000000-0000-4000-8000-000000000001'))
    open(); await credentials(); submit()
    expect(await screen.findByText(content => content.includes(String(text)))).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Sign In' })).toBeEnabled()
    expect(screen.getByLabelText('Password')).toHaveValue('')
    expect(screen.getByLabelText('Password')).toHaveFocus()
  })
  it.each([
    ['BORROWER', '/borrower/home'], ['STAFF', '/staff/dashboard'], ['ADMIN', '/staff/dashboard'],
  ] as const)('navigates %s using safe server account metadata', async (role, destination) => {
    vi.mocked(authApi.login).mockImplementation(async () => { authenticate(role); return { user: account(role) } })
    const router = open(); await credentials(); submit()
    await waitFor(() => expect(router.state.location.pathname).toBe(destination))
    expect(await screen.findByText('This feature is not available yet')).toBeInTheDocument()
    expect(authApi.login).toHaveBeenCalledWith({ email: 'synthetic@example.invalid', password: 'unit-only-input' })
  })
  it('preserves an allowed deep link after an anonymous redirect and login', async () => {
    vi.mocked(authApi.login).mockImplementation(async () => { authenticate('BORROWER'); return { user: account('BORROWER') } })
    const router = open('/borrower/account'); await credentials()
    expect(router.state.location.pathname).toBe('/login'); submit()
    expect(await screen.findByRole('heading', { name: 'Account details' })).toBeInTheDocument()
    expect(router.state.location.pathname).toBe('/borrower/account')
  })
  it.each(['https://untrusted.example', '//untrusted.example', '/staff/dashboard', '/borrower/../admin/reports', '/__preview/borrower/home'])('rejects unsafe Borrower return destination %s', destination => {
    expect(loginDestination('BORROWER', destination)).toBe('/borrower/home')
  })
  it.each([
    ['BORROWER', '/staff/dashboard', false], ['STAFF', '/admin/reports', false], ['ADMIN', '/admin/reports', true],
    ['STAFF', '/staff/dashboard', true], ['BORROWER', '/borrower/home', true], ['ADMIN', '/borrower/home', false],
  ] as const)('guards %s access to %s', async (role, path, allowed) => {
    authenticate(role); open(path)
    expect(await screen.findByText(allowed ? 'This feature is not available yet' : 'Access denied')).toBeInTheDocument()
    if (!allowed) expect(screen.queryByText('This feature is not available yet')).not.toBeInTheDocument()
  })
  it('withholds all protected content while current database authority is loading', async () => {
    authenticate('ADMIN'); const gate = deferred<AuthUser>(); vi.mocked(authApi.me).mockReturnValue(gate.promise)
    open('/admin/reports')
    expect(await screen.findByText('Confirming account access…')).toBeInTheDocument()
    expect(screen.queryByRole('heading', { name: 'Reports' })).not.toBeInTheDocument()
    gate.resolve(account('BORROWER'))
    expect(await screen.findByText('Access denied')).toBeInTheDocument()
    expect(useAuthStore.getState().user?.role).toBe('BORROWER')
  })
  it('denies inactive metadata and provides recoverable access errors for transient failures', async () => {
    authenticate('BORROWER'); vi.mocked(authApi.me).mockRejectedValueOnce(new ApiRequestError('Unable to connect.', 0)).mockResolvedValueOnce({ ...account('BORROWER'), is_active: false })
    open('/borrower/home')
    expect(await screen.findByRole('heading', { name: 'Unable to confirm access' })).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: 'Try again' }))
    expect(await screen.findByText('Access denied')).toBeInTheDocument()
    expect(screen.queryByText('This feature is not available yet')).not.toBeInTheDocument()
  })
  it('clears private memory immediately during sign out and permits retry after network failure', async () => {
    authenticate('STAFF'); open('/staff/dashboard'); await screen.findByText('This feature is not available yet')
    queryClient.setQueryData(['users', 'private'], { confidential: true })
    const gate = deferred<void>(); vi.mocked(authApi.logout).mockReturnValueOnce(gate.promise)
    fireEvent.click(screen.getByRole('button', { name: 'Sign Out' }))
    expect(await screen.findByText('Signing out…')).toBeInTheDocument()
    expect(useAuthStore.getState().accessToken).toBeNull()
    expect(queryClient.getQueryData(['users', 'private'])).toBeUndefined()
    gate.resolve(); await waitFor(() => expect(useSessionActionStore.getState().logout).toBe('idle'))
    expect(await screen.findByRole('button', { name: 'Sign In' })).toBeEnabled()
    vi.mocked(authApi.logout).mockRejectedValueOnce(new ApiRequestError('Unable to connect.', 0))
    authenticate('BORROWER'); await screen.findByText('This feature is not available yet')
    fireEvent.click(screen.getByRole('button', { name: 'Sign Out' }))
    expect(await screen.findByText(/Server sign-out could not be confirmed/)).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: 'Retry sign out' }))
    await waitFor(() => expect(useSessionActionStore.getState().logout).toBe('idle'))
    expect(useAuthStore.getState().isAuthenticated).toBe(false)
  })
  it('switches both login and forbidden experiences using the shared theme control', async () => {
    open(); await screen.findByLabelText('Email')
    fireEvent.click(screen.getByRole('button', { name: 'Switch to dark theme' }))
    expect(document.documentElement).toHaveClass('dark')
  })
})
