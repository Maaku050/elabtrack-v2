import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useLocation, Navigate } from 'react-router-dom'
import { Eye, EyeOff } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { AppBrand, AppButton, ThemeControl, SurfaceCard } from '@/components/application/visual'
import { useAuthStore } from '@/stores/auth-store'
import { useSessionActionStore } from '@/stores/session-action-store'
import { loginSchema, type LoginSchema } from '../schemas/login.schema'
import { useLogin, useLogout } from '../hooks/use-auth'
import { loginDestination } from '../navigation'
import { normalizeApiError, apiErrorMessage } from '@/lib/api-error'

export function LoginPage() {
  const location = useLocation()
  const requested: unknown = (location.state as { from?: unknown } | null)?.from
  const login = useLogin(requested)
  const signOut = useLogout()
  const feedback = useSessionActionStore(s => s.logout)
  const logoutMessage = useSessionActionStore(s => s.message)
  const user = useAuthStore(s => s.user)
  const sessionEnded = useAuthStore(s => s.sessionEnded)
  const status = useAuthStore(s => s.status)
  const [show, setShow] = useState(false)
  const form = useForm<LoginSchema>({ resolver: zodResolver(loginSchema), defaultValues: { email: '', password: '' } })
  const busy = login.isPending || feedback === 'pending'
  if (status === 'authenticated' && user && !login.isPending) return <Navigate to={loginDestination(user.role, requested)} replace />
  const error = login.error ? normalizeApiError(login.error) : null
  const message = error?.status === 401 ? 'Unable to sign in. Check your email and password, or contact FSMO if your account is unavailable.'
    : error?.status === 0 ? 'Unable to connect. Check your connection and try again.'
    : error?.status === 429 ? 'Too many sign-in attempts. Please wait before trying again.'
    : error ? apiErrorMessage(error) : null
  async function submit(values: LoginSchema) {
    if (busy) return
    try { await login.mutateAsync(values) } catch { form.setValue('password', ''); form.setFocus('password') }
  }
  return <main className="auth-page"><header className="auth-top"><AppBrand /><ThemeControl /></header>
    <div className="login-content"><div className="login-heading"><p>FSMO / Account access</p><h1>Welcome back</h1></div>
    <SurfaceCard className="auth-card"><h2>Sign in to eLabTrack</h2>
      {sessionEnded && feedback === 'idle' && !(location.state as { signedOut?: boolean } | null)?.signedOut && <p role="status">Your session is no longer available. Sign in again to continue.</p>}
      {feedback === 'pending' && <p role="status">Signing out…</p>}
      {feedback === 'failed' && <div role="alert" className="auth-error"><p>Signed out on this device. Server sign-out could not be confirmed. Reconnect and try again before leaving this device.</p><p>{logoutMessage}</p><AppButton variant="outline" onClick={() => { void signOut() }}>Retry sign out</AppButton></div>}
      {message && <div role="alert" className="auth-error">{message}</div>}
      <form onSubmit={form.handleSubmit(submit)} noValidate aria-busy={busy}>
        <div className="auth-field"><label htmlFor="login-email">Email</label><Input id="login-email" type="email" autoComplete="username" inputMode="email" maxLength={255} disabled={busy} aria-invalid={!!form.formState.errors.email} aria-describedby={form.formState.errors.email ? 'email-error' : undefined} {...form.register('email')} />
          {form.formState.errors.email && <p id="email-error" className="field-error" role="alert">{form.formState.errors.email.message}</p>}</div>
        <div className="auth-field"><label htmlFor="login-password">Password</label><div className="password-field"><Input id="login-password" type={show ? 'text' : 'password'} autoComplete="current-password" maxLength={72} disabled={busy} aria-invalid={!!form.formState.errors.password} aria-describedby={form.formState.errors.password ? 'password-error' : undefined} {...form.register('password')} />
          <AppButton type="button" variant="ghost" aria-label={show ? 'Hide password' : 'Show password'} aria-pressed={show} onClick={() => setShow(!show)}>{show ? <EyeOff size={20} aria-hidden="true" /> : <Eye size={20} aria-hidden="true" />}</AppButton></div>
          {form.formState.errors.password && <p id="password-error" className="field-error" role="alert">{form.formState.errors.password.message}</p>}</div>
        <AppButton type="submit" className="auth-submit" disabled={busy}>{login.isPending ? 'Signing in…' : 'Sign In'}</AppButton>
      </form>
      <p className="auth-help">Accounts are provisioned by FSMO. Contact FSMO if you need an account or help signing in.</p>
    </SurfaceCard></div><footer className="auth-footer">Food Service Management Organization</footer>
  </main>
}
