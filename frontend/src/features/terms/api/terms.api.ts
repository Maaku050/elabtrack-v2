import type { TermsAcceptance, TermsStatus, TermsVersion } from '../types'
// The existing centralized client owns authentication, cancellation and retries.
const transport = async () => (await import('@/lib/api-client')).apiClient
export const termsApi = {
 status: async (signal?: AbortSignal): Promise<TermsStatus> => (await transport()).get('/terms/status', { signal }),
 current: async (signal?: AbortSignal): Promise<TermsVersion> => (await transport()).get('/terms/current', { signal }),
 accept: async (versionId: string): Promise<TermsAcceptance> => (await transport()).post(`/terms/${encodeURIComponent(versionId)}/accept`, {}),
}
