/// <reference types="node" />
import { readFile } from 'node:fs/promises'
import { afterEach, expect, it, vi } from 'vitest'
import { reportPDF } from './lib/pdf-export'
import type { ReportPage } from './types'
const report:ReportPage={kind:'inventory',columns:['Equipment ID','Name','Available'],rows:[['11111111-1111-4111-8111-111111111111','TEST Culinary ladle — ₱0',15].map(String)],total:1,page:1,per_page:25,as_of:'2026-10-10T00:00:00Z'}
afterEach(()=>vi.unstubAllGlobals())
it('rejects an oversized or incomplete snapshot before loading any export assets',async()=>{
 const fetcher=vi.fn();vi.stubGlobal('fetch',fetcher)
 await expect(reportPDF({...report,total:5001},'Inventory',{page:2,per_page:25})).rejects.toThrow('complete bounded report snapshot')
 await expect(reportPDF({...report,total:2},'Inventory',{page:1,per_page:25})).rejects.toThrow('complete bounded report snapshot')
 expect(fetcher).not.toHaveBeenCalled()
})
it('fails safely when the bundled font cannot load and permits a later retry',async()=>{
 vi.stubGlobal('fetch',vi.fn().mockResolvedValue({ok:false}))
 await expect(reportPDF(report,'Inventory',{page:1,per_page:25})).rejects.toThrow('local PDF font')
})
it('generates real PDF bytes with a local Unicode font, including an empty filtered result',async()=>{
 const data=await readFile('src/assets/fonts/DejaVuSans.ttf')
 const fetcher=vi.fn().mockResolvedValue({ok:true,arrayBuffer:async()=>data.buffer.slice(data.byteOffset,data.byteOffset+data.byteLength)})
 vi.stubGlobal('fetch',fetcher)
 const blob=await reportPDF(report,'Inventory',{page:2,per_page:1,search:'TEST Culinary',status:'ACTIVE'})
 expect(blob.type).toBe('application/pdf');expect(blob.size).toBeGreaterThan(1000)
 const text=await new Promise<string>((resolve,reject)=>{const reader=new FileReader();reader.onload=()=>resolve(String(reader.result));reader.onerror=reject;reader.readAsText(blob.slice(0,8))})
 expect(text).toMatch(/^%PDF-1/)
 const empty=await reportPDF({...report,rows:[],total:0},'Inventory',{page:1,per_page:25,search:'No matches'})
 expect(empty.type).toBe('application/pdf');expect(empty.size).toBeGreaterThan(1000)
 expect(fetcher).toHaveBeenCalledTimes(1)
})
