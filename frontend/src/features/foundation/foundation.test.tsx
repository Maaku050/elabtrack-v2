import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { createMemoryRouter, RouterProvider } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { routes } from '@/app/router'
import { useAuthStore } from '@/stores/auth-store'
import { useUIStore } from '@/stores/ui-store'
import { apiClient } from '@/lib/api-client'
import { ApiRequestError } from '@/lib/api-error'

function renderRoute(path = '/') {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  const router = createMemoryRouter(routes, { initialEntries: [path] })
  render(<QueryClientProvider client={client}><RouterProvider router={router} /></QueryClientProvider>)
}

beforeEach(() => {
  useAuthStore.getState().clear()
  vi.spyOn(apiClient, 'get')
  useUIStore.getState().setTheme('light')
})
afterEach(() => vi.restoreAllMocks())

describe('Phase 0 application boundaries', () => {
  it('loads the anonymous login and applies the theme without fetching protected data', async () => {
    renderRoute()
    expect(await screen.findByRole('heading', { name: 'Sign in to your account' })).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: 'Switch to dark theme' }))
    expect(document.documentElement).toHaveClass('dark')
    expect(apiClient.get).not.toHaveBeenCalled()
  })

  it('checks the unauthenticated API only on request and shows minimal process liveness', async () => {
    vi.mocked(apiClient.get).mockResolvedValue({ status: 'ok', service: 'elabtrack-v2' })
    renderRoute('/status')
    expect(apiClient.get).not.toHaveBeenCalled()
    fireEvent.click(screen.getByRole('button', { name: 'Check connection' }))
    expect(await screen.findByText('API connected.')).toBeInTheDocument()
    expect(screen.getByText('elabtrack-v2')).toBeInTheDocument()
    expect(screen.getByText(/checks API liveness/)).toBeInTheDocument()
    expect(apiClient.get).toHaveBeenCalledWith('/health', expect.objectContaining({ withCredentials: false, attachAuth: false, retryAuth: false, signal: expect.any(AbortSignal) }))
  })

  it('handles backend failure and a manual retry without exposing internal responses', async () => {
    vi.mocked(apiClient.get).mockRejectedValueOnce(new ApiRequestError('Something went wrong.', 503, 'SERVICE_UNAVAILABLE'))
      .mockResolvedValueOnce({ status: 'ok', service: 'elabtrack-v2' })
    renderRoute('/status')
    fireEvent.click(screen.getByRole('button', { name: 'Check connection' }))
    expect(await screen.findByText(/Unable to connect/)).toBeInTheDocument()
    expect(screen.queryByText('private database detail')).not.toBeInTheDocument()
    await waitFor(() => expect(screen.getByRole('button', { name: 'Check connection' })).toBeEnabled())
    fireEvent.click(screen.getByRole('button', { name: 'Check connection' }))
    expect(await screen.findByText('API connected.')).toBeInTheDocument()
  })

  it('rejects malformed health payloads', async () => {
    vi.mocked(apiClient.get).mockResolvedValue({})
    renderRoute('/status')
    fireEvent.click(screen.getByRole('button', { name: 'Check connection' }))
    expect(await screen.findByText(/Unable to connect/)).toBeInTheDocument()
    expect(screen.queryByText('API connected.')).not.toBeInTheDocument()
  })

  it.each(['/register', '/app/dashboard', '/app/components', '/missing'])('keeps removed starter route %s outside the application', (path) => {
    renderRoute(path)
    expect(screen.getByRole('heading', { name: /Page not found/ })).toBeInTheDocument()
    expect(apiClient.get).not.toHaveBeenCalled()
  })
})
