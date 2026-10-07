import { z } from 'zod'
import { ApiRequestError } from '@/lib/api-error'

const healthSchema = z.object({ status: z.literal('ok'), service: z.literal('elabtrack-v2') })

// Central transport/envelope/error normalization; diagnostic remains cookie-free,
// bearer-free and outside automatic session refresh. Liveness only, not DB readiness.
export async function getHealth(signal?: AbortSignal) {
  const { apiClient } = await import('@/lib/api-client')
  const data = await apiClient.get<unknown>('/health', { signal, withCredentials: false, attachAuth: false, retryAuth: false })
  const result = healthSchema.safeParse(data)
  if (!result.success) throw new ApiRequestError('Invalid service response.', 502, 'INVALID_RESPONSE')
  return result.data
}
