import { useRef, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Plus } from 'lucide-react'
import { AppButton, ToneBadge } from '@/components/application/visual'
import { ManagementPage, FormSection, FormField, FormActions, Notice } from '@/components/application/management'
import { DirectoryPanel, DirectoryTableRegion } from '@/components/application/directory'
import { ServerPagination } from '@/components/application/server-pagination'
import { Input } from '@/components/ui/input'
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/components/ui/dialog'
import { useDirectoryState } from '@/hooks/use-directory-state'
import { apiErrorMessage } from '@/lib/api-error'
import { inventoryApi } from '../api/inventory.api'
import { useCategories, useInventoryCommand } from '../hooks/use-inventory'
import type { Category } from '../types'
const schema = z.object({ name: z.string().trim().min(1, 'Enter a category name.').max(100), is_active: z.boolean() })
type Values = z.infer<typeof schema>
function CategoryForm({ record, done, cancel, pending }: { record?: Category; done: () => void; cancel: () => void; pending: (value:boolean)=>void }) {
 const key = useRef({ payload: '', key: crypto.randomUUID() }); const form = useForm<Values>({ resolver: zodResolver(schema), defaultValues: { name: record?.name ?? '', is_active: record?.is_active ?? true } })
 const command = useInventoryCommand<Values, Category>(v => { const payload = JSON.stringify(v); if (key.current.payload !== payload) key.current = { payload, key: crypto.randomUUID() }; return inventoryApi.category(record?.id, { ...v, expected_version: record?.version ?? 0 }, key.current.key) })
 return <form className="management-dialog-form" noValidate onSubmit={form.handleSubmit(async v => { if (command.isPending) return; pending(true);try { await command.mutateAsync(v); form.reset(); done() } catch { form.setFocus('name') } finally {pending(false)} })} aria-busy={command.isPending}>{command.error && <Notice tone="danger">{apiErrorMessage(command.error)}</Notice>}<FormSection title="Category information" disabled={command.isPending}><FormField id="category-name" label="Category name" required error={form.formState.errors.name?.message} wide><Input {...form.register('name')} /></FormField><label className="management-checkbox form-field-wide"><input type="checkbox" {...form.register('is_active')} />Available for new equipment</label></FormSection><p className="management-note">Existing equipment keeps its category when the category is made inactive.</p><FormActions><AppButton type="submit" disabled={command.isPending}>{command.isPending?'Saving…':'Save Category'}</AppButton><AppButton type="button" variant="outline" disabled={command.isPending} onClick={cancel}>Cancel</AppButton></FormActions></form>
}
export function CategoryManagement() {
 const view=useDirectoryState(),page=view.page();const [editing,setEditing]=useState<Category|undefined>(),[open,setOpen]=useState(false),[revision,setRevision]=useState(0),[notice,setNotice]=useState(''),[saving,setSaving]=useState(false);const query=useCategories(page)
 function edit(record?:Category){setEditing(record);setRevision(r=>r+1);setOpen(true)}
 return <ManagementPage title="Equipment Categories" description="Manage FSMO catalog categories." domain="Equipment management" back={{to:'/staff/inventory',label:'Back to inventory'}} actions={<AppButton onClick={()=>edit()}><Plus aria-hidden="true" />Add Category</AppButton>}>{notice && <Notice tone="success">{notice}</Notice>}<DirectoryPanel title="Categories" summary="Names and availability for new equipment" toolbar={<span className="management-note">{query.data?.length??'—'} categories on this page</span>} pending={query.isPending} error={query.isError?apiErrorMessage(query.error):undefined} empty={!query.data?.length} onRetry={()=>{void query.refetch()}} footer={<ServerPagination label="Categories" page={page} count={query.data?.length??0} hasNext={query.data?.length===100} pending={query.isFetching} onPage={p=>view.update({page:String(p)},false)} />}>
 <DirectoryTableRegion label="Equipment categories table"><Table className="management-compact-table"><TableHeader><TableRow><TableHead scope="col">Category</TableHead><TableHead scope="col">Status</TableHead><TableHead scope="col"><span className="sr-only">Actions</span></TableHead></TableRow></TableHeader><TableBody>{query.data?.map(c=><TableRow key={c.id}><TableCell><strong>{c.name}</strong></TableCell><TableCell><ToneBadge tone={c.is_active?'success':'neutral'}>{c.is_active?'Active':'Inactive'}</ToneBadge></TableCell><TableCell><AppButton variant="outline" aria-label={`Edit ${c.name}`} onClick={()=>edit(c)}>Edit</AppButton></TableCell></TableRow>)}</TableBody></Table></DirectoryTableRegion></DirectoryPanel>
 <Dialog open={open} onOpenChange={next=>{if(!saving)setOpen(next)}}><DialogContent className="management-dialog" showCloseButton={!saving}><DialogHeader><DialogTitle>{editing?'Edit Category':'Add Category'}</DialogTitle><DialogDescription>Maintain a catalog category without changing equipment relationships.</DialogDescription></DialogHeader><CategoryForm key={`${editing?.id??'new'}-${revision}`} record={editing} pending={setSaving} done={()=>{setNotice(editing?'Category updated.':'Category created.');setOpen(false)}} cancel={()=>setOpen(false)} /></DialogContent></Dialog></ManagementPage>
}
