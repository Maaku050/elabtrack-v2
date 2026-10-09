import type { Borrowing, BorrowingFilter, BorrowingPage, RequestInput, DecisionInput } from '../types'
const transport = async () => (await import('@/lib/api-client')).apiClient
const options = (key: string) => ({ headers: { 'Idempotency-Key': key } })
export const borrowingApi = {
 list: async (filter: BorrowingFilter, signal?: AbortSignal): Promise<BorrowingPage> => (await transport()).get('/borrowings', { params: filter, signal }),
 detail: async (id: string, signal?: AbortSignal): Promise<Borrowing> => (await transport()).get(`/borrowings/${id}`, { signal }),
 submit: async (input: RequestInput, key: string): Promise<Borrowing> => (await transport()).post('/borrowings', input, options(key)),
 direct: async (input: RequestInput, key: string): Promise<Borrowing> => (await transport()).post('/borrowings/direct-checkout', input, options(key)),
 decide: async (id: string, action: 'cancel' | 'deny' | 'approve', input: DecisionInput, key: string): Promise<Borrowing> => (await transport()).post(`/borrowings/${id}/${action}`, input, options(key)),
}
