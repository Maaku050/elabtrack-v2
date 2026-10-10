import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react'
import { createMemoryRouter, RouterProvider } from 'react-router-dom'
import { QueryClientProvider } from '@tanstack/react-query'
import { routes } from '@/app/router'
import { queryClient } from '@/app/query-client'
import { useAuthStore } from '@/stores/auth-store'
import { useSessionActionStore } from '@/stores/session-action-store'
import { useUIStore } from '@/stores/ui-store'
import { ApiRequestError } from '@/lib/api-error'
import { borrowingApi } from '@/features/borrowing/api/borrowing.api'
import { profileApi } from '@/features/profile/api/profile.api'
import { reportingApi } from '@/features/reporting/api/reporting.api'
import { termsApi } from '@/features/terms/api/terms.api'
import { authApi } from './api/auth.api'
import { loginDestination } from './navigation'
import type { AuthUser, UserRole } from '@/types/common'

const account = (role: UserRole): AuthUser => ({ id: '00000000-0000-0000-0000-000000000001', email: 'synthetic@example.invalid', name: 'Synthetic Account', role, is_active: true })
function authenticate(role: UserRole) {
  useAuthStore.getState().setSession({ access_token: 'unit-memory-only', expires_at: '2030-01-01T00:00:00Z', token_type: 'Bearer', user: account(role) })
}
async function open(path = '/login', from?: string) {
  const router = createMemoryRouter(routes, { initialEntries: [{ pathname: path, state: from ? { from } : undefined }] })
  const ready = router.state.initialized ? Promise.resolve() : new Promise<void>(resolve => {
    const unsubscribe = router.subscribe(state => { if (state.initialized) { unsubscribe(); resolve() } })
  })
  await act(async () => {
    render(<QueryClientProvider client={queryClient}><RouterProvider router={router} /></QueryClientProvider>)
    await ready
  })
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
  useSessionActionStore.setState({ sessionNotice: null, signingIn: false, logout: 'idle', message: undefined })
  useUIStore.getState().setTheme('light')
  vi.spyOn(termsApi, 'status').mockResolvedValue({ state: 'accepted', current_terms: { id: 'terms-unit-v1', version: 'TEST-1', title: 'SYNTHETIC UNIT TEST TERMS', body: 'TEST ONLY', content_hash: '0'.repeat(64), published_at: '2026-01-01T00:00:00Z' }, acceptance: { id: 'unit-receipt', terms_version_id: 'terms-unit-v1', accepted_at: '2026-01-01T00:00:00Z' }, has_previous_acceptance: true, acceptance_required: false, can_initiate_borrowing: true })
  vi.spyOn(borrowingApi,'list').mockResolvedValue({items:[],total:0,page:1,per_page:3});vi.spyOn(profileApi,'own').mockResolvedValue({name:'Synthetic',email:'synthetic@example.invalid',role:'BORROWER',borrower_type:'FACULTY',student_id:'',course:'',contact_number:''});vi.spyOn(profileApi,'metadata').mockResolvedValue({account_id:'unit',image_id:null,version:0})
  vi.spyOn(reportingApi, 'dashboard').mockResolvedValue({ metrics: { active_loans: 0 }, recent: [], as_of: '2026-10-10T00:00:00Z' })
  vi.spyOn(reportingApi, 'definitions').mockResolvedValue([])
  vi.spyOn(authApi, 'login')
  vi.spyOn(authApi, 'me').mockImplementation(async () => useAuthStore.getState().user ?? account('BORROWER'))
  vi.spyOn(authApi, 'logout').mockResolvedValue(undefined)
})
afterEach(() => { vi.restoreAllMocks(); vi.unstubAllGlobals(); queryClient.clear() })

describe('Phase 4A product authentication', () => {
  it('validates fields, exposes accessible errors, and never submits invalid credentials', async () => {
    await open(); await screen.findByLabelText('Email'); submit()
    expect(await screen.findByText('Enter a valid email address')).toBeInTheDocument()
    expect(screen.getByLabelText('Email')).toHaveAttribute('aria-invalid', 'true')
    expect(screen.getByLabelText('Password')).toHaveAttribute('aria-describedby', 'password-error')
    expect(authApi.login).not.toHaveBeenCalled()
    expect(screen.queryByRole('link', { name: /register|forgot/i })).not.toBeInTheDocument()
  })
  it('prevents duplicate submissions and supports password visibility', async () => {
    const gate = deferred<{ user: AuthUser }>(); vi.mocked(authApi.login).mockReturnValue(gate.promise)
    await open(); await credentials()
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
    await open(); await credentials(); submit()
    expect(await screen.findByText(content => content.includes(String(text)))).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Sign In' })).toBeEnabled()
    expect(screen.getByLabelText('Password')).toHaveValue('')
    expect(screen.getByLabelText('Password')).toHaveFocus()
  })
  it.each([
    ['BORROWER', '/borrower/home'], ['STAFF', '/staff/dashboard'], ['ADMIN', '/staff/dashboard'],
  ] as const)('navigates %s using safe server account metadata', async (role, destination) => {
    vi.mocked(authApi.login).mockImplementation(async () => { authenticate(role); return { user: account(role) } })
    const router = await open(); await credentials(); submit()
    await waitFor(() => expect(router.state.location.pathname).toBe(destination))
    expect(await screen.findByRole('heading', { name: role === 'BORROWER' ? /Hello,/ : 'Dashboard', level: 1 })).toBeInTheDocument()
    expect(authApi.login).toHaveBeenCalledWith({ email: 'synthetic@example.invalid', password: 'unit-only-input' })
  })
  it('preserves an allowed deep link after an anonymous redirect and login', async () => {
    vi.mocked(authApi.login).mockImplementation(async () => { authenticate('BORROWER'); return { user: account('BORROWER') } })
    const router = await open('/borrower/account'); await credentials()
    expect(router.state.location.pathname).toBe('/login'); submit()
    expect(await screen.findByRole('heading', { name: 'Personal Information' })).toBeInTheDocument()
    expect(router.state.location.pathname).toBe('/borrower/account')
  })
  it.each(['https://untrusted.example', '//untrusted.example', '/staff/dashboard', '/borrower/../admin/reports', '/__preview/borrower/home'])('rejects unsafe Borrower return destination %s', destination => {
    expect(loginDestination('BORROWER', destination)).toBe('/borrower/home')
  })
  it.each([
    ['BORROWER', '/staff/dashboard', false], ['STAFF', '/admin/reports', false], ['ADMIN', '/admin/reports', true],
    ['STAFF', '/staff/dashboard', true], ['BORROWER', '/borrower/home', true], ['ADMIN', '/borrower/home', false],
  ] as const)('guards %s access to %s', async (role, path, allowed) => {
    authenticate(role); await open(path)
    const title = path === '/admin/reports' ? 'Reports' : path === '/borrower/home' ? /Hello,/ : 'Dashboard'
    expect(await screen.findByRole('heading', { name: allowed ? title : 'Access denied', level: 1 })).toBeInTheDocument()
    if (!allowed) expect(screen.queryByRole('heading', { name: title, level: 1 })).not.toBeInTheDocument()
  })
  it('withholds all protected content while current database authority is loading', async () => {
    authenticate('ADMIN'); const gate = deferred<AuthUser>(); vi.mocked(authApi.me).mockReturnValue(gate.promise)
    await open('/admin/reports')
    expect(await screen.findByText('Confirming account access…')).toBeInTheDocument()
    expect(screen.queryByRole('heading', { name: 'Reports' })).not.toBeInTheDocument()
    gate.resolve(account('BORROWER'))
    expect(await screen.findByText('Access denied')).toBeInTheDocument()
    expect(useAuthStore.getState().user?.role).toBe('BORROWER')
  })
  it('denies inactive metadata and provides recoverable access errors for transient failures', async () => {
    authenticate('BORROWER'); vi.mocked(authApi.me).mockRejectedValueOnce(new ApiRequestError('Unable to connect.', 0)).mockResolvedValueOnce({ ...account('BORROWER'), is_active: false })
    await open('/borrower/home')
    expect(await screen.findByRole('heading', { name: 'Unable to confirm access' })).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: 'Try again' }))
    expect(await screen.findByText('Access denied')).toBeInTheDocument()
    expect(screen.queryByRole('heading', { name: /Hello,/, level: 1 })).not.toBeInTheDocument()
  })
  it('clears private memory immediately during sign out and permits retry after network failure', async () => {
    authenticate('STAFF'); const logoutRouter=await open('/staff/dashboard'); await screen.findByRole('heading', { name: 'Dashboard', level: 1 })
    queryClient.setQueryData(['users', 'private'], { confidential: true })
    const gate = deferred<void>(); vi.mocked(authApi.logout).mockReturnValueOnce(gate.promise)
    fireEvent.click(screen.getByRole('button',{name:'Account actions'}));await screen.findByRole('menuitem',{name:'Sign Out'})
    fireEvent.click(screen.getByRole('menuitem', { name: 'Sign Out' }))
    expect(await screen.findByText('Signing out…')).toBeInTheDocument()
    expect(useAuthStore.getState().accessToken).toBeNull()
    expect(queryClient.getQueryData(['users', 'private'])).toBeUndefined()
    gate.resolve(); await waitFor(() => expect(useSessionActionStore.getState().logout).toBe('idle'))
    expect(await screen.findByRole('button', { name: 'Sign In' })).toBeEnabled()
    vi.mocked(authApi.logout).mockRejectedValueOnce(new ApiRequestError('Unable to connect.', 0))
    authenticate('BORROWER'); await screen.findByRole('heading', { name: /Hello,/, level: 1 });await act(async()=>{await logoutRouter.navigate('/borrower/account')});await screen.findByRole('button',{name:'Sign Out'})
    fireEvent.click(screen.getByRole('button', { name: 'Sign Out' }))
    expect(await screen.findByText(/Server sign-out could not be confirmed/)).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: 'Retry sign out' }))
    await waitFor(() => expect(useSessionActionStore.getState().logout).toBe('idle'))
    expect(useAuthStore.getState().isAuthenticated).toBe(false)
  })
  it('switches both login and forbidden experiences using the shared theme control', async () => {
    await open(); await screen.findByLabelText('Email')
    fireEvent.click(screen.getByRole('button', { name: 'Switch to dark theme' }))
    expect(document.documentElement).toHaveClass('dark')
  })
})


describe('pre-Phase 7 persistent operational layout', () => {
  it('keeps chrome and UI state while route authority is pending, denied, and recovered', async () => {
    authenticate('ADMIN'); useUIStore.getState().setSidebar(false); useUIStore.getState().setTheme('dark')
    const router = await open('/staff/dashboard')
    await screen.findByRole('heading', { name: 'Dashboard' })
    const sidebar = document.querySelector('.staff-sidebar'), header = document.querySelector('.staff-top-bar'), main = document.querySelector('#staff-content')
    const gate = deferred<AuthUser>(); vi.mocked(authApi.me).mockReturnValueOnce(gate.promise)
    await act(async () => { await router.navigate('/admin/reports') })
    expect(await screen.findByText('Confirming account access…')).toBeInTheDocument()
    expect(screen.queryByRole('heading', { name: 'Reports' })).not.toBeInTheDocument()
    expect(document.querySelector('.staff-sidebar')).toBe(sidebar)
    expect(document.querySelector('.staff-top-bar')).toBe(header)
    expect(document.querySelector('#staff-content')).toBe(main)
    await act(async () => { gate.resolve({ ...account('ADMIN'), is_active: false }) })
    expect(await screen.findByText('Access denied')).toBeInTheDocument()
    expect(screen.queryByRole('heading', { name: 'Reports' })).not.toBeInTheDocument()
    expect(document.querySelector('.staff-sidebar')).toBe(sidebar)
    vi.mocked(authApi.me).mockResolvedValue(account('ADMIN'))
    await act(async () => { await router.navigate('/staff/requests') })
    expect(await screen.findByRole('heading', { name: 'Requests & Borrowings' })).toBeInTheDocument()
    expect(document.querySelector('.staff-top-bar')).toBe(header)
    expect(document.querySelector('.reconstruction-shell')).toHaveAttribute('data-sidebar-open', 'false')
    expect(document.documentElement).toHaveClass('dark')
  })
  it('withholds cached route content during revalidation without replacing chrome', async () => {
    authenticate('ADMIN'); const router = await open('/admin/reports')
    await screen.findByRole('heading', { name: 'Reports' })
    const header = document.querySelector('.staff-top-bar')
    await act(async () => { await router.navigate('/staff/dashboard') })
    await screen.findByRole('heading', { name: 'Dashboard' })
    const gate = deferred<AuthUser>(); vi.mocked(authApi.me).mockReturnValueOnce(gate.promise)
    await act(async () => { await router.navigate('/admin/reports') })
    expect(await screen.findByText('Confirming account access…')).toBeInTheDocument()
    expect(screen.queryByRole('heading', { name: 'Reports' })).not.toBeInTheDocument()
    expect(document.querySelector('.staff-top-bar')).toBe(header)
    await act(async () => { gate.resolve(account('STAFF')) })
    expect(await screen.findByText('Access denied')).toBeInTheDocument()
    expect(screen.queryByRole('heading', { name: 'Reports' })).not.toBeInTheDocument()
  })
  it('keeps the shell during an authority error and successful retry', async () => {
    authenticate('STAFF'); const router = await open('/staff/dashboard')
    await screen.findByRole('heading', { name: 'Dashboard' })
    const sidebar = document.querySelector('.staff-sidebar')
    vi.mocked(authApi.me).mockRejectedValueOnce(new ApiRequestError('Unable to connect.', 0))
    await act(async () => { await router.navigate('/staff/requests') })
    expect(await screen.findByRole('heading', { name: 'Unable to confirm access' })).toBeInTheDocument()
    expect(screen.queryByRole('heading', { name: 'Requests & Borrowings' })).not.toBeInTheDocument()
    expect(document.querySelector('.staff-sidebar')).toBe(sidebar)
    fireEvent.click(screen.getByRole('button', { name: 'Try again' }))
    expect(await screen.findByRole('heading', { name: 'Requests & Borrowings' })).toBeInTheDocument()
    expect(document.querySelector('.staff-sidebar')).toBe(sidebar)
  })
})


describe('culinary login session notice', () => {
  const notice = 'Your session is no longer available. Sign in again to continue.'
  it('keeps a normal anonymous entry neutral even if the legacy ended flag remains', async () => {
    useAuthStore.setState({ sessionEnded: true }); await open('/login')
    await screen.findByLabelText('Email')
    expect(screen.queryByText(notice)).not.toBeInTheDocument()
  })
  it('shows an actual involuntary session invalidation and retains the form', async () => {
    useSessionActionStore.getState().setSessionNotice('invalidated'); await open('/login')
    expect(await screen.findByText(notice)).toHaveAttribute('role', 'status')
    expect(screen.getByLabelText('Email')).toBeEnabled()
    expect(screen.getByLabelText('Password')).toHaveAttribute('autoComplete', 'current-password')
  })
  it('does not describe a recoverable network failure as an expired session', async () => {
    useSessionActionStore.getState().setSessionNotice('invalidated')
    useAuthStore.getState().setStatus('error'); await open('/login')
    await screen.findByLabelText('Email'); expect(screen.queryByText(notice)).not.toBeInTheDocument()
  })
})
