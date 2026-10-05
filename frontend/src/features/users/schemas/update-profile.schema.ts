import { z } from 'zod'

export const updateProfileSchema = z.object({
  name: z.string().min(1, 'Name is required').max(100, 'Name is too long').optional(),
  email: z.string().email('Enter a valid email address').optional(),
})

export type UpdateProfileSchema = z.infer<typeof updateProfileSchema>
