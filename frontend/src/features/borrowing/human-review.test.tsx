import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter, useLocation } from 'react-router-dom'
import { afterEach, beforeEach, expect, it, vi } from 'vitest'
import { BorrowingDirectory } from './pages/borrowing-directory'
import { CatalogCart } from './pages/catalog-cart'
import { WorkspacePage } from '@/features/workspace/workspace-page'
import { useAuthStore } from '@/stores/auth-store'
import { useCartStore } from './stores/cart-store'
import { borrowingApi } from './api/borrowing.api'
import { inventoryApi } from '@/features/inventory/api/inventory.api'
import { profileApi } from '@/features/profile/api/profile.api'
import { termsApi } from '@/features/terms/api/terms.api'
import type { Borrowing } from './types'
import type { Equipment } from '@/features/inventory/types'
const logout=vi.fn()
vi.mock('@/features/auth/hooks/use-auth',()=>({useLogout:()=>logout}))
vi.mock('@/features/notifications/components/notification-bell',()=>({NotificationBell:()=>null}))
const equipment:Equipment={id:'11111111-1111-4111-8111-111111111111',name:'TEST Ladle',description:'',category_id:null,category_name:'Utensils',status:'ACTIVE',stock:{available:5,reserved:0,checked_out:0,damaged_held:0,total_tracked:5},stock_sequence:1,metadata_version:1,image_id:null,created_at:'',updated_at:''}
function Location(){return <output aria-label="URL">{useLocation().pathname+useLocation().search}</output>}
function mount(page:React.ReactNode,url:string){const client=new QueryClient({defaultOptions:{queries:{retry:false}}});render(<QueryClientProvider client={client}><MemoryRouter initialEntries={[url]}>{page}<Location/></MemoryRouter></QueryClientProvider>);return client}
beforeEach(()=>{useAuthStore.setState({user:{id:'test-actor',role:'STAFF',name:'TEST Actor',email:'test@example.invalid',is_active:true}});useCartStore.getState().clear();logout.mockReset();vi.spyOn(profileApi,'metadata').mockResolvedValue({account_id:'test',version:0,image_id:null});vi.spyOn(profileApi,'own').mockResolvedValue({name:'TEST Actor',email:'test@example.invalid',role:'BORROWER',borrower_type:'STUDENT',student_id:'TEST-001',course:'TEST Culinary',contact_number:''});vi.spyOn(termsApi,'status').mockResolvedValue({state:'unpublished',current_terms:null,acceptance:null,has_previous_acceptance:false,acceptance_required:true,can_initiate_borrowing:false});vi.spyOn(inventoryApi,'categories').mockResolvedValue([]);vi.spyOn(inventoryApi,'list').mockResolvedValue({items:[equipment],page:1,per_page:25,total:1,totals:equipment.stock})})
afterEach(()=>{cleanup();vi.restoreAllMocks();useAuthStore.getState().clear();useCartStore.getState().clear()})
it('obtains status totals from complete server results with search/borrower intersections and resets paging on selection',async()=>{
 const list=vi.spyOn(borrowingApi,'list').mockImplementation(async f=>({items:[],page:f.page,per_page:f.per_page,total:f.status==='PENDING'?31:f.status?5:61} as never));mount(<BorrowingDirectory/>,'/staff/requests?search=TEST&borrower_id=test-target&page=2')
 const pending=await screen.findByRole('button',{name:'Pending 31'});expect(list).toHaveBeenCalledWith({page:1,per_page:1,search:'TEST',status:'PENDING',borrower_id:'test-target'},expect.any(AbortSignal));fireEvent.click(pending);await waitFor(()=>expect(list).toHaveBeenCalledWith({page:1,per_page:25,status:'PENDING',search:'TEST',borrower_id:'test-target'},expect.any(AbortSignal)));expect(screen.getByLabelText('URL')).toHaveTextContent('status=PENDING');fireEvent.click(screen.getByRole('button',{name:'Clear filters'}));await waitFor(()=>expect(list).toHaveBeenCalledWith(expect.objectContaining({per_page:25,status:'',search:'',borrower_id:'test-target'}),expect.any(AbortSignal)))
})
it('never manufactures failed status counts and includes them in normal borrowing invalidation',async()=>{
 const list=vi.spyOn(borrowingApi,'list').mockImplementation(async f=>{if(f.per_page===1&&f.status==='DENIED')throw Error('Unavailable');return {items:[],page:f.page,per_page:f.per_page,total:41}});const client=mount(<BorrowingDirectory/>,'/staff/requests');await screen.findByLabelText('Count unavailable');expect(screen.queryByRole('button',{name:/Denied 0/})).not.toBeInTheDocument();const before=list.mock.calls.length;await client.invalidateQueries({queryKey:['borrowings']});await waitFor(()=>expect(list.mock.calls.length).toBeGreaterThan(before));expect(list).toHaveBeenCalledWith(expect.objectContaining({per_page:1,status:'PENDING'}),expect.any(AbortSignal))
})
it('offers real Review/Open detail links and retains full reference discoverability',async()=>{
 const record={id:'test-loan',reference:'BR-11111111-1111-4111-8111-111111111111',borrower_id:'test-borrower',borrower_name:'TEST Student',borrower_type:'STUDENT',student_id:'TEST-001',status:'PENDING',created_at:'2026-10-10T00:00:00Z',expires_at:'2026-10-11T00:00:00Z',items:[{equipment_id:equipment.id,name:equipment.name,quantity:3}]} as Borrowing;vi.spyOn(borrowingApi,'list').mockResolvedValue({items:[record],total:1,page:1,per_page:25});mount(<BorrowingDirectory/>,'/staff/requests');const link=await screen.findByRole('link',{name:'Review '+record.reference});expect(link).toHaveAttribute('href','/staff/requests/test-loan');expect(screen.getByRole('region',{name:'Borrowings table'})).toHaveAttribute('tabindex','0');expect(screen.getByTitle(record.reference)).toBeInTheDocument()
})
it('uses a portalled accessible catalog filter and retains real sort filtering on close',async()=>{
 useAuthStore.setState({user:{id:'test-actor',role:'BORROWER',name:'TEST Actor',email:'test@example.invalid',is_active:true}});mount(<CatalogCart/>,'/borrower/equipment');fireEvent.click(await screen.findByRole('button',{name:'Filter & Sort'}));const category=await screen.findByRole('combobox',{name:'Category'});expect(category.closest('[data-slot=popover-content]')).toBeTruthy();expect(category.closest('.selection-browser')).toBeNull();fireEvent.change(screen.getByRole('combobox',{name:'Sort'}),{target:{value:'available'}});await waitFor(()=>expect(inventoryApi.list).toHaveBeenCalledWith(expect.objectContaining({sort:'available',page:1}),expect.any(AbortSignal)));fireEvent.keyDown(category,{key:'Escape'});await waitFor(()=>expect(screen.queryByRole('combobox',{name:'Category'})).not.toBeInTheDocument());expect(screen.getByRole('button',{name:'Filter & Sort'})).toHaveAttribute('aria-expanded','false')
})
it('keeps quantity bounds and adjacent Add controls, then allows cart editing and removal',async()=>{
 useAuthStore.setState({user:{id:'test-actor',role:'BORROWER',name:'TEST Actor',email:'test@example.invalid',is_active:true}});mount(<CatalogCart/>,'/borrower/equipment');const add=await screen.findByRole('button',{name:'Add TEST Ladle to cart'});expect(add.parentElement?.children).toHaveLength(4);expect(screen.getByRole('button',{name:'Decrease TEST Ladle quantity'})).toBeDisabled();fireEvent.click(screen.getByRole('button',{name:'Increase TEST Ladle quantity'}));fireEvent.click(add);expect(useCartStore.getState().lines[0].quantity).toBe('2');fireEvent.click(screen.getByRole('button',{name:'Edit TEST Ladle in cart'}));const quantity=await screen.findByRole('spinbutton',{name:'Cart TEST Ladle quantity'});fireEvent.change(quantity,{target:{value:'5'}});expect(screen.getByRole('button',{name:'Increase Cart TEST Ladle quantity'})).toBeDisabled();fireEvent.click(screen.getByRole('button',{name:'Remove TEST Ladle from cart'}));expect(useCartStore.getState().lines).toHaveLength(0)
})
it('has exactly one Account sign-out in the heading and preserves its existing logout hook',async()=>{
 useAuthStore.setState({user:{id:'test-actor',role:'BORROWER',name:'TEST Actor',email:'test@example.invalid',is_active:true}});mount(<WorkspacePage/>,'/borrower/account');const button=screen.getByRole('button',{name:'Sign Out'});expect(screen.getAllByRole('button',{name:'Sign Out'})).toHaveLength(1);expect(button.closest('.page-heading')).toBeTruthy();expect(button).toHaveClass('account-header-signout');fireEvent.click(button);expect(logout).toHaveBeenCalledTimes(1);await screen.findByText('TEST-001');expect(screen.getByText('Borrowing terms').closest('.account-terms-section')).toBeTruthy()
})
