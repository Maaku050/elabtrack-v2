import { z } from 'zod'

const healthSchema = z.object({
  status: z.enum(['ok', 'degraded']),
  services: z.record(z.string(), z.enum(['healthy', 'unhealthy', 'unknown'])),
})

// The retained health endpoint predates the normal response envelope.
// This unauthenticated read never attaches the provisional starter session.
export async function getHealth(signal?: AbortSignal) {
  const base = (import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1').replace(/\/$/, '')
  const response = await fetch(`${base}/health`, { signal, credentials: 'omit' })
  if (!response.ok) throw new Error('Service unavailable. Check that the API and database are running, then try again.')
  return healthSchema.parse(await response.json())
}
