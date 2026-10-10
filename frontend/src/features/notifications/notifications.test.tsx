import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter } from 'react-router-dom'
import { afterEach, beforeEach, expect, it, vi } from 'vitest'
import { NotificationCenter } from './pages/notification-center'
import { NotificationBell } from './components/notification-bell'
import { notificationsApi } from './api/notifications.api'
import { useAuthStore } from '@/stores/auth-store'
const notice = { id: 'notification', borrowing_id: 'loan', kind: 'SUBMITTED', title: 'Request submitted', body: 'Your request is pending.', scope: 'BORROWER', created_at: '2026-10-10T00:00:00Z', read_at: null } as const
async function mount(element: React.ReactNode) { const client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); await act(async () => { render(<QueryClientProvider client={client}><MemoryRouter>{element}</MemoryRouter></QueryClientProvider>) }) }
beforeEach(() => { useAuthStore.setState({ user: { id: 'borrower', email: 'fictional@example.invalid', name: 'TEST Faculty', role: 'BORROWER', is_active: true }, isAuthenticated: true }); vi.spyOn(notificationsApi, 'count').mockResolvedValue({ unread: 1 }); vi.spyOn(notificationsApi, 'list').mockResolvedValue({ items: [notice], total: 1, unread: 1, page: 1, per_page: 25 }) })
afterEach(() => { cleanup(); vi.restoreAllMocks(); useAuthStore.setState({ user: null, isAuthenticated: false }) })
it('uses borrower links, persisted read updates and server filters', async () => { const mark = vi.spyOn(notificationsApi, 'mark').mockResolvedValue({ ...notice, read_at: '2026-10-10T01:00:00Z' }); await mount(<NotificationCenter />); expect(await screen.findByRole('link', { name: 'Open Borrowing' })).toHaveAttribute('href', '/borrower/borrowings/loan'); fireEvent.click(screen.getByRole('button', { name: 'Mark read: Request submitted' })); await waitFor(() => expect(mark).toHaveBeenCalledWith('notification', true)); fireEvent.change(screen.getByLabelText('Read state'), { target: { value: 'unread' } }); await waitFor(() => expect(notificationsApi.list).toHaveBeenCalledWith({ page: 1, per_page: 25, read: 'unread' }, expect.any(AbortSignal))) })
it('uses operational borrowing links for Staff without exposing recipient selectors', async () => { useAuthStore.setState({ user: { ...useAuthStore.getState().user!, role: 'STAFF' } }); await mount(<NotificationCenter />); expect(await screen.findByRole('link', { name: 'Open Borrowing' })).toHaveAttribute('href', '/staff/requests/loan'); expect(screen.queryByLabelText(/recipient/i)).not.toBeInTheDocument() })
it('does not turn an unavailable count into a false zero', async () => { vi.mocked(notificationsApi.count).mockRejectedValue(new Error('unavailable')); await mount(<NotificationBell />); await waitFor(() => expect(screen.getByRole('link', { name: 'Notifications' })).toHaveAttribute('aria-description', 'Unread count unavailable')) })
it('caps the visible count while retaining the actual accessible unread count', async () => { vi.mocked(notificationsApi.count).mockResolvedValue({ unread: 123 }); await mount(<NotificationBell />); expect(await screen.findByText('99+')).toBeInTheDocument(); expect(screen.getByRole('link', { name: 'Notifications' })).toHaveAttribute('aria-description', '123 unread') })

it('marks all through one durable command, updates rows/count/bell, and prevents pending duplicates',async()=>{
 let finish!:(value:{updated:number})=>void
 const all=vi.spyOn(notificationsApi,'markAll').mockImplementation(()=>new Promise(r=>{finish=r}))
 await mount(<><NotificationBell/><NotificationCenter/></>)
 const button=await screen.findByRole('button',{name:'Mark All Read'})
 fireEvent.click(button);fireEvent.click(button)
 await waitFor(()=>expect(all).toHaveBeenCalledTimes(1));expect(button).toBeDisabled()
 vi.mocked(notificationsApi.list).mockResolvedValue({items:[{...notice,read_at:'2026-10-10T02:00:00Z'}],total:1,unread:0,page:1,per_page:25})
 vi.mocked(notificationsApi.count).mockResolvedValue({unread:0})
 act(()=>finish({updated:1}));await screen.findByText('1 notifications marked read.')
 await waitFor(()=>expect(screen.getByRole('button',{name:'Mark All Read'})).toBeDisabled())
 expect(screen.getByRole('button',{name:'Mark unread: Request submitted'})).toBeEnabled()
 await waitFor(()=>expect(screen.getAllByRole('link',{name:'Notifications'}).find(e=>e.classList.contains('notification-bell'))).toHaveAttribute('aria-description','0 unread'))
})
it('does not submit bulk reads for an empty unread state',async()=>{
 vi.mocked(notificationsApi.list).mockResolvedValue({items:[],total:0,unread:0,page:1,per_page:25})
 const all=vi.spyOn(notificationsApi,'markAll');await mount(<NotificationCenter/>)
 expect(await screen.findByRole('button',{name:'Mark All Read'})).toBeDisabled();expect(all).not.toHaveBeenCalled()
})
it('retains actual unread state and permits retry after a bulk-read failure',async()=>{
 const all=vi.spyOn(notificationsApi,'markAll').mockRejectedValueOnce(new Error('Network unavailable'))
 await mount(<NotificationCenter/>);fireEvent.click(await screen.findByRole('button',{name:'Mark All Read'}))
 await screen.findByRole('alert');await waitFor(()=>expect(screen.getByRole('button',{name:'Mark All Read'})).toBeEnabled())
 expect(screen.getByText(/1 unread/)).toBeInTheDocument();expect(all).toHaveBeenCalledTimes(1)
})
