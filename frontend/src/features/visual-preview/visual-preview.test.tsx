import { useState } from 'react'
import { beforeEach, afterEach, describe, expect, it, vi } from 'vitest'
import { fireEvent, render, screen, within } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { QuantitySelector } from '@/components/application/equipment'
import { BorrowingStatusBadge, ThemeControl } from '@/components/application/visual'
import { StaffSidebar } from '@/components/application/shells'
import { useUIStore } from '@/stores/ui-store'
import { useAuthStore } from '@/stores/auth-store'
import { apiClient } from '@/lib/api-client'
import { EquipmentCatalogPreview } from './borrower-pages'
import { PendingRequestsPreview } from './staff-pages'

beforeEach(() => {
  useUIStore.getState().setTheme('light')
  vi.stubGlobal('matchMedia', vi.fn(() => ({ matches: false, addEventListener: vi.fn(), removeEventListener: vi.fn() })))
})
afterEach(() => { vi.unstubAllGlobals(); vi.restoreAllMocks() })

function QuantityHarness() {
  const [value, setValue] = useState(1)
  return <QuantitySelector name="Preview spoon" value={value} max={2} onChange={setValue} />
}

describe('Phase 3B visual boundaries', () => {
  it('allows local integer selection, disables bounds, and rejects invalid quantities with linked errors', () => {
    render(<QuantityHarness />)
    const input = screen.getByRole('spinbutton', { name: 'Quantity for Preview spoon' })
    fireEvent.click(screen.getByRole('button', { name: 'Add one Preview spoon' }))
    expect(input).toHaveValue(2)
    expect(screen.getByRole('button', { name: 'Add one Preview spoon' })).toBeDisabled()
    fireEvent.change(input, { target: { value: '1.5' } })
    expect(screen.getByRole('alert')).toHaveTextContent('whole number')
    expect(input).toHaveAttribute('aria-invalid', 'true')
    expect(input).toHaveAccessibleDescription(/whole number/)
    expect(input).toHaveValue(2)
    fireEvent.change(input, { target: { value: '0' } })
    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Remove one Preview spoon' })).toBeDisabled()
  })

  it('changes appearance through the existing UI store without changing the session', () => {
    const before = useAuthStore.getState()
    render(<ThemeControl />)
    fireEvent.click(screen.getByRole('button', { name: 'Switch to dark theme' }))
    expect(document.documentElement).toHaveClass('dark')
    expect(useUIStore.getState().theme).toBe('dark')
    expect(useAuthStore.getState()).toBe(before)
  })

  it('presents canonical lifecycle labels without separate approval and release states', () => {
    render(<><BorrowingStatusBadge status="PENDING" /><BorrowingStatusBadge status="CHECKED_OUT" /><BorrowingStatusBadge status="COMPLETED" /></>)
    expect(screen.getByText('Pending')).toBeInTheDocument()
    expect(screen.getByText('Active')).toBeInTheDocument()
    expect(screen.getByText('Completed')).toBeInTheDocument()
    expect(screen.queryByText('Approved')).not.toBeInTheDocument()
    expect(screen.queryByText('Released')).not.toBeInTheDocument()
  })

  it('keeps Admin additions out of Staff navigation and does not grant a session when rendering Admin', () => {
    const before = useAuthStore.getState()
    const props = { active: 'Dashboard', dashboardHref: '/dashboard', requestsHref: '/requests', onUnavailable: vi.fn() }
    const { rerender } = render(<MemoryRouter><StaffSidebar {...props} audience="staff" /></MemoryRouter>)
    expect(screen.queryByRole('button', { name: 'Administration' })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Reports' })).not.toBeInTheDocument()
    rerender(<MemoryRouter><StaffSidebar {...props} audience="admin" /></MemoryRouter>)
    expect(screen.getByRole('button', { name: 'Administration' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Reports' })).toBeInTheDocument()
    expect(useAuthStore.getState()).toBe(before)
  })

  it('filters catalog fixtures, preserves local selection, and offers no submit or payment action', async () => {
    const get = vi.spyOn(apiClient, 'get'), post = vi.spyOn(apiClient, 'post')
    render(<MemoryRouter><EquipmentCatalogPreview /></MemoryRouter>)
    fireEvent.change(screen.getByRole('searchbox', { name: 'Search equipment' }), { target: { value: 'nonexistent' } })
    expect(screen.getByText('No matching equipment')).toBeInTheDocument()
    expect(screen.getByText('2 types · 3 units selected')).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: 'Clear filters' }))
    expect(screen.getByRole('heading', { name: 'Wooden Spoon' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Out of stock' })).toBeDisabled()
    fireEvent.click(screen.getByRole('button', { name: /Review Request/ }))
    const dialog = await screen.findByRole('dialog', { name: 'Selected equipment preview' })
    expect(within(dialog).getByText(/not reserved/)).toBeInTheDocument()
    expect(within(dialog).queryByRole('button', { name: /submit|pay|approve/i })).not.toBeInTheDocument()
    expect(get).not.toHaveBeenCalled()
    expect(post).not.toHaveBeenCalled()
  })

  it('paginates and filters the queue locally; Review is read-only and cannot issue equipment', async () => {
    const post = vi.spyOn(apiClient, 'post')
    render(<MemoryRouter><PendingRequestsPreview /></MemoryRouter>)
    expect(screen.getByText(/Showing 1–5 of 8/)).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: 'Next page' }))
    expect(screen.getByText(/Showing 6–8 of 8/)).toBeInTheDocument()
    fireEvent.change(screen.getByRole('searchbox', { name: 'Search pending requests' }), { target: { value: 'Alex' } })
    expect(screen.getByText(/Showing 1–1 of 1/)).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: 'Review PREVIEW-0215' }))
    const dialog = await screen.findByRole('dialog', { name: 'Request review preview' })
    expect(within(dialog).getByText(/approval and physical release are one transition/)).toBeInTheDocument()
    expect(within(dialog).queryByRole('button', { name: /approve|deny|return|issue/i })).not.toBeInTheDocument()
    expect(post).not.toHaveBeenCalled()
  })
})
