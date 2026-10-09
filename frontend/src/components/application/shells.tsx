import { useState, type ReactNode } from 'react'
import { Link } from 'react-router-dom'
import { Bell, LayoutDashboard, ListChecks, Package, Users, ChartNoAxesCombined, Settings, Menu, Search, ChevronDown, ChevronRight, UserRound, LogOut, type LucideIcon } from 'lucide-react'
import { SidebarProvider, Sidebar } from '@/components/ui/sidebar'
import { Sheet, SheetContent, SheetTitle, SheetDescription, SheetHeader } from '@/components/ui/sheet'
import { DropdownMenu, DropdownMenuTrigger, DropdownMenuContent, DropdownMenuItem } from '@/components/ui/dropdown-menu'
import { useUIStore } from '@/stores/ui-store'
import { AppBrand, AppButton, ThemeControl } from './visual'

export interface NavigationDestination { label: string; icon: LucideIcon; href?: string; onSelect?: () => void }
interface WorkspaceIdentity { name: string; email: string }
interface WorkspaceOptions { identity?: WorkspaceIdentity; destinationHrefs?: Record<string, string>; footerAction?: ReactNode; onSignOut?: () => void }

function Destination({ item, active, onNavigate }: { item: NavigationDestination; active: boolean; onNavigate?: () => void }) {
  const Icon = item.icon
  const content = <><Icon size={21} aria-hidden="true" /><span>{item.label}</span></>
  return item.href
    ? <Link title={item.label} aria-label={item.label} aria-current={active ? 'page' : undefined} to={item.href} onClick={onNavigate}>{content}</Link>
    : <button type="button" title={item.label} aria-label={item.label} onClick={() => { onNavigate?.(); item.onSelect?.() }}>{content}</button>
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

export function StaffSidebar({ audience, active, dashboardHref, requestsHref, onUnavailable, onNavigate, identity, destinationHrefs, footerAction }: {
  audience: 'staff' | 'admin'; active: string; dashboardHref: string; requestsHref: string; onUnavailable: (label: string) => void; onNavigate?: () => void
} & WorkspaceOptions) {
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
  for (const item of [...operational, ...administration]) if (destinationHrefs?.[item.label]) item.href = destinationHrefs[item.label]
  return <div className="staff-sidebar-content"><div className="staff-brand"><AppBrand compact /></div>
    <nav aria-label={`${audience === 'admin' ? 'Admin' : 'Staff'} navigation`}>
      {operational.map(item => <Destination key={item.label} item={item} active={active === item.label} onNavigate={onNavigate} />)}
      {audience === 'admin' && <div className="admin-navigation">{administration.map(item => <Destination key={item.label} item={item} active={active === item.label} onNavigate={onNavigate} />)}</div>}
    </nav>
    <div className="staff-sidebar-footer"><ThemeControl showLabel />{footerAction}<div className={`staff-identity${identity ? " product-identity" : ""}`}><span className="account-avatar">{audience === 'admin' ? 'AD' : 'ST'}</span><div><strong>{identity?.name ?? `Preview ${audience === 'admin' ? 'Admin' : 'Staff'}`}</strong><small>{identity ? (audience === 'admin' ? 'Admin account' : 'Staff account') : 'Synthetic identity'}</small></div></div></div>
  </div>
}

export function StaffTopBar({ onMenu, onUnavailable, audience, identity, active, onSignOut }: { audience: 'staff' | 'admin'; onMenu: () => void; onUnavailable: (label: string) => void; identity?: WorkspaceIdentity; active: string; onSignOut?: () => void }) {
  return <header className="staff-top-bar"><AppButton variant="ghost" aria-label="Toggle navigation" onClick={onMenu}><Menu aria-hidden="true" size={22} /></AppButton>
    {identity ? <nav className="workspace-context" aria-label="Page context"><span>FSMO</span><ChevronRight size={14} aria-hidden="true" /><span>{active}</span></nav> : <AppButton variant="ghost" className="workspace-search" onClick={() => onUnavailable('Workspace search')}><Search size={18} aria-hidden="true" /><span>Search borrowers, equipment, or request ID</span></AppButton>}
    <div className="staff-top-bar-right">{!identity && <span className="preview-tag">Preview</span>}<AppButton variant="ghost" title={identity ? 'Notifications — not available yet' : 'Notifications preview'} aria-label={identity ? 'Notifications' : 'Notifications preview'} onClick={() => onUnavailable('Notifications')}><Bell size={20} aria-hidden="true" /></AppButton>
      <DropdownMenu><DropdownMenuTrigger render={<AppButton variant="ghost" />} aria-label={identity ? 'Account actions' : 'Account preview'}><span className="account-avatar">{audience === 'admin' ? 'AD' : 'ST'}</span><span className={`toolbar-account-copy${identity ? " product-identity" : ""}`}><strong>{identity?.name ?? `Preview ${audience === 'admin' ? 'Admin' : 'Staff'}`}</strong><small>{audience === 'admin' ? 'Admin' : 'Staff'} {identity ? 'account' : 'preview'}</small></span><ChevronDown size={14} aria-hidden="true" /></DropdownMenuTrigger><DropdownMenuContent align="end" className="workspace-account-menu"><DropdownMenuItem onClick={() => onUnavailable('Account')}><UserRound aria-hidden="true" />Your account</DropdownMenuItem>{onSignOut && <DropdownMenuItem onClick={onSignOut}><LogOut aria-hidden="true" />Sign Out</DropdownMenuItem>}</DropdownMenuContent></DropdownMenu>
    </div>
  </header>
}

/** Audience chooses visual navigation only; never grants session/role authority. */
export function StaffShell({ children, audience = 'staff', active, dashboardHref, requestsHref, onUnavailable, identity, destinationHrefs, footerAction, onSignOut }: {
  children: ReactNode; audience?: 'staff' | 'admin'; active: string; dashboardHref: string; requestsHref: string; onUnavailable: (label: string) => void
} & WorkspaceOptions) {
  const [drawer, setDrawer] = useState(false)
  const sidebarOpen = useUIStore(s => s.sidebarOpen)
  const toggleSidebar = useUIStore(s => s.toggleSidebar)
  const sidebarProps = { audience, active, dashboardHref, requestsHref, onUnavailable, identity, destinationHrefs, footerAction }
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
    <div className="staff-workspace"><StaffTopBar active={active} onSignOut={onSignOut} audience={audience} identity={identity} onMenu={menu} onUnavailable={onUnavailable} /><main id="staff-content" tabIndex={-1} className="staff-content">{children}</main></div>
  </SidebarProvider>
}
