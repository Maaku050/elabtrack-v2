import { render, screen, fireEvent } from '@testing-library/react'
import { MemoryRouter, useLocation } from 'react-router-dom'
import { describe, it, expect, vi } from 'vitest'
import { ServerPagination } from './server-pagination'
import { pageRange } from '@/lib/page-range'
import { useDirectoryState } from '@/hooks/use-directory-state'

function mount(props: Partial<Parameters<typeof ServerPagination>[0]> = {}) {
  const onPage = vi.fn()
  render(<MemoryRouter initialEntries={['/staff/borrowers?status=ACTIVE&page=2']}><ServerPagination label="Borrowers" page={2} totalPages={10} total={250} count={25} onPage={onPage} {...props} /></MemoryRouter>)
  return onPage
}
describe('server directory pagination', () => {
  it('bounds links for very large totals and includes both endpoints', () => {
    expect(pageRange(500, 1000)).toEqual([1, 'ellipsis', 499, 500, 501, 'ellipsis', 1000])
    expect(pageRange(1, 1)).toEqual([1])
    expect(pageRange(2, 3)).toEqual([1, 2, 3])
    expect(pageRange(999, 1000)).toEqual([1, 'ellipsis', 997, 998, 999, 1000])
  })
  it('preserves filters in links, exposes current page and uses SPA callbacks', () => {
    const onPage = mount()
    expect(screen.getByRole('link', { name: 'Go to page 2' })).toHaveAttribute('aria-current', 'page')
    const next = screen.getByRole('link', { name: 'Go to next page' })
    expect(next).toHaveAttribute('href', '/staff/borrowers?status=ACTIVE&page=3')
    fireEvent.click(next); expect(onPage).toHaveBeenCalledWith(3)
    expect(screen.getByRole('status')).toHaveTextContent('26–50 of 250 borrowers')
  })
  it('disables both unavailable directions and excludes them from keyboard navigation', () => {
    const onPage = mount({ page: 1, totalPages: 1, total: 0, count: 0 })
    for (const label of ['Go to previous page', 'Go to next page']) {
      const control = screen.getByLabelText(label)
      expect(control).toHaveAttribute('aria-disabled', 'true'); expect(control).toHaveAttribute('tabindex', '-1'); expect(control).not.toHaveAttribute('href'); fireEvent.click(control)
    }
    expect(onPage).not.toHaveBeenCalled(); expect(screen.getByRole('status')).toHaveTextContent('0–0 of 0 borrowers')
  })
  it('does not invent a final page for unknown category totals', () => {
    const onPage = mount({ total: undefined, totalPages: undefined, label: 'Categories', pageKey: 'category_page', page: 1, hasNext: true, count: 100 })
    expect(screen.getByRole('status')).toHaveTextContent('Total pages unavailable')
    expect(screen.queryByRole('link', { name: 'Go to page 1' })).not.toBeInTheDocument()
    fireEvent.click(screen.getByRole('link', { name: 'Go to next page' })); expect(onPage).toHaveBeenCalledWith(2)
  })
  it('prevents changing pages while the authoritative page is fetching', () => {
    const onPage = mount({ pending: true })
    fireEvent.click(screen.getByLabelText('Go to next page')); fireEvent.click(screen.getByLabelText('Go to page 3'))
    expect(onPage).not.toHaveBeenCalled()
  })
})
function URLHarness() {
  const { value, page, update } = useDirectoryState(), location = useLocation()
  return <><output>{`${page()}|${value('search')}|${location.search}`}</output><button onClick={() => update({ search: '1' })}>Search</button><button onClick={() => update({ page: '1' }, false)}>First</button></>
}
describe('URL directory state', () => {
  it('resets page for filters while keeping numeric-looking search text', () => {
    render(<MemoryRouter initialEntries={['/staff/borrowers?page=3&status=ACTIVE']}><URLHarness /></MemoryRouter>)
    fireEvent.click(screen.getByText('Search')); expect(screen.getByRole('status')).toHaveTextContent('1|1|?status=ACTIVE&search=1')
  })
  it('removes only the default page parameter, retaining other view state', () => {
    render(<MemoryRouter initialEntries={['/staff/borrowers?page=3&search=1']}><URLHarness /></MemoryRouter>)
    fireEvent.click(screen.getByText('First')); expect(screen.getByRole('status')).toHaveTextContent('1|1|?search=1')
  })
  it.each(['-1', 'NaN', '1.5', '1000001'])('rejects invalid URL page %s', invalid => {
    render(<MemoryRouter initialEntries={['/staff/borrowers?page=' + invalid]}><URLHarness /></MemoryRouter>)
    expect(screen.getByRole('status').textContent).toMatch(/^1\|/)
  })
})
