import { toast } from '@/components/ui/toast'

export function useToast() {
  return {
    success: (title: string, description?: string) => toast.add({ title, description, type: 'success' }),
    error: (title: string, description?: string) => toast.add({ title, description, type: 'error', priority: 'high' }),
    info: (title: string, description?: string) => toast.add({ title, description, type: 'info' }),
  }
}
