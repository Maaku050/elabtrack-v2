import { useDeferredValue, useState } from 'react'
import { Link } from 'react-router-dom'
import { BorrowerShell } from '@/components/application/shells'
import { borrowerDestinations } from '@/components/application/borrower-navigation'
import { AppButton, SurfaceCard, ThemeControl } from '@/components/application/visual'
import { ManagementPage, ManagementCard } from '@/components/application/management'
import { DirectorySearch, DirectorySelect } from '@/components/application/directory'
import { ServerPagination } from '@/components/application/server-pagination'
import { AvailabilityBadge } from '@/components/application/equipment'
import { EmptyState, ErrorState, LoadingState } from '@/components/feedback/states'
import { useDirectoryState } from '@/hooks/use-directory-state'
import { apiErrorMessage } from '@/lib/api-error'
import { useCategories, useEquipment } from '../hooks/use-inventory'
import { EquipmentImage } from './equipment-image'
import { InventoryDirectory } from './inventory-directory'
export function BorrowerCatalog() {
 const view=useDirectoryState(),page=view.page(),categoryPage=view.page('category_page'),search=view.value('search'),deferred=useDeferredValue(search),category=view.value('category_id'),available=view.value('available_only')==='true',sort=view.value('sort')==='available'?'available':'name'
 const [categoryLabel,setCategoryLabel]=useState('Selected category')
 const categories=useCategories(categoryPage),query=useEquipment({page,per_page:25,search:deferred,category_id:category||undefined,status:undefined,available_only:available,sort}),data=query.data
 return <BorrowerShell active="Equipment" destinations={borrowerDestinations}><ManagementPage title="Equipment Catalog" domain="Equipment catalog" description="Browse FSMO equipment and current availability." actions={<><AppButton nativeButton={false} render={<Link to="/borrower/borrowings/new" />}>New Request</AppButton><ThemeControl /></>}><ManagementCard title="Browse equipment"><section className="directory-filters" aria-label="Equipment catalog filters"><DirectorySearch label="Search equipment" placeholder="Equipment name" value={search} onChange={e=>view.update({search:e.target.value},true,true)} /><DirectorySelect label="Category" value={category} onChange={e=>{setCategoryLabel(e.target.selectedOptions[0]?.textContent??'Selected category');view.update({category_id:e.target.value})}}><option value="">All categories</option>{category&&!categories.data?.some(c=>c.id===category)&&<option value={category}>{categoryLabel}</option>}{categories.data?.map(c=><option key={c.id} value={c.id}>{c.name}{!c.is_active?' (inactive)':''}</option>)}</DirectorySelect><DirectorySelect label="Sort" value={sort} onChange={e=>view.update({sort:e.target.value})}><option value="name">Name</option><option value="available">Available quantity</option></DirectorySelect><label className="management-checkbox"><input type="checkbox" checked={available} onChange={e=>view.update({available_only:e.target.checked?'true':''})} />Available for borrowing</label></section>
 {categories.isError&&<ErrorState description={apiErrorMessage(categories.error)} onRetry={()=>{void categories.refetch()}} />}{(categories.data?.length===100||categoryPage>1)&&<ServerPagination label="Catalog categories" pageKey="category_page" page={categoryPage} count={categories.data?.length??0} hasNext={categories.data?.length===100} pending={categories.isFetching} onPage={p=>view.update({category_page:String(p)},false)} />}</ManagementCard>
 {query.isPending?<LoadingState label="Loading equipment…" />:query.isError?<ErrorState description={apiErrorMessage(query.error)} onRetry={()=>{void query.refetch()}} />:!data?.items.length?<ManagementCard><EmptyState title={search||category||available?'No equipment matches':'No equipment yet'} description="Contact FSMO for assistance." /></ManagementCard>:<div className="inventory-catalog">{data.items.map(v=><SurfaceCard className="equipment-card inventory-catalog-card" key={v.id}><EquipmentImage record={v} /><div className="equipment-copy"><h2>{v.name}</h2><p>{v.category_name||'Uncategorized'}</p><p className="inventory-description">{v.description}</p></div><div className="equipment-controls"><AvailabilityBadge available={v.stock.available} state={v.stock.available>0?'available':'unavailable'} /><AppButton variant="outline" nativeButton={false} render={<Link to={`/borrower/equipment/${v.id}`} />}>View Details</AppButton></div></SurfaceCard>)}</div>}
 {data&&!query.isError&&<ServerPagination label="Equipment types" page={data.page} total={data.total} perPage={data.per_page} totalPages={Math.ceil(data.total/data.per_page)} count={data.items.length} pending={query.isFetching} onPage={p=>view.update({page:String(p)},false)} />}</ManagementPage></BorrowerShell>
}
export function EquipmentList({borrower=false}:{borrower?:boolean}){return borrower?<BorrowerCatalog />:<InventoryDirectory />}
