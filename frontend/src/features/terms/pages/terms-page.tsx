import { useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { Checkbox } from '@/components/ui/checkbox'
import { AppButton, SurfaceCard } from '@/components/application/visual'
import { LoadingState } from '@/components/feedback/states'
import { apiErrorMessage, normalizeApiError } from '@/lib/api-error'
import { useAuthStore } from '@/stores/auth-store'
import { loginDestination } from '@/features/auth/navigation'
import { TermsFrame, TermsReadError } from '../components/terms-frame'
import { useAcceptTerms, useTermsStatus } from '../hooks/use-terms'
import type { TermsVersion } from '../types'

function AcceptanceForm({ version, pending, onAccept }: { version: TermsVersion; pending: boolean; onAccept: () => void }) {
 const [checked, setChecked] = useState(false)
 return <form className="terms-accept-form" aria-busy={pending} onSubmit={event => { event.preventDefault(); if (checked && !pending) onAccept() }}>
  <label className="terms-consent"><Checkbox id="terms-consent" checked={checked} onCheckedChange={value => setChecked(value)} disabled={pending} aria-describedby="terms-version" />
   <span>I have read and accept Version {version.version}.</span></label>
  <AppButton type="submit" disabled={!checked || pending}>{pending ? 'Recording acceptance…' : 'Accept and continue'}</AppButton>
 </form>
}

export function TermsPage() {
 const status = useTermsStatus(), accept = useAcceptTerms()
 const location = useLocation(), navigate = useNavigate()
 const requested: unknown = (location.state as { from?: unknown } | null)?.from
 const destination = requested === '/borrower/terms' ? '/borrower/home' : loginDestination('BORROWER', requested)
 async function submit(version: TermsVersion) {
  const { generation, user } = useAuthStore.getState()
  try {
   await accept.mutateAsync(version.id)
   const now = useAuthStore.getState()
   if (now.isAuthenticated && now.generation === generation && now.user?.id === user?.id) navigate(destination, { replace: true })
  } catch { /* Safe persistent error below; no fabricated acceptance or retry. */ }
 }
 if (status.isPending || status.isFetching) return <TermsFrame title="Borrowing terms"><SurfaceCard className="terms-card"><LoadingState label="Loading borrowing terms…" /></SurfaceCard></TermsFrame>
 if (status.isError) return <TermsReadError error={status.error} retry={() => { void status.refetch() }} />
 const data = status.data, version = data?.current_terms
 if (!version || data.state === 'unpublished') return <TermsFrame title="Borrowing terms"><SurfaceCard className="terms-card"><h2>Official terms are not available yet</h2><p>FSMO has not published the borrowing terms for this application.</p><p>New borrowing requests remain unavailable until official terms are published and you accept them. You can still view your account, existing borrowings and sign out.</p><p role="status">No terms acceptance has been recorded by opening this screen.</p></SurfaceCard></TermsFrame>
 const accepted = data.state === 'accepted' && !!data.acceptance && data.acceptance.terms_version_id === version.id
 const error = accept.error ? normalizeApiError(accept.error) : null
 return <TermsFrame title={accepted ? 'Borrowing terms' : data.state === 'updated' ? 'Updated terms' : 'Review borrowing terms'}>
  <SurfaceCard className="terms-card"><h2>{version.title}</h2><p id="terms-version" className="terms-version">Version {version.version}</p>
   <span className={`terms-status ${accepted ? 'terms-status-accepted' : ''}`}>{accepted ? 'Accepted' : 'Acceptance required'}</span>
   <p>{accepted ? 'Your acceptance applies to this exact terms version.' : data.state === 'updated' ? 'The borrowing terms have been updated. Review and accept this version before initiating another borrowing request.' : 'Read the complete borrowing terms before accepting this version.'}</p>
   <article className="terms-document" aria-label={`Complete terms, Version ${version.version}`}>{version.body}</article>
   {error && <div role="alert" className="auth-error">{error.code === 'TERMS_VERSION_CHANGED' ? 'The terms changed while you were reviewing them. Read the new version and select its acceptance checkbox again.' : error.status === 0 ? 'Unable to connect. Your acceptance could not be confirmed. Reconnect and try again.' : apiErrorMessage(error)}</div>}
   {accepted ? <div className="terms-accepted"><p role="status">Acceptance recorded on {new Date(data.acceptance!.accepted_at).toLocaleString()}.</p><Link className="app-button" to={destination}>Continue to workspace</Link></div>
    : <AcceptanceForm key={version.id} version={version} pending={accept.isPending} onAccept={() => { void submit(version) }} />}
  </SurfaceCard>
 </TermsFrame>
}
