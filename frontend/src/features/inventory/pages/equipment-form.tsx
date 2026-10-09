import { useRef, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { useForm, useWatch } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Input } from '@/components/ui/input'
import { AppButton } from '@/components/application/visual'
import { ManagementPage, ManagementCard, FormSection, FormField, FormActions, Notice } from '@/components/application/management'
import { NativeSelect } from '@/components/ui/native-select'
import { Textarea } from '@/components/ui/textarea'
import { ServerPagination } from '@/components/application/server-pagination'
import { useDirectoryState } from '@/hooks/use-directory-state'
import { ErrorState, LoadingState } from '@/components/feedback/states'
import { apiErrorMessage } from '@/lib/api-error'
import { inventoryApi } from '../api/inventory.api'
import { useCategories, useEquipmentDetail, useInventoryCommand } from '../hooks/use-inventory'
import { equipmentSchema, type EquipmentValues } from '../schemas/inventory.schema'
import { useCatalogImage } from '../hooks/use-catalog-image'
import { CatalogImagePicker } from '../components/catalog-image-picker'
import type { Equipment } from '../types'
function Form({ record }: { record?: Equipment }) {
 const image = useCatalogImage(), submitting = useRef(false), [saved, setSaved] = useState<Equipment | null>(null)
 const navigate = useNavigate(), view=useDirectoryState(),categoryPage=view.page('category_page'), [categoryLabel, setCategoryLabel] = useState(record?.category_name ?? ''); const categories = useCategories(categoryPage); const key = useRef({ payload: '', key: crypto.randomUUID() })
 const form = useForm<EquipmentValues>({ resolver: zodResolver(equipmentSchema), defaultValues: { name: record?.name ?? '', description: record?.description ?? '', category_id: record?.category_id ?? '', opening_quantity: 0, reason: '' } })
 const selectedCategory=useWatch({control:form.control,name:'category_id'});const categoryField=form.register('category_id')
 const command = useInventoryCommand<EquipmentValues, Equipment>(async value => { const input = { name: value.name, description: value.description, category_id: value.category_id || null, expected_version: record?.metadata_version ?? 0 }; const payload = JSON.stringify({ ...input, ...(!record ? { opening_quantity: value.opening_quantity, reason: value.reason } : {}) }); if (key.current.payload !== payload) key.current = { payload, key: crypto.randomUUID() }; return record ? inventoryApi.edit(record.id, input, key.current.key) : inventoryApi.create({ ...input, opening_quantity: value.opening_quantity, reason: value.reason }, key.current.key) })
 async function submit(value: EquipmentValues) {
  if (submitting.current || command.isPending || image.command.isPending || image.validating) return
  submitting.current = true
  try {
   let metadata = saved
   if (!metadata) { metadata = await command.mutateAsync(value); setSaved(metadata) }
   if (image.file) metadata = await image.save(metadata)
   navigate(`/staff/inventory/${metadata.id}`, { replace: true })
  } catch { if (!saved && !command.isSuccess) form.setFocus('name') } finally { submitting.current = false }
 }
 const busy=command.isPending||image.command.isPending||image.validating,locked=busy||!!saved
 return <ManagementCard>{command.error && <Notice tone="danger">{apiErrorMessage(command.error)}</Notice>}<form onSubmit={event => { void form.handleSubmit(submit)(event) }} noValidate aria-busy={busy}>
 <FormSection title="Equipment information" disabled={locked}><FormField id="equipment-name" label="Name" required error={form.formState.errors.name?.message}><Input {...form.register('name')} /></FormField><FormField id="equipment-category" label="Category (optional)"><NativeSelect {...categoryField} onChange={e=>{void categoryField.onChange(e);setCategoryLabel(e.target.selectedOptions[0]?.textContent ?? '')}}><option value="">Uncategorized</option>{selectedCategory && !categories.data?.some(c => c.id === selectedCategory) && <option value={selectedCategory}>{categoryLabel}</option>}{categories.data?.filter(c => c.is_active || c.id === record?.category_id).map(c => <option key={c.id} value={c.id}>{c.name}</option>)}</NativeSelect></FormField><FormField id="equipment-description" label="Description (optional)" error={form.formState.errors.description?.message} wide><Textarea {...form.register('description')} maxLength={2000} /></FormField></FormSection>
 {categories.isError && <Notice tone="danger">{apiErrorMessage(categories.error)}<AppButton type="button" variant="outline" onClick={()=>{void categories.refetch()}}>Retry categories</AppButton></Notice>}
 {(categories.data?.length === 100 || categoryPage > 1) && <ServerPagination label="Category options" pageKey="category_page" page={categoryPage} count={categories.data?.length??0} hasNext={categories.data?.length===100} pending={categories.isFetching||locked} onPage={p=>view.update({category_page:String(p)},false)} />}
 {!record && <FormSection title="Opening physical stock" disabled={locked}><FormField id="equipment-opening" label="Opening available quantity" required error={form.formState.errors.opening_quantity?.message}><Input type="number" min={0} max={2147483647} {...form.register('opening_quantity',{valueAsNumber:true})} /></FormField><FormField id="equipment-reason" label="Opening stock explanation" hint="Required when the opening quantity is greater than zero." error={form.formState.errors.reason?.message}><Input {...form.register('reason')} /></FormField></FormSection>}
 <CatalogImagePicker record={saved ?? record} file={image.file} error={image.error} busy={busy} onSelect={image.select} onError={image.setError} />
 {saved && image.command.error && <Notice tone="warning" title="Equipment information was saved. The image was not confirmed.">Your equipment record is preserved. Retry the image only, or continue without changing it.</Notice>}{image.command.error && <Notice tone="danger">{apiErrorMessage(image.command.error)}</Notice>}
 <p className="management-note">The catalog image is optional. Opening stock is recorded in inventory history.</p><FormActions><AppButton type="submit" disabled={busy||!!image.error}>{command.isPending ? 'Saving information…' : image.command.isPending ? 'Uploading image…' : saved ? image.file ? 'Retry Image Only' : 'Continue to Equipment' : 'Save Equipment'}</AppButton>{saved ? <AppButton type="button" variant="outline" disabled={image.command.isPending} onClick={() => navigate(`/staff/inventory/${saved.id}`, { replace: true })}>Continue to Equipment</AppButton> : <AppButton nativeButton={false} variant="outline" render={<Link to={record?`/staff/inventory/${record.id}`:'/staff/inventory'} />}>Cancel</AppButton>}</FormActions></form></ManagementCard>

}
export function EquipmentForm() { const { id } = useParams(); const query = useEquipmentDetail(id ?? ''); return <ManagementPage title={id ? 'Edit Equipment' : 'Add Equipment'} description="Maintain catalog information and images." domain="Equipment management" back={{to:id?`/staff/inventory/${id}`:'/staff/inventory',label:id?'Back to equipment':'Back to inventory'}}>{!id ? <Form /> : query.isPending ? <LoadingState /> : query.isError ? <ErrorState description={apiErrorMessage(query.error)} /> : <Form key={id} record={query.data} />}</ManagementPage> }
