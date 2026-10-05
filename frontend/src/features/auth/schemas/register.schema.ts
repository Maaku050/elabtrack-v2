import { z } from 'zod'

export const registerSchema = z
  .object({
    email: z.string().email('Enter a valid email address'),
    name: z.string().min(1, 'Name is required').max(100, 'Name is too long'),
    password: z.string().min(8, 'Password must be at least 8 characters').max(72, 'Password is too long'),
    confirmPassword: z.string(),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword'],
  })

export type RegisterSchema = z.infer<typeof registerSchema>
