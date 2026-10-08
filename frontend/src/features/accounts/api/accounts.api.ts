import type { AccountFilter, AccountInput, AccountPage, AccountRecord, AuditEvent, ProvisioningPolicy, RosterBatch } from '../types'
const transport = async () => (await import('@/lib/api-client')).apiClient
const path = (staff = false) => staff ? '/staff-accounts' : '/borrowers'
const options = (key: string) => ({ headers: { 'Idempotency-Key': key } })
export const accountsApi = {
 list: async (filter: AccountFilter, staff: boolean, signal?: AbortSignal): Promise<AccountPage> => (await transport()).get(path(staff), { params: filter, signal }),
 detail: async (id: string, staff: boolean, signal?: AbortSignal): Promise<AccountRecord> => (await transport()).get(`${path(staff)}/${encodeURIComponent(id)}`, { signal }),
 policy: async (signal?: AbortSignal): Promise<ProvisioningPolicy> => (await transport()).get('/borrowers/policy', { signal }),
 audit: async (id: string, staff: boolean, page: number, signal?: AbortSignal): Promise<AuditEvent[]> => (await transport()).get(`${path(staff)}/${encodeURIComponent(id)}/audit`, { params: { page }, signal }),
 create: async (input: AccountInput, staff: boolean, key: string): Promise<AccountRecord> => (await transport()).post(path(staff), input, options(key)),
 profile: async (id: string, input: { name: string; course: string; contact_number: string; expected_updated_at: string }, key: string): Promise<AccountRecord> => (await transport()).patch(`/borrowers/${id}`, input, options(key)),
 status: async (account: AccountRecord, staff: boolean, key: string): Promise<AccountRecord> => (await transport()).patch(`${path(staff)}/${account.id}/status`, { active: !account.is_active, confirm: true, expected_updated_at: account.updated_at }, options(key)),
 resend: async (id: string, staff: boolean, key: string): Promise<AccountRecord> => (await transport()).post(`${path(staff)}/${id}/activation`, {}, options(key)),
 template: async (): Promise<Blob> => (await transport()).download('/borrowers/student-template'),
 prepare: async (file: File, operation: string): Promise<RosterBatch> => { const form = new FormData(); form.append('roster', file); form.append('operation', operation); return (await transport()).post('/borrowers/rosters', form, { headers: { 'Content-Type': undefined } }) },
 confirm: async (id: string, rows: number[]): Promise<RosterBatch> => (await transport()).post(`/borrowers/rosters/${id}/confirm`, { rows, confirm: true }),
 activate: async (token: string, password: string): Promise<void> => (await transport()).post('/auth/activate', { token, password }, { attachAuth: false, retryAuth: false }),
}
