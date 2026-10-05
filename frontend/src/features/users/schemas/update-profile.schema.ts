import { z } from 'zod'

export const updateProfileSchema = z.object({
  name: z.string().min(1, 'Name is required').max(100, 'Name is too long'),
})

export type UpdateProfileSchema = z.infer<typeof updateProfileSchema>
