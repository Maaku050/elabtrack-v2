import { render, screen, fireEvent, waitFor, cleanup } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter, useLocation } from 'react-router-dom'
import { beforeEach, afterEach, expect, it, vi } from 'vitest'
import { InventoryDirectory } from './inventory-directory'
import { inventoryApi } from '../api/inventory.api'
import { useAuthStore } from '@/stores/auth-store'
import type { EquipmentPage } from '../types'

const empty: EquipmentPage = { items: [], total: 0, page: 1, per_page: 25, totals: { available: 0, reserved: 0, checked_out: 0, damaged_held: 0, total_tracked: 0 } }
function Location() { return <output aria-label="Directory URL">{useLocation().search}</output> }
function mount(url = '/staff/inventory') {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  render(<QueryClientProvider client={client}><MemoryRouter initialEntries={[url]}><InventoryDirectory /><Location /></MemoryRouter></QueryClientProvider>)
  return client
}
beforeEach(() => {
  useAuthStore.setState({ user: { id: 'test-admin', email: 'test@example.invalid', name: 'TEST Admin', role: 'ADMIN', is_active: true } })
  vi.spyOn(inventoryApi, 'categories').mockResolvedValue([{ id: 'category-a', name: 'Utensils', is_active: true, version: 1 }])
  vi.spyOn(inventoryApi, 'list').mockImplementation(async filter => ({ ...empty, page: filter.page }))
})
afterEach(() => { cleanup(); vi.restoreAllMocks(); useAuthStore.setState({ user: null }) })
it('uses the borrowing label and sends availability independently of the selected status/category', async () => {
  mount('/staff/inventory?status=INACTIVE&category_id=category-a&page=2&sort=available')
  await waitFor(() => expect(inventoryApi.list).toHaveBeenCalledWith(expect.objectContaining({ page: 2, status: 'INACTIVE', category_id: 'category-a', available_only: false, sort: 'available' }), expect.any(AbortSignal)))
  fireEvent.click(screen.getByRole('checkbox', { name: 'Available for borrowing' }))
  await waitFor(() => expect(inventoryApi.list).toHaveBeenLastCalledWith(expect.objectContaining({ page: 1, status: 'INACTIVE', category_id: 'category-a', available_only: true, sort: 'available' }), expect.any(AbortSignal)))
  expect(screen.getByLabelText('Status')).toHaveValue('INACTIVE')
  expect(screen.getByLabelText('Category')).toHaveValue('category-a')
  expect(await screen.findByText('0 matching equipment types')).toBeInTheDocument()
  expect(screen.getByLabelText('Directory URL')).toHaveTextContent('available_only=true')
  fireEvent.click(screen.getByRole('checkbox', { name: 'Available for borrowing' }))
  await waitFor(() => expect(inventoryApi.list).toHaveBeenLastCalledWith(expect.objectContaining({ status: 'INACTIVE', category_id: 'category-a', available_only: false }), expect.any(AbortSignal)))
})
it('resets pagination when switching filters and clears all independent inventory filters', async () => {
  mount('/staff/inventory?page=3&search=Spoon&status=ACTIVE&category_id=category-a&available_only=true&sort=available')
  await waitFor(() => expect(inventoryApi.list).toHaveBeenCalled())
  fireEvent.change(screen.getByLabelText('Status'), { target: { value: 'ARCHIVED' } })
  await waitFor(() => expect(inventoryApi.list).toHaveBeenLastCalledWith(expect.objectContaining({ page: 1, status: 'ARCHIVED', search: 'Spoon', available_only: true }), expect.any(AbortSignal)))
  fireEvent.click(screen.getByRole('button', { name: 'Clear filters' }))
  await waitFor(() => expect(inventoryApi.list).toHaveBeenLastCalledWith({ page: 1, per_page: 25, search: '', category_id: undefined, status: undefined, available_only: false, sort: 'name' }, expect.any(AbortSignal)))
  expect(screen.getByLabelText('Directory URL')).toBeEmptyDOMElement()
  expect(screen.getByRole('checkbox', { name: 'Available for borrowing' })).not.toBeChecked()
  expect(screen.getByLabelText('Status')).toHaveValue('')
})
it('intersects search with availability and preserves the checkbox while typing', async () => {
  mount('/staff/inventory?available_only=true&page=2')
  await waitFor(() => expect(inventoryApi.list).toHaveBeenCalled())
  fireEvent.change(screen.getByRole('searchbox', { name: 'Search equipment' }), { target: { value: 'Spoon' } })
  await waitFor(() => expect(inventoryApi.list).toHaveBeenLastCalledWith(expect.objectContaining({ search: 'Spoon', available_only: true, page: 1 }), expect.any(AbortSignal)))
  expect(screen.getByRole('checkbox', { name: 'Available for borrowing' })).toBeChecked()
})
it('restores filter/page state on remount and displays the complete server page/count without client filtering', async () => {
  const url = '/staff/inventory?available_only=true&status=ACTIVE&search=Spoon&category_id=category-a&page=2'
  // Deliberately retain a server row with zero stock: the frontend must render
  // authoritative pages, rather than hiding defects in a page-local filter.
  vi.mocked(inventoryApi.list).mockResolvedValue({ ...empty, total: 27, page: 2, items: [{ id: 'server-row', name: 'Server page Spoon', description: '', category_id: null, category_name: '', status: 'ACTIVE', stock: empty.totals, stock_sequence: 0, metadata_version: 1, image_id: null, created_at: '', updated_at: '' }] })
  let client = mount(url)
  expect(await screen.findByRole('link', { name: 'Server page Spoon' })).toBeInTheDocument()
  expect(screen.getByText('27 matching equipment types')).toBeInTheDocument()
  expect(screen.getByLabelText('Status')).toHaveValue('ACTIVE')
  expect(screen.getByRole('checkbox', { name: 'Available for borrowing' })).toBeChecked()
  expect(screen.getByRole('link', { name: 'Go to previous page' })).toHaveAttribute('href', expect.stringContaining('available_only=true'))
  cleanup(); client.clear(); vi.mocked(inventoryApi.list).mockClear()
  client = mount(url)
  await waitFor(() => expect(inventoryApi.list).toHaveBeenCalledWith(expect.objectContaining({ page: 2, search: 'Spoon', status: 'ACTIVE', category_id: 'category-a', available_only: true }), expect.any(AbortSignal)))
  expect(await screen.findByText('27 matching equipment types')).toBeInTheDocument()
  client.clear()
})
