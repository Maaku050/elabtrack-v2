import { useState, type ReactNode } from 'react'
import { Link } from 'react-router-dom'
import { Home, Package, List, UserRound, Bell } from 'lucide-react'
import { AppButton, ThemeControl } from '@/components/application/visual'
import { BorrowerShell } from '@/components/application/shells'
import { PreviewDialog } from '@/components/application/feedback'

export function PreviewNotice({ staff = false }: { staff?: boolean }) {
  return <aside className="preview-notice" aria-label="Preview information"><p><strong>Development visual preview</strong> · Synthetic data. Nothing is reserved or recorded. Fixture snapshot: 08 Oct 2026, 10:00 AM (Asia/Manila).</p>
    <div>{!staff && <ThemeControl />}<AppButton variant="ghost" nativeButton={false} render={<Link to={staff ? '/__preview/borrower/home' : '/__preview/staff/dashboard'} />}>{staff ? 'Borrower preview' : 'Staff preview'}</AppButton>
      {staff && <AppButton variant="ghost" nativeButton={false} render={<Link to="/__preview/admin/dashboard" />}>Admin shell preview</AppButton>}
    </div>
  </aside>
}

export function BorrowerPreviewFrame({ active, children, headerAction }: { active: 'Home' | 'Equipment'; children: ReactNode; headerAction?: ReactNode }) {
  const [disclosure, setDisclosure] = useState<string | null>(null)
  const destinations = [
    { label: 'Home', icon: Home, href: '/__preview/borrower/home' },
    { label: 'Equipment', icon: Package, href: '/__preview/borrower/equipment' },
    { label: 'My Borrowings', icon: List, onSelect: () => setDisclosure('My Borrowings') },
    { label: 'Account', icon: UserRound, onSelect: () => setDisclosure('Account') },
  ]
  return <><BorrowerShell active={active} destinations={destinations} preview action={headerAction ?? <AppButton variant="ghost" className="notification-button" aria-label="Notifications preview" onClick={() => setDisclosure('Notifications')}><Bell size={23} aria-hidden="true" /><span className="notification-dot" aria-hidden="true" /></AppButton>}>
    {children}<PreviewNotice />
  </BorrowerShell><PreviewDialog title={`${disclosure ?? 'Feature'} preview`} open={disclosure !== null} onOpenChange={open => { if (!open) setDisclosure(null) }}>This screen belongs to a later feature phase. This preview does not load accounts, borrowings, or notifications.</PreviewDialog></>
}
