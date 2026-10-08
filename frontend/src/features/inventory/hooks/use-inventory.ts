import { useMutation, useQuery } from '@tanstack/react-query'
import { useAuthStore } from '@/stores/auth-store'
import { queryClient } from '@/app/query-client'
import { ApiRequestError } from '@/lib/api-error'
import { inventoryApi } from '../api/inventory.api'
import type { EquipmentFilter } from '../types'
export function useEquipment(filter: EquipmentFilter) { const actor = useAuthStore(s => s.user?.id); return useQuery({ queryKey: ['inventory', actor, 'list', filter], queryFn: ({ signal }) => inventoryApi.list(filter, signal), enabled: !!actor, meta: { authenticated: true } }) }
export function useEquipmentDetail(id: string) { const actor = useAuthStore(s => s.user?.id); return useQuery({ queryKey: ['inventory', actor, 'detail', id], queryFn: ({ signal }) => inventoryApi.detail(id, signal), enabled: !!actor && !!id, meta: { authenticated: true } }) }
export function useCategories(page = 1) { const actor = useAuthStore(s => s.user?.id); return useQuery({ queryKey: ['inventory', actor, 'categories', page], queryFn: ({ signal }) => inventoryApi.categories(page, signal), enabled: !!actor, meta: { authenticated: true } }) }
export function useMovements(id: string, page: number) { const actor = useAuthStore(s => s.user); return useQuery({ queryKey: ['inventory', actor?.id, 'movements', id, page], queryFn: ({ signal }) => inventoryApi.movements(id, page, signal), enabled: !!actor && actor.role !== 'BORROWER' && !!id, meta: { authenticated: true } }) }
export function useInventoryCommand<I, O>(fn: (input: I) => Promise<O>) { return useMutation({ mutationFn: async (input: I) => { const generation = useAuthStore.getState().generation; const data = await fn(input); if (generation !== useAuthStore.getState().generation) throw new ApiRequestError('Session changed.', 401, 'SESSION_CHANGED'); return data }, onSuccess: async () => { await queryClient.invalidateQueries({ queryKey: ['inventory'] }) }, retry: false }) }
