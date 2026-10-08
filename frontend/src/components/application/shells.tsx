import { useState, type ReactNode } from 'react'
import { Link } from 'react-router-dom'
import { Bell, LayoutDashboard, ListChecks, Package, Users, ChartNoAxesCombined, Settings, Menu, Search, ChevronDown, type LucideIcon } from 'lucide-react'
import { SidebarProvider, Sidebar } from '@/components/ui/sidebar'
import { Sheet, SheetContent, SheetTitle, SheetDescription, SheetHeader } from '@/components/ui/sheet'
import { useUIStore } from '@/stores/ui-store'
import { AppBrand, AppButton, ThemeControl } from './visual'

export interface NavigationDestination { label: string; icon: LucideIcon; href?: string; onSelect?: () => void }

function Destination({ item, active, onNavigate }: { item: NavigationDestination; active: boolean; onNavigate?: () => void }) {
  const Icon = item.icon
  const content = <><Icon size={21} aria-hidden="true" /><span>{item.label}</span></>
  return item.href
    ? <Link aria-label={item.label} aria-current={active ? 'page' : undefined} to={item.href} onClick={onNavigate}>{content}</Link>
    : <button type="button" aria-label={item.label} onClick={() => { onNavigate?.(); item.onSelect?.() }}>{content}</button>
}

export function BorrowerTopBar({ action, preview = false }: { action?: ReactNode; preview?: boolean }) {
  return <header className="borrower-top-bar"><div className="borrower-top-bar-inner"><AppBrand /><div className="top-bar-actions">
    {preview && <span className="preview-tag">Preview</span>}{action}
  </div></div></header>
}

export function BorrowerBottomNavigation({ destinations, active }: { destinations: readonly NavigationDestination[]; active: string }) {
  return <nav className="borrower-bottom-nav" aria-label="Borrower navigation"><div>{destinations.map(item => <Destination key={item.label} item={item} active={active === item.label} />)}</div></nav>
}

export function BorrowerShell({ children, destinations, active, action, preview = false }: { children: ReactNode; destinations: readonly NavigationDestination[]; active: string; action?: ReactNode; preview?: boolean }) {
  return <div className="borrower-shell"><a className="skip-link" href="#borrower-content">Skip to content</a>
    <BorrowerTopBar action={action} preview={preview} /><main id="borrower-content" tabIndex={-1} className="borrower-content">{children}</main>
    <BorrowerBottomNavigation destinations={destinations} active={active} />
  </div>
}

export function StaffSidebar({ audience, active, dashboardHref, requestsHref, onUnavailable, onNavigate }: {
  audience: 'staff' | 'admin'; active: string; dashboardHref: string; requestsHref: string; onUnavailable: (label: string) => void; onNavigate?: () => void
}) {
  const operational: NavigationDestination[] = [
    { label: 'Dashboard', icon: LayoutDashboard, href: dashboardHref },
    { label: 'Requests & Borrowings', icon: ListChecks, href: requestsHref },
    { label: 'Inventory', icon: Package, onSelect: () => onUnavailable('Inventory') },
    { label: 'Borrowers', icon: Users, onSelect: () => onUnavailable('Borrowers') },
  ]
  const administration: NavigationDestination[] = [
    { label: 'Reports', icon: ChartNoAxesCombined, onSelect: () => onUnavailable('Reports') },
    { label: 'Administration', icon: Settings, onSelect: () => onUnavailable('Administration') },
  ]
  return <div className="staff-sidebar-content"><div className="staff-brand"><AppBrand compact /></div>
    <nav aria-label={`${audience === 'admin' ? 'Admin' : 'Staff'} navigation`}>
      {operational.map(item => <Destination key={item.label} item={item} active={active === item.label} onNavigate={onNavigate} />)}
      {audience === 'admin' && <div className="admin-navigation">{administration.map(item => <Destination key={item.label} item={item} active={active === item.label} onNavigate={onNavigate} />)}</div>}
    </nav>
    <div className="staff-sidebar-footer"><ThemeControl showLabel /><div className="staff-identity"><span className="account-avatar">{audience === 'admin' ? 'AD' : 'ST'}</span><div><strong>Preview {audience === 'admin' ? 'Admin' : 'Staff'}</strong><small>Synthetic identity</small></div></div></div>
  </div>
}

export function StaffTopBar({ onMenu, onUnavailable, audience }: { audience: 'staff' | 'admin'; onMenu: () => void; onUnavailable: (label: string) => void }) {
  return <header className="staff-top-bar"><AppButton variant="ghost" aria-label="Toggle navigation" onClick={onMenu}><Menu aria-hidden="true" size={22} /></AppButton>
    <AppButton variant="ghost" className="workspace-search" onClick={() => onUnavailable('Workspace search')}><Search size={18} aria-hidden="true" /><span>Search borrowers, equipment, or request ID</span></AppButton>
    <div className="staff-top-bar-right"><span className="preview-tag">Preview</span><AppButton variant="ghost" aria-label="Notifications preview" onClick={() => onUnavailable('Notifications')}><Bell size={20} aria-hidden="true" /></AppButton>
      <AppButton variant="ghost" aria-label="Account preview" onClick={() => onUnavailable('Account')}><span className="account-avatar">{audience === 'admin' ? 'AD' : 'ST'}</span><span className="toolbar-account-copy"><strong>Preview {audience === 'admin' ? 'Admin' : 'Staff'}</strong><small>{audience === 'admin' ? 'Admin' : 'Staff'} preview</small></span><ChevronDown size={14} aria-hidden="true" /></AppButton>
    </div>
  </header>
}

/** Audience chooses visual navigation only; never grants session/role authority. */
export function StaffShell({ children, audience = 'staff', active, dashboardHref, requestsHref, onUnavailable }: {
  children: ReactNode; audience?: 'staff' | 'admin'; active: string; dashboardHref: string; requestsHref: string; onUnavailable: (label: string) => void
}) {
  const [drawer, setDrawer] = useState(false)
  const sidebarOpen = useUIStore(s => s.sidebarOpen)
  const toggleSidebar = useUIStore(s => s.toggleSidebar)
  const sidebarProps = { audience, active, dashboardHref, requestsHref, onUnavailable }
  function menu() {
    if (window.matchMedia('(max-width: 1023px)').matches) setDrawer(true)
    else toggleSidebar()
  }
  return <SidebarProvider className="staff-shell" data-sidebar-open={sidebarOpen}>
    <a className="skip-link" href="#staff-content">Skip to content</a>
    <Sidebar collapsible="none" className="staff-sidebar"><StaffSidebar {...sidebarProps} /></Sidebar>
    <Sheet open={drawer} onOpenChange={setDrawer}><SheetContent side="left" className="staff-navigation-sheet">
      <SheetHeader className="sr-only"><SheetTitle>Staff navigation</SheetTitle><SheetDescription>Choose a workspace destination. Escape closes navigation.</SheetDescription></SheetHeader>
      <StaffSidebar {...sidebarProps} onNavigate={() => setDrawer(false)} />
    </SheetContent></Sheet>
    <div className="staff-workspace"><StaffTopBar audience={audience} onMenu={menu} onUnavailable={onUnavailable} /><main id="staff-content" tabIndex={-1} className="staff-content">{children}</main></div>
  </SidebarProvider>
}
