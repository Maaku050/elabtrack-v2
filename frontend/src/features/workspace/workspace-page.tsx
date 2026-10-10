import { Mail, GraduationCap, IdCard, Phone, ShieldCheck } from 'lucide-react'
import { useOwnAccount } from '@/features/profile/hooks/use-profile'
import { Notice } from '@/components/application/management'
import { apiErrorMessage } from '@/lib/api-error'
import { ProfileEditor } from '@/features/profile/components/profile-editor'
import { Link, useLocation } from 'react-router-dom'
import { LogOut, LayoutDashboard, ArrowRight, Package, Users, UserRound } from 'lucide-react'
import { borrowerDestinations } from '@/components/application/borrower-navigation'
import { AccountTerms } from '@/features/terms/components/account-terms'
import { BorrowerShell } from '@/components/application/shells'
import { AppButton, PageHeading, SurfaceCard } from '@/components/application/visual'
import { EmptyState } from '@/components/feedback/states'
import { useAuthStore } from '@/stores/auth-store'
import { useLogout } from '@/features/auth/hooks/use-auth'

const labels: Record<string, string> = {
  '/borrower/home': 'Home', '/borrower/equipment': 'Equipment Catalog', '/borrower/borrowings': 'My Borrowings',
  '/borrower/account': 'Account', '/borrower/notifications': 'Notifications',
  '/staff/dashboard': 'Dashboard', '/staff/requests': 'Requests & Borrowings', '/staff/inventory': 'Inventory',
  '/staff/borrowers': 'Borrowers', '/staff/account': 'Account', '/staff/notifications': 'Notifications',
  '/admin/reports': 'Reports', '/admin/administration': 'Administration',
}

export function WorkspacePage() {
  const user = useAuthStore(s => s.user)
  const location = useLocation()
  const logout = useLogout()
  const own=useOwnAccount()
  if (!user) return null
  const title = labels[location.pathname] ?? 'Workspace'
  const borrower = user.role === 'BORROWER'
  const account = location.pathname.endsWith('/account')
  const signOut = <AppButton variant="outline" className={`workspace-signout ${borrower && account ? 'account-header-signout' : ''}`} aria-label="Sign Out" onClick={() => { void logout() }}><LogOut size={18} aria-hidden="true" /><span>Sign Out</span></AppButton>
  const body = <>
    <PageHeading title={title} description={account ? 'Your FSMO account.' : `Welcome, ${user.name}.`} action={account && borrower ? signOut : undefined} />
    {account ? <SurfaceCard className="personal-account"><ProfileEditor key={user.id} category={own.data?.borrower_type}/>{own.error&&<Notice tone="danger">{apiErrorMessage(own.error)}<AppButton onClick={()=>{void own.refetch()}}>Retry account information</AppButton></Notice>}<section className="personal-information"><h2>Personal Information</h2><p className="management-note">Your recorded account information. Contact FSMO for authorized corrections.</p><dl className="account-information-grid">{[[Mail,'Email',user.email],[ShieldCheck,user.role==='BORROWER'?'Borrower type':'Account role',user.role==='BORROWER'?(own.data?.borrower_type||'Borrower'):user.role],[IdCard,'Student ID',own.data?.borrower_type==='STUDENT'?own.data.student_id:''],[GraduationCap,'Course / program',own.data?.course],[Phone,'Contact number',own.data?.contact_number]].filter(([, ,value])=>value).map(([Icon,label,value])=>{const I=Icon as typeof Mail;return <div key={String(label)}><I aria-hidden="true"/><div><dt>{String(label)}</dt><dd>{String(value)}</dd></div></div>})}</dl></section></SurfaceCard>
      : !borrower && title === 'Dashboard' ? <DashboardWelcome /> : <SurfaceCard className="workspace-placeholder"><EmptyState title="This feature is not available yet" description="Your account is connected. Contact FSMO for assistance while this workspace is being prepared." /></SurfaceCard>}
    {account && borrower && <div className="account-terms-section"><AccountTerms/></div>}
  </>
  return borrower ? <BorrowerShell active={title === 'Equipment Catalog' ? 'Equipment' : title} destinations={borrowerDestinations}
>{body}</BorrowerShell>
    : body
}

function DashboardWelcome() {
  return <><SurfaceCard className="dashboard-welcome"><span className="dashboard-welcome-icon"><LayoutDashboard size={28} aria-hidden="true" /></span><p className="section-eyebrow">Your FSMO workspace</p><h2>A clear view of your operation.</h2><p>Dashboard insights will be available in a later phase. For now, manage equipment and borrower accounts from the workspaces below.</p><span className="dashboard-phase-note">This feature is not available yet</span></SurfaceCard>
    <div className="dashboard-shortcuts">{[{ title: 'Inventory', description: 'Manage equipment and physical stock.', href: '/staff/inventory', icon: Package }, { title: 'Borrowers', description: 'View borrower accounts and their status.', href: '/staff/borrowers', icon: Users }, { title: 'Your account', description: 'Review your current account details.', href: '/staff/account', icon: UserRound }].map(({ title, description, href, icon: Icon }) => <Link key={href} to={href} className="dashboard-shortcut"><Icon size={22} aria-hidden="true" /><div><h3>{title}</h3><p>{description}</p></div><ArrowRight size={18} aria-hidden="true" /></Link>)}</div></>
}
