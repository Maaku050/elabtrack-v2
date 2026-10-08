import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Input } from '@/components/ui/input'
import { AppBrand, AppButton, ThemeControl, SurfaceCard } from '@/components/application/visual'
import { apiErrorMessage } from '@/lib/api-error'
import { accountsApi } from '../api/accounts.api'
const schema = z.object({ password: z.string().min(8, 'Use at least 8 characters.').refine(value => new TextEncoder().encode(value).length <= 72, 'Use at most 72 UTF-8 bytes.'), confirmation: z.string() }).refine(value => value.password === value.confirmation, { path: ['confirmation'], message: 'Passwords must match.' })
export function ActivationPage() {
 const [token, setToken] = useState(() => new URLSearchParams(window.location.hash.slice(1)).get('token') ?? ''); const [complete, setComplete] = useState(false)
 const form = useForm<z.infer<typeof schema>>({ resolver: zodResolver(schema), defaultValues: { password: '', confirmation: '' } })
 useEffect(() => { window.history.replaceState(window.history.state, '', window.location.pathname + window.location.search) }, [])
 const mutation = useMutation({ mutationFn: (password: string) => accountsApi.activate(token, password), retry: false })
 async function submit(value: z.infer<typeof schema>) { if (mutation.isPending) return; try { await mutation.mutateAsync(value.password); setToken(''); setComplete(true); mutation.reset() } catch { form.setFocus('password') } finally { form.reset() } }
 return <main className="auth-page"><header className="auth-top"><AppBrand /><ThemeControl /></header><div className="login-content"><div className="login-heading"><p>FSMO / Account activation</p><h1>Create your password</h1></div><SurfaceCard className="auth-card"><h2>{complete ? 'Account activated' : 'Activate your eLabTrack account'}</h2>{complete ? <><p role="status">Your password is established. Sign in with your email and eLabTrack password. Borrowers must review current officially published FSMO terms before borrowing.</p><AppButton  nativeButton={false} render={<Link to="/login" />}> Continue to Sign In</AppButton></> : !token ? <p role="alert">Open the activation link sent to your email. If the link is missing or expired, contact an FSMO administrator.</p> : <><p>This is a separate eLabTrack password. Do not enter your institutional email password.</p>{mutation.error && <p className="auth-error" role="alert">{apiErrorMessage(mutation.error)}</p>}<form noValidate onSubmit={form.handleSubmit(submit)} aria-busy={mutation.isPending}>{(['password', 'confirmation'] as const).map(field => <div className="auth-field" key={field}><label htmlFor={`activate-${field}`}>{field === 'password' ? 'New eLabTrack password' : 'Confirm password'}</label><Input id={`activate-${field}`} type="password" autoComplete="new-password" disabled={mutation.isPending} aria-invalid={!!form.formState.errors[field]} aria-describedby={form.formState.errors[field] ? `activate-error-${field}` : undefined} {...form.register(field)} />{form.formState.errors[field] && <p id={`activate-error-${field}`} role="alert" className="field-error">{form.formState.errors[field]?.message}</p>}</div>)}<AppButton type="submit" disabled={mutation.isPending}>{mutation.isPending ? 'Activating…' : 'Activate Account'}</AppButton></form></>}<p className="auth-help">Need another link? Contact FSMO.</p></SurfaceCard></div><footer className="auth-footer">Food Service Management Organization</footer></main>
}
