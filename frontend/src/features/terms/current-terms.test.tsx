import type { PropsWithChildren } from 'react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { act, cleanup, renderHook, waitFor } from '@testing-library/react'
import { QueryClientProvider } from '@tanstack/react-query'
import { queryClient } from '@/app/query-client'
import { useAuthStore } from '@/stores/auth-store'
import { ApiRequestError } from '@/lib/api-error'
import { termsApi } from './api/terms.api'
import { useCurrentTerms } from './hooks/use-terms'
import type { TermsVersion } from './types'

const version: TermsVersion = { id: 'unit-version', version: 'TEST', title: 'TEST ONLY', body: 'Synthetic test', content_hash: '0'.repeat(64), published_at: '2026-01-01T00:00:00Z' }
const wrapper = ({ children }: PropsWithChildren) => <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
// Subscribe to data as the consuming page does during render.
const open = () => renderHook(() => ({ ...useCurrentTerms() }), { wrapper })
beforeEach(() => {
 queryClient.clear(); useAuthStore.getState().clear()
 useAuthStore.getState().setSession({ access_token: 'unit-memory-only', token_type: 'Bearer', expires_at: '2030-01-01T00:00:00Z', user: { id: 'unit-admin', role: 'ADMIN', name: 'Synthetic Admin', email: 'unit@example.invalid', is_active: true } })
 vi.spyOn(termsApi, 'current').mockRejectedValue(new ApiRequestError('Unavailable policy.', 503, 'TERMS_NOT_PUBLISHED'))
})
afterEach(() => { cleanup(); vi.restoreAllMocks(); queryClient.clear(); useAuthStore.getState().clear() })

describe('administrative publication information', () => {
 it('caches expected absence across remounts without retrying or implying acceptance', async () => {
  const first = open(); await waitFor(() => expect(first.result.current.isSuccess).toBe(true))
  expect(first.result.current.data).toBeNull(); first.unmount()
  const second = open(); await waitFor(() => expect(second.result.current.isSuccess).toBe(true))
  expect(second.result.current.data).toBeNull(); expect(termsApi.current).toHaveBeenCalledTimes(1)
  expect(queryClient.getQueriesData({ queryKey: ['terms', 'status'] })).toEqual([])
 })
 it('rechecks after the bounded freshness window and can discover publication', async () => {
  const first = open(); await waitFor(() => expect(first.result.current.isSuccess).toBe(true)); first.unmount()
  queryClient.setQueryData(['terms', 'unit-admin', 'current-publication'], null, { updatedAt: Date.now() - 30_001 })
  vi.mocked(termsApi.current).mockResolvedValue(version)
  const next = open(); await waitFor(() => expect(next.result.current.data).toEqual(version))
  expect(termsApi.current).toHaveBeenCalledTimes(2)
 })
 it('explicit terms invalidation discovers a new publication immediately', async () => {
  const current = open(); await waitFor(() => expect(current.result.current.isSuccess).toBe(true))
  vi.mocked(termsApi.current).mockResolvedValue(version)
  await act(async () => { await queryClient.invalidateQueries({ queryKey: ['terms'] }) })
  await waitFor(() => expect(current.result.current.data).toEqual(version))
  expect(termsApi.current).toHaveBeenCalledTimes(2)
 })
 it.each([[0, 'NETWORK_ERROR'], [500, 'INTERNAL_ERROR'], [503, 'SERVICE_UNAVAILABLE'], [401, 'TOKEN_INVALID'], [403, 'FORBIDDEN'], [500, 'TERMS_NOT_PUBLISHED']])('retains genuine failure %s/%s and does not automatically retry', async (status, code) => {
  const error = new ApiRequestError('Safe failure.', Number(status), String(code))
  vi.mocked(termsApi.current).mockRejectedValue(error)
  const current = open(); await waitFor(() => expect(current.result.current.isError).toBe(true))
  expect(current.result.current.error).toBe(error); expect(current.result.current.data).toBeUndefined()
  expect(termsApi.current).toHaveBeenCalledTimes(1)
 })
 it.each(['STAFF', 'BORROWER'] as const)('does not fetch publication for %s', role => {
  useAuthStore.getState().setUser({ ...useAuthStore.getState().user!, role })
  open(); expect(termsApi.current).not.toHaveBeenCalled()
 })
 it('clears cached policy information on logout', async () => {
  const current = open(); await waitFor(() => expect(current.result.current.isSuccess).toBe(true)); current.unmount()
  useAuthStore.getState().clear()
  expect(queryClient.getQueriesData({ queryKey: ['terms'] })).toEqual([])
 })
 it('clears publication information when the account loses Admin authority', async () => {
  const current = open(); await waitFor(() => expect(current.result.current.isSuccess).toBe(true)); current.unmount()
  useAuthStore.getState().setUser({ ...useAuthStore.getState().user!, role: 'STAFF' })
  expect(queryClient.getQueriesData({ queryKey: ['terms'] })).toEqual([])
 })
})
