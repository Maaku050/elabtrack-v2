import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Input } from '@/components/ui/input'
import { AppBrand, AppButton, ThemeControl } from '@/components/application/visual'
import { ManagementCard, FormSection, FormField, FormActions, Notice } from '@/components/application/management'
import { apiErrorMessage } from '@/lib/api-error'
import { accountsApi } from '../api/accounts.api'
const schema = z.object({ password: z.string().min(8, 'Use at least 8 characters.').refine(value => new TextEncoder().encode(value).length <= 72, 'Use at most 72 UTF-8 bytes.'), confirmation: z.string() }).refine(value => value.password === value.confirmation, { path: ['confirmation'], message: 'Passwords must match.' })
export function ActivationPage() {
 const [token, setToken] = useState(() => new URLSearchParams(window.location.hash.slice(1)).get('token') ?? ''); const [complete, setComplete] = useState(false)
 const form = useForm<z.infer<typeof schema>>({ resolver: zodResolver(schema), defaultValues: { password: '', confirmation: '' } })
 useEffect(() => { window.history.replaceState(window.history.state, '', window.location.pathname + window.location.search) }, [])
 const mutation = useMutation({ mutationFn: (password: string) => accountsApi.activate(token, password), retry: false })
 async function submit(value: z.infer<typeof schema>) { if (mutation.isPending) return; try { await mutation.mutateAsync(value.password); setToken(''); setComplete(true); mutation.reset() } catch { form.setFocus('password') } finally { form.reset() } }
 return <main className="auth-page"><header className="auth-top"><AppBrand /><ThemeControl /></header><div className="login-content activation-form"><div className="login-heading"><p>FSMO / Account activation</p><h1>Create your password</h1></div><ManagementCard title={complete ? 'Account activated' : 'Activate your eLabTrack account'}>{complete ? <><Notice tone="success">Your password is established. Sign in with your email and eLabTrack password. Borrowers must review current officially published FSMO terms before borrowing.</Notice><AppButton nativeButton={false} render={<Link to="/login" />}>Continue to Sign In</AppButton></> : !token ? <Notice tone="warning">Open the activation link sent to your email. If the link is missing or expired, contact an FSMO administrator.</Notice> : <><p>This is a separate eLabTrack password. Do not enter your institutional email password.</p>{mutation.error && <Notice tone="danger">{apiErrorMessage(mutation.error)}</Notice>}<form noValidate onSubmit={form.handleSubmit(submit)} aria-busy={mutation.isPending}><FormSection title="Your eLabTrack password" disabled={mutation.isPending}>{(['password', 'confirmation'] as const).map(field => <FormField key={field} id={`activate-${field}`} required label={field === 'password' ? 'New eLabTrack password' : 'Confirm password'} error={form.formState.errors[field]?.message}><Input type="password" autoComplete="new-password" {...form.register(field)} /></FormField>)}</FormSection><FormActions><AppButton type="submit" disabled={mutation.isPending}>{mutation.isPending ? 'Activating…' : 'Activate Account'}</AppButton></FormActions></form></>}<p className="auth-help">Need another link? Contact FSMO.</p></ManagementCard></div><footer className="auth-footer">Food Service Management Organization</footer></main>
}
