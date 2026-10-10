import type { Borrowing, BorrowingFilter, BorrowingPage, RequestInput, DecisionInput, ReturnInput, ReplacementInput, FineInput } from '../types'
const transport = async () => (await import('@/lib/api-client')).apiClient
const options = (key: string) => ({ headers: { 'Idempotency-Key': key } })
export const borrowingApi = {
 eligibility: async (id:string):Promise<{eligible:boolean}> => (await transport()).get(`/borrowings/eligibility/${id}`),
 history: async <T=Record<string,unknown>>(id:string,kind:string,page:number,perPage=25,signal?:AbortSignal):Promise<import('../types').HistoryPage<T>> => (await transport()).get(`/borrowings/${id}/history`,{params:{kind,page,per_page:perPage},signal}),
 list: async (filter: BorrowingFilter, signal?: AbortSignal): Promise<BorrowingPage> => (await transport()).get('/borrowings', { params: filter, signal }),
 detail: async (id: string, signal?: AbortSignal): Promise<Borrowing> => (await transport()).get(`/borrowings/${id}`, { signal }),
 submit: async (input: RequestInput, key: string): Promise<Borrowing> => (await transport()).post('/borrowings', input, options(key)),
 direct: async (input: RequestInput, key: string): Promise<Borrowing> => (await transport()).post('/borrowings/direct-checkout', input, options(key)),
 decide: async (id: string, action: 'cancel' | 'deny' | 'approve', input: DecisionInput, key: string): Promise<Borrowing> => (await transport()).post(`/borrowings/${id}/${action}`, input, options(key)),
 recordReturn: async (id: string, input: ReturnInput, key: string): Promise<Borrowing> => (await transport()).post(`/borrowings/${id}/returns`, input, options(key)),
 replace: async (id: string, input: ReplacementInput, key: string): Promise<Borrowing> => (await transport()).post(`/borrowings/${id}/replacements`, input, options(key)),
 clearFine: async (id: string, input: FineInput, key: string): Promise<Borrowing> => (await transport()).post(`/borrowings/${id}/fine-clearances`, input, options(key)),
}
