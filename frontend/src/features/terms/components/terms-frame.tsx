import type { ReactNode } from 'react'
import { Link } from 'react-router-dom'
import { Bell, LogOut } from 'lucide-react'
import { BorrowerShell } from '@/components/application/shells'
import { borrowerDestinations } from '@/components/application/borrower-navigation'
import { AppButton, PageHeading, ThemeControl, SurfaceCard } from '@/components/application/visual'
import { useLogout } from '@/features/auth/hooks/use-auth'
import { apiErrorMessage } from '@/lib/api-error'

export function TermsFrame({ title, children }: { title: string; children: ReactNode }) {
 const logout = useLogout()
 return <BorrowerShell active="Account" destinations={borrowerDestinations} action={<Link className="app-button" to="/borrower/notifications" aria-label="Notifications"><Bell size={20} aria-hidden="true" /></Link>}>
  <div className="terms-layout"><p className="terms-eyebrow">FSMO / Account access</p>
  <PageHeading title={title} description="Review your account’s borrowing terms." action={<div className="workspace-actions"><ThemeControl /><AppButton variant="outline" onClick={() => { void logout() }}><LogOut size={18} aria-hidden="true" />Sign Out</AppButton></div>} />
  {children}<div className="terms-readonly-links"><Link className="app-button" to="/borrower/account">View account</Link><Link className="app-button" to="/borrower/borrowings">View existing borrowings</Link></div></div>
 </BorrowerShell>
}
export function TermsReadError({ error, retry }: { error: unknown; retry: () => void }) {
 return <TermsFrame title="Borrowing terms"><SurfaceCard className="terms-card"><h2>Unable to load borrowing terms</h2><p role="alert">{apiErrorMessage(error)}</p><p>Acceptance could not be confirmed. You can still view your account and existing borrowings.</p><AppButton variant="outline" onClick={retry}>Try again</AppButton></SurfaceCard></TermsFrame>
}
