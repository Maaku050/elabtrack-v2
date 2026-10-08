import { StrictMode } from 'react'
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { expect, it, vi } from 'vitest'
import { AxiosError } from 'axios'
import { Providers } from './providers'
import { useAuthStore } from '@/stores/auth-store'

const mock = vi.hoisted(() => ({ calls: 0 }))
vi.mock('@/lib/api-client', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/lib/api-client')>()
  return { ...actual, apiClient: new actual.ApiClient(async (config) => {
    mock.calls++
    if (mock.calls === 1) throw new AxiosError('Network unavailable', 'ERR_NETWORK', config)
    return { config, status: 200, statusText: 'OK', headers: {}, data: { success: true, data: {
      access_token: 'synthetic-memory-access', expires_at: '2026-10-06T12:00:00Z', token_type: 'Bearer',
      user: { id: '00000000-0000-0000-0000-000000000001', email: 'synthetic@example.invalid', name: 'Current Account', role: 'BORROWER', is_active: true },
    } } }
  }) }
})

it('mounts one bootstrap in StrictMode, shows recoverable failure and retries deliberately', async () => {
  await import('@/lib/api-client')
  useAuthStore.getState().clear(); useAuthStore.setState({ status: 'idle' })
  render(<StrictMode><Providers><h1>Public foundation</h1></Providers></StrictMode>)
  expect(screen.getByRole('status')).toHaveTextContent('Restoring session')
  expect(screen.queryByRole('heading')).not.toBeInTheDocument()
  await screen.findByText('Unable to restore your session.')
  expect(screen.getByRole('heading', { name: 'Public foundation' })).toBeInTheDocument()
  expect(mock.calls).toBe(1)
  fireEvent.click(screen.getByRole('button', { name: 'Retry session' }))
  await waitFor(() => expect(useAuthStore.getState().status).toBe('authenticated'))
  expect(mock.calls).toBe(2)
  expect(screen.queryByText('Unable to restore your session.')).not.toBeInTheDocument()
})
