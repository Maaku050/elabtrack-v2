import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { createMemoryRouter, RouterProvider } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { routes } from '@/app/router'
import { useUIStore } from '@/stores/ui-store'

function renderRoute(path = '/') {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  const router = createMemoryRouter(routes, { initialEntries: [path] })
  render(<QueryClientProvider client={client}><RouterProvider router={router} /></QueryClientProvider>)
}

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn())
  useUIStore.getState().setTheme('light')
})
afterEach(() => vi.unstubAllGlobals())

describe('Phase 0 application boundaries', () => {
  it('loads without authentication, applies the theme, and navigates without fetching automatically', async () => {
    renderRoute()
    expect(screen.getByRole('heading', { name: 'eLabTrack V2' })).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: 'Switch to dark theme' }))
    expect(document.documentElement).toHaveClass('dark')
    fireEvent.click(screen.getByRole('link', { name: 'Check service connection' }))
    expect(await screen.findByRole('heading', { name: 'eLabTrack V2 service connection' })).toBeInTheDocument()
    expect(fetch).not.toHaveBeenCalled()
  })

  it('checks the unauthenticated API only on request and shows backend dependency status', async () => {
    vi.mocked(fetch).mockResolvedValue(new Response(JSON.stringify({ status: 'ok', services: { api: 'healthy', database: 'unknown' } })))
    renderRoute('/status')
    expect(fetch).not.toHaveBeenCalled()
    fireEvent.click(screen.getByRole('button', { name: 'Check connection' }))
    expect(await screen.findByText('API connected.')).toBeInTheDocument()
    expect(screen.getByText('database: unknown')).toBeInTheDocument()
    expect(fetch).toHaveBeenCalledWith(expect.stringMatching(/\/api\/v1\/health$/), expect.objectContaining({ credentials: 'omit', signal: expect.any(AbortSignal) }))
  })

  it('handles backend failure and a manual retry without exposing internal responses', async () => {
    vi.mocked(fetch).mockResolvedValueOnce(new Response('private database detail', { status: 503 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ status: 'ok', services: { api: 'healthy', database: 'healthy' } })))
    renderRoute('/status')
    fireEvent.click(screen.getByRole('button', { name: 'Check connection' }))
    expect(await screen.findByText(/Unable to connect/)).toBeInTheDocument()
    expect(screen.queryByText('private database detail')).not.toBeInTheDocument()
    await waitFor(() => expect(screen.getByRole('button', { name: 'Check connection' })).toBeEnabled())
    fireEvent.click(screen.getByRole('button', { name: 'Check connection' }))
    expect(await screen.findByText('API connected.')).toBeInTheDocument()
  })

  it('rejects malformed health payloads', async () => {
    vi.mocked(fetch).mockResolvedValue(new Response(JSON.stringify({ success: true, data: {} })))
    renderRoute('/status')
    fireEvent.click(screen.getByRole('button', { name: 'Check connection' }))
    expect(await screen.findByText(/Unable to connect/)).toBeInTheDocument()
    expect(screen.queryByText('API connected.')).not.toBeInTheDocument()
  })

  it.each(['/register', '/app/dashboard', '/app/components', '/missing'])('keeps removed starter route %s outside the application', (path) => {
    renderRoute(path)
    expect(screen.getByRole('heading', { name: /Page not found/ })).toBeInTheDocument()
    expect(fetch).not.toHaveBeenCalled()
  })
})
