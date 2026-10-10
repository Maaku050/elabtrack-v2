import type { Dashboard, ReportDefinition, ReportFilter, ReportPage } from '../types'
const transport = async () => (await import('@/lib/api-client')).apiClient
export const reportingApi = {
 dashboard: async (signal?: AbortSignal, days = 7): Promise<Dashboard> => (await transport()).get('/reporting/dashboard', { signal, params: { days } }),
 exportData: async (kind: string, filter: ReportFilter): Promise<ReportPage> => (await transport()).get(`/reporting/${encodeURIComponent(kind)}/export-data`, {params:filter}),
 definitions: async (signal?: AbortSignal): Promise<ReportDefinition[]> => (await transport()).get('/reporting/kinds', { signal }),
 report: async (kind: string, filter: ReportFilter, signal?: AbortSignal): Promise<ReportPage> => (await transport()).get(`/reporting/${encodeURIComponent(kind)}`, { params: filter, signal }),
 export: async (kind: string, filter: ReportFilter): Promise<Blob> => { const params = new URLSearchParams(); for (const [key, value] of Object.entries(filter)) if (value !== undefined && value !== '') params.set(key, String(value)); return (await transport()).download(`/reporting/${encodeURIComponent(kind)}.csv?${params}`) },
}
