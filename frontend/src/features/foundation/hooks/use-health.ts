import { useQuery } from '@tanstack/react-query'
import { getHealth } from '@/features/foundation/api/health.api'

export function useHealth() {
  return useQuery({
    queryKey: ['foundation', 'health'],
    queryFn: ({ signal }) => getHealth(signal),
    enabled: false,
    retry: false,
  })
}
