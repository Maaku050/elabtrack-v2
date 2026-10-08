import { beforeEach, afterEach, describe, expect, it, vi } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { createMemoryRouter, RouterProvider } from 'react-router-dom'
import { QueryClientProvider } from '@tanstack/react-query'
import { routes } from '@/app/router'
import { queryClient } from '@/app/query-client'
import { authApi } from '@/features/auth/api/auth.api'
import { useAuthStore } from '@/stores/auth-store'
import { useSessionActionStore } from '@/stores/session-action-store'
import { useUIStore } from '@/stores/ui-store'
import { ApiRequestError } from '@/lib/api-error'
import { termsApi } from './api/terms.api'
import type { TermsStatus, TermsVersion, TermsAcceptance } from './types'

const id = '00000000-0000-4000-8000-000000000001'
const v1: TermsVersion = { id, version: 'TEST-1', title: 'SYNTHETIC TEST TERMS', body: 'TEST ONLY. This is not institutional policy.\nLiteral <script> text is readable, never executed.', content_hash: '0'.repeat(64), published_at: '2026-01-01T00:00:00Z' }
const receipt: TermsAcceptance = { id: 'unit-receipt', terms_version_id: id, accepted_at: '2026-01-02T00:00:00Z' }
let status: TermsStatus
function accepted(): TermsStatus { return { ...status, state: 'accepted', acceptance: receipt, has_previous_acceptance: true, acceptance_required: false, can_initiate_borrowing: true } }
function open(path = '/borrower/terms') {
 const router = createMemoryRouter(routes, { initialEntries: [path] })
 render(<QueryClientProvider client={queryClient}><RouterProvider router={router} /></QueryClientProvider>)
 return router
}
async function checkConsent() { fireEvent.click(await screen.findByRole('checkbox', { name: 'I have read and accept Version TEST-1.' })) }
function deferred<T>() { let resolve!: (value: T) => void; const promise = new Promise<T>(r => { resolve = r }); return { promise, resolve } }

beforeEach(() => {
 vi.stubGlobal('matchMedia', vi.fn(() => ({ matches: false, addEventListener: vi.fn(), removeEventListener: vi.fn() })))
 queryClient.clear(); useAuthStore.getState().clear()
 useAuthStore.getState().setSession({ access_token: 'unit-memory-only', token_type: 'Bearer', expires_at: '2030-01-01T00:00:00Z', user: { id: 'unit-account', role: 'BORROWER', name: 'Synthetic Borrower', email: 'unit@example.invalid', is_active: true } })
 useSessionActionStore.setState({ signingIn: false, logout: 'idle', message: undefined }); useUIStore.getState().setTheme('light')
 status = { state: 'required', current_terms: v1, acceptance: null, has_previous_acceptance: false, acceptance_required: true, can_initiate_borrowing: false }
 vi.spyOn(authApi, 'me').mockImplementation(async () => useAuthStore.getState().user!)
 vi.spyOn(authApi, 'logout').mockResolvedValue(undefined)
 vi.spyOn(termsApi, 'status').mockImplementation(async () => status)
 vi.spyOn(termsApi, 'accept').mockImplementation(async () => { status = accepted(); return receipt })
})
afterEach(() => { vi.restoreAllMocks(); vi.unstubAllGlobals(); queryClient.clear() })

describe('Phase 4B authoritative terms experience', () => {
 it('redirects first-use Home to complete terms without consent inferred from opening', async () => {
  const router = open('/borrower/home')
  expect(await screen.findByRole('heading', { name: 'Review borrowing terms' }, { timeout: 5000 })).toBeInTheDocument()
  expect(router.state.location.pathname).toBe('/borrower/terms')
  expect(screen.getByRole('checkbox')).not.toBeChecked()
  expect(screen.getByRole('button', { name: 'Accept and continue' })).toBeDisabled()
  expect(termsApi.accept).not.toHaveBeenCalled()
 })
 it('withholds workspace content while acceptance status is loading', async () => {
  const gate = deferred<TermsStatus>(); vi.mocked(termsApi.status).mockReturnValue(gate.promise)
  open('/borrower/home'); expect(await screen.findByText('Checking borrowing terms…')).toBeInTheDocument()
  expect(screen.queryByText('This feature is not available yet')).not.toBeInTheDocument()
  gate.resolve(status); await screen.findByRole('checkbox')
 })
 it('permits returning current-terms borrower to enter Home', async () => {
  status = accepted(); open('/borrower/home')
  expect(await screen.findByText('This feature is not available yet')).toBeInTheDocument()
  expect(screen.queryByRole('checkbox')).not.toBeInTheDocument()
 })
 it('does not treat malformed accepted metadata as current consent', async () => {
  status = { ...accepted(), acceptance: { ...receipt, terms_version_id: 'different-version' } };open('/borrower/home')
  await screen.findByRole('checkbox');expect(screen.queryByText('This feature is not available yet')).not.toBeInTheDocument()
 })
 it('explains updated terms without discarding historical acceptance', async () => {
  status = { ...status, state: 'updated', has_previous_acceptance: true };open()
  expect(await screen.findByRole('heading', { name: 'Updated terms' })).toBeInTheDocument()
  expect(screen.getByText(/The borrowing terms have been updated/)).toBeInTheDocument()
  expect(screen.getByRole('checkbox')).not.toBeChecked()
 })
 it('records only the reviewed version, prevents pending duplicates and proceeds after server success', async () => {
  const gate = deferred<TermsAcceptance>();vi.mocked(termsApi.accept).mockImplementation(async () => { await gate.promise; status = accepted();return receipt })
  const router = open();await checkConsent();fireEvent.click(screen.getByRole('button', { name: 'Accept and continue' }))
  const pending = await screen.findByRole('button', { name: 'Recording acceptance…' });expect(pending).toBeDisabled();fireEvent.click(pending)
  expect(termsApi.accept).toHaveBeenCalledTimes(1);expect(termsApi.accept).toHaveBeenCalledWith(id)
  expect(router.state.location.pathname).toBe('/borrower/terms')
  gate.resolve(receipt);await waitFor(() => expect(router.state.location.pathname).toBe('/borrower/home'))
  await screen.findByText('This feature is not available yet')
 })
 it.each([0, 500])('does not infer consent after failed acceptance %s', async code => {
  vi.mocked(termsApi.accept).mockRejectedValue(new ApiRequestError('Something went wrong. Please try again later.', code, 'INTERNAL_ERROR', undefined, id))
  open();await checkConsent();fireEvent.click(screen.getByRole('button', { name: 'Accept and continue' }))
  expect(await screen.findByRole('alert')).toHaveTextContent(code === 0 ? 'Your acceptance could not be confirmed' : `Reference: ${id}`)
  expect(screen.getByRole('button', { name: 'Accept and continue' })).toBeEnabled()
  expect(useAuthStore.getState().isAuthenticated).toBe(true)
 })
 it('clears checked consent when the server reports a different published version', async () => {
  const v2 = { ...v1, id: '00000000-0000-4000-8000-000000000002', version: 'TEST-2' }
  vi.mocked(termsApi.accept).mockImplementation(async () => { status = { ...status, state: 'updated', current_terms: v2 };throw new ApiRequestError('Terms changed.',409,'TERMS_VERSION_CHANGED') })
  open();await checkConsent();fireEvent.click(screen.getByRole('button', { name: 'Accept and continue' }))
  expect(await screen.findByRole('checkbox', { name: 'I have read and accept Version TEST-2.' })).not.toBeChecked()
  expect(screen.getByRole('button', { name: 'Accept and continue' })).toBeDisabled()
  expect(screen.getByRole('alert')).toHaveTextContent('Read the new version')
  expect(termsApi.accept).toHaveBeenCalledWith(id)
 })
 it('shows current accepted evidence on deliberate terms review', async () => {
  status = accepted();open();expect(await screen.findByText('Accepted', { exact: true })).toBeInTheDocument()
  expect(screen.getByText(/Acceptance recorded on/)).toBeInTheDocument();expect(screen.queryByRole('checkbox')).not.toBeInTheDocument()
 })
 it('handles no official terms safely without an acceptance control', async () => {
  status = { ...status, state: 'unpublished', current_terms: null };open('/borrower/home')
  expect(await screen.findByRole('heading', { name: 'Official terms are not available yet' })).toBeInTheDocument()
  expect(screen.queryByRole('checkbox')).not.toBeInTheDocument()
  expect(screen.getByRole('link', { name: 'View account' })).toHaveAttribute('href','/borrower/account')
  expect(screen.getByRole('link', { name: 'View existing borrowings' })).toBeInTheDocument();expect(screen.getByRole('button', { name: 'Sign Out' })).toBeEnabled()
 })
 it.each(['/borrower/account', '/borrower/borrowings', '/borrower/notifications'])('keeps permitted read-only destination %s accessible without acceptance', async path => {
  status = { ...status, state: 'unpublished', current_terms: null };const router = open(path)
  expect(await screen.findByRole('heading', { name: path.endsWith('account') ? 'Account' : path.endsWith('borrowings') ? 'My Borrowings' : 'Notifications', level: 1 })).toBeInTheDocument()
  expect(router.state.location.pathname).toBe(path)
 })
 it('offers bounded manual recovery for terms network failure', async () => {
  vi.mocked(termsApi.status).mockRejectedValueOnce(new ApiRequestError('Unable to connect.',0)).mockImplementation(async () => status)
  open();expect(await screen.findByRole('heading', { name: 'Unable to load borrowing terms' })).toBeInTheDocument()
  expect(termsApi.status).toHaveBeenCalledTimes(1)
  fireEvent.click(screen.getByRole('button', { name: 'Try again' }));await screen.findByRole('checkbox');expect(termsApi.status).toHaveBeenCalledTimes(2)
 })
 it('logs out without accepting and clears private terms state', async () => {
  open();await screen.findByRole('checkbox');fireEvent.click(screen.getByRole('button', { name: 'Sign Out' }))
  expect(await screen.findByLabelText('Email')).toBeInTheDocument();expect(useAuthStore.getState().accessToken).toBeNull()
  expect(queryClient.getQueriesData({ queryKey: ['terms'] })).toEqual([]);expect(termsApi.accept).not.toHaveBeenCalled()
 })
 it('does not resurrect a session when it expires during acceptance', async () => {
  vi.mocked(termsApi.accept).mockImplementation(async () => { useAuthStore.getState().clear();throw new ApiRequestError('Session expired.',401,'TOKEN_INVALID') })
  open();await checkConsent();fireEvent.click(screen.getByRole('button', { name: 'Accept and continue' }))
  expect(await screen.findByLabelText('Email')).toBeInTheDocument();expect(useAuthStore.getState().isAuthenticated).toBe(false)
 })
 it('denies inactive accounts before displaying terms or submitting consent', async () => {
  vi.mocked(authApi.me).mockResolvedValue({ ...useAuthStore.getState().user!, is_active: false });open()
  expect(await screen.findByText('Access denied')).toBeInTheDocument();expect(termsApi.status).not.toHaveBeenCalled()
 })
 it.each(['STAFF','ADMIN'] as const)('leaves %s navigation unaffected without borrower acceptance', async role => {
  useAuthStore.setState({ user: { ...useAuthStore.getState().user!,role } });open('/staff/dashboard')
  await screen.findByText('This feature is not available yet');expect(termsApi.status).not.toHaveBeenCalled()
 })
 it('renders the complete long document as text and supports both themes', async () => {
  status = { ...status, current_terms: { ...v1, body: v1.body + '\n' + 'Long synthetic test paragraph.\n'.repeat(300) + 'END OF TEST DOCUMENT' } }
  open();const article = await screen.findByRole('article', { name: 'Complete terms, Version TEST-1' })
  expect(article).toHaveTextContent('END OF TEST DOCUMENT');expect(article.querySelector('script')).toBeNull();expect(screen.queryByRole('dialog')).not.toBeInTheDocument()
  fireEvent.click(screen.getByRole('button', { name: 'Switch to dark theme' }));expect(useUIStore.getState().theme).toBe('dark')
  fireEvent.click(screen.getByRole('button', { name: 'Switch to light theme' }));expect(useUIStore.getState().theme).toBe('light')
 })
})
