import type { Adjustment, Category, Equipment, EquipmentFilter, EquipmentPage, Metadata, Movement } from '../types'
const transport = async () => (await import('@/lib/api-client')).apiClient
const options = (key: string) => ({ headers: { 'Idempotency-Key': key } })
export const inventoryApi = {
 list: async (filter: EquipmentFilter, signal?: AbortSignal): Promise<EquipmentPage> => (await transport()).get('/equipment', { params: filter, signal }),
 detail: async (id: string, signal?: AbortSignal): Promise<Equipment> => (await transport()).get(`/equipment/${id}`, { signal }),
 categories: async (page: number, signal?: AbortSignal): Promise<Category[]> => (await transport()).get('/equipment-categories', { params: { page }, signal }),
 category: async (id: string | undefined, input: { name: string; is_active: boolean; expected_version: number }, key: string): Promise<Category> => id ? (await transport()).patch(`/equipment-categories/${id}`, input, options(key)) : (await transport()).post('/equipment-categories', input, options(key)),
 create: async (input: Metadata & { opening_quantity: number; reason: string }, key: string): Promise<Equipment> => (await transport()).post('/equipment', input, options(key)),
 edit: async (id: string, input: Metadata, key: string): Promise<Equipment> => (await transport()).patch(`/equipment/${id}`, input, options(key)),
 status: async (record: Equipment, status: string, key: string): Promise<Equipment> => (await transport()).patch(`/equipment/${record.id}/status`, { status, expected_version: record.metadata_version, confirm: true }, options(key)),
 adjust: async (id: string, input: Adjustment, key: string): Promise<Equipment> => (await transport()).post(`/equipment/${id}/adjustments`, input, options(key)),
 movements: async (id: string, page: number, signal?: AbortSignal): Promise<Movement[]> => (await transport()).get(`/equipment/${id}/movements`, { params: { page }, signal }),
 image: async (record: Equipment, signal?: AbortSignal): Promise<Blob> => (await transport()).download(`/equipment/${record.id}/images/${record.image_id}`, signal),
 saveImage: async (record: Equipment, file: File, key: string): Promise<Equipment> => { const form = new FormData(); form.append('image', file); form.append('expected_version', String(record.metadata_version)); return (await transport()).post(`/equipment/${record.id}/image`, form, { headers: { 'Idempotency-Key': key, 'Content-Type': undefined } }) },
}
