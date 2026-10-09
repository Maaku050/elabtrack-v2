import { act, renderHook } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { afterEach, expect, it, vi } from 'vitest'
import type { ReactNode } from 'react'
import { inventoryApi } from '../api/inventory.api'
import { useCatalogImage, validateCatalogImage } from './use-catalog-image'
import type { Equipment } from '../types'
const record: Equipment = { id: 'equipment-id', name: 'TEST pot', description: '', category_id: null, category_name: '', status: 'ACTIVE', stock: { available: 15, reserved: 0, checked_out: 0, damaged_held: 0, total_tracked: 15 }, stock_sequence: 1, metadata_version: 2, image_id: null, created_at: '', updated_at: '' }
function file(type = 'image/png', size = 16, signature = [137,80,78,71,13,10,26,10]) { const f = new File(['test'], 'test.png', { type });Object.defineProperty(f,'size',{ value:size });vi.spyOn(f,'slice').mockReturnValue({arrayBuffer:async()=>new Uint8Array(signature).buffer} as Blob);return f }
afterEach(() => { vi.restoreAllMocks();vi.unstubAllGlobals() })
it('rejects oversized, unsupported, disguised and invalid-dimension images before saving', async () => {
 await expect(validateCatalogImage(file('image/png',524289))).rejects.toThrow('512 KiB')
 await expect(validateCatalogImage(file('image/svg+xml'))).rejects.toThrow('PNG or JPEG')
 await expect(validateCatalogImage(file('image/png',16,[1,2,3]))).rejects.toThrow('contents must match')
 const close = vi.fn();vi.stubGlobal('createImageBitmap',vi.fn().mockResolvedValue({width:2049,height:1,close}))
 await expect(validateCatalogImage(file())).rejects.toThrow('2048');expect(close).toHaveBeenCalledOnce()
})
it('retries only the image with the same frozen version and key after an uncertain upload failure', async () => {
 vi.stubGlobal('createImageBitmap',vi.fn().mockResolvedValue({width:32,height:32,close:vi.fn()}))
 const save=vi.spyOn(inventoryApi,'saveImage').mockRejectedValueOnce(new Error('network lost')).mockResolvedValueOnce({...record,metadata_version:3,image_id:'new-image'})
 const client=new QueryClient({defaultOptions:{mutations:{retry:false}}}),wrapper=({children}:{children:ReactNode})=><QueryClientProvider client={client}>{children}</QueryClientProvider>
 const {result}=renderHook(()=>useCatalogImage(),{wrapper});const f=file()
 await act(async()=>{await result.current.select(f)})
 await act(async()=>{await expect(result.current.save(record)).rejects.toThrow('network lost')})
 expect(result.current.file).toBe(f)
 await act(async()=>{await result.current.save({...record,metadata_version:3})})
 expect(save).toHaveBeenCalledTimes(2);expect(save.mock.calls[1]).toEqual(save.mock.calls[0]);expect(save.mock.calls[0][0].metadata_version).toBe(2);expect(result.current.file).toBeNull()
})
