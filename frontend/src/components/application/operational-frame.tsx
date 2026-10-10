import { ProfileAvatar } from '@/features/profile/components/profile-avatar'
import type { CSSProperties, ReactNode } from 'react'
import { Link } from 'react-router-dom'
import { ChevronUp, LayoutDashboard, ListChecks, Package, Users, ChartNoAxesCombined, Settings, UserRound, LogOut, Moon, Sun, X, type LucideIcon } from 'lucide-react'
import { SidebarProvider, Sidebar, SidebarHeader, SidebarContent, SidebarGroup, SidebarGroupLabel, SidebarMenu, SidebarMenuItem, SidebarMenuButton, SidebarFooter, SidebarTrigger, useSidebar } from '@/components/ui/sidebar'
import { Breadcrumb, BreadcrumbList, BreadcrumbItem, BreadcrumbLink, BreadcrumbPage, BreadcrumbSeparator } from '@/components/ui/breadcrumb'
import { DropdownMenu, DropdownMenuTrigger, DropdownMenuContent, DropdownMenuItem, DropdownMenuSeparator } from '@/components/ui/dropdown-menu'
import { AppButton } from './visual'
import { useUIStore } from '@/stores/ui-store'
import { NotificationBell } from '@/features/notifications/components/notification-bell'
import seal from '@/assets/brand/fsmo-seal-reference.png'

interface Identity { id?: string; name: string; email: string; role: string }
const workspace = [
  ['Dashboard', '/staff/dashboard', LayoutDashboard],
  ['Requests & Borrowings', '/staff/requests', ListChecks],
  ['Inventory', '/staff/inventory', Package],
  ['Borrowers', '/staff/borrowers', Users],
] as const
const administration = [
  ['Reports', '/admin/reports', ChartNoAxesCombined],
  ['Administration', '/admin/administration', Settings],
] as const

function NavigationGroup({ label, items, active }: { label: string; items: readonly (readonly [string, string, LucideIcon])[]; active: string }) {
  const { setOpenMobile } = useSidebar()
  return <SidebarGroup className="fsmo-navigation-group"><SidebarGroupLabel>{label}</SidebarGroupLabel><SidebarMenu>
    {items.map(([name, href, Icon]) => <SidebarMenuItem key={href}><SidebarMenuButton className="fsmo-navigation-link" tooltip={name} isActive={active === name} aria-label={name} render={<Link to={href} aria-current={active === name ? 'page' : undefined} onClick={() => setOpenMobile(false)} />}><Icon aria-hidden="true" /><span>{name}</span></SidebarMenuButton></SidebarMenuItem>)}
  </SidebarMenu></SidebarGroup>
}

function NavigationBody({ identity, active, onSignOut }: { identity: Identity; active: string; onSignOut: () => void }) {
  const { isMobile, state, setOpenMobile } = useSidebar()
  const theme = useUIStore(s => s.theme), toggleTheme = useUIStore(s => s.toggleTheme)
  const collapsed = !isMobile && state === 'collapsed'
  const initials = identity.name.trim().split(/\s+/).slice(0, 2).map(v => v[0]).join('').toUpperCase()
  return <div className="fsmo-navigation-body" data-collapsed={collapsed}>
    <SidebarHeader className="fsmo-brand-header"><Link to="/staff/dashboard" aria-label="FSMO dashboard" className="fsmo-brand-link" onClick={() => setOpenMobile(false)}><img src={seal} alt="FSMO seal" width="44" height="44" /><span><strong>eLabTrack</strong><small>FSMO workspace</small></span></Link>{isMobile && <AppButton className="fsmo-mobile-close" variant="ghost" aria-label="Close navigation" onClick={() => setOpenMobile(false)}><X aria-hidden="true" /></AppButton>}</SidebarHeader>
    <SidebarContent className="fsmo-navigation-content"><nav aria-label={`${identity.role === 'ADMIN' ? 'Admin' : 'Staff'} navigation`}><NavigationGroup label="Workspace" items={identity.role === 'STAFF' ? [...workspace, ['Reports', '/staff/reports', ChartNoAxesCombined] as const] : workspace} active={active} />{identity.role === 'ADMIN' && <NavigationGroup label="Administration" items={administration} active={active} />}</nav></SidebarContent>
    <SidebarFooter className="fsmo-navigation-footer"><SidebarMenu><SidebarMenuItem><SidebarMenuButton className="fsmo-navigation-link" tooltip={theme === 'light' ? 'Dark mode' : 'Light mode'} onClick={toggleTheme} aria-label={`Switch to ${theme === 'light' ? 'dark' : 'light'} theme`}>{theme === 'light' ? <Moon aria-hidden="true" /> : <Sun aria-hidden="true" />}<span>{theme === 'light' ? 'Dark mode' : 'Light mode'}</span></SidebarMenuButton></SidebarMenuItem><SidebarMenuItem>
      <DropdownMenu><DropdownMenuTrigger aria-label="Account actions" render={<SidebarMenuButton size="lg" className="fsmo-account-trigger" />}><span className="fsmo-account-avatar">{identity.id?<ProfileAvatar id={identity.id} name={identity.name} />:initials}</span><span className="fsmo-account-copy"><strong>{identity.name}</strong><small>{identity.role === 'ADMIN' ? 'Admin' : 'Staff'} account</small></span><ChevronUp className="fsmo-account-chevron" aria-hidden="true" /></DropdownMenuTrigger><DropdownMenuContent side={isMobile ? 'top' : 'right'} align="end" className="fsmo-account-menu"><div className="fsmo-menu-identity"><strong>{identity.name}</strong><span>{identity.email}</span></div><DropdownMenuSeparator /><DropdownMenuItem onClick={() => setOpenMobile(false)} render={<Link to="/staff/account" />}><UserRound aria-hidden="true" />Your account</DropdownMenuItem><DropdownMenuSeparator /><DropdownMenuItem variant="destructive" onClick={onSignOut}><LogOut aria-hidden="true" />Sign Out</DropdownMenuItem></DropdownMenuContent></DropdownMenu>
    </SidebarMenuItem></SidebarMenu></SidebarFooter>
  </div>
}

function FrameContent({ active, children }: { active: string; children: ReactNode }) {
  const { open, openMobile, isMobile } = useSidebar()
  return <div className="fsmo-workspace"><header className="staff-top-bar fsmo-header"><SidebarTrigger aria-label="Toggle navigation" aria-expanded={isMobile ? openMobile : open} className="fsmo-sidebar-trigger" /><span className="fsmo-header-divider" aria-hidden="true" /><Breadcrumb aria-label="Page context"><BreadcrumbList><BreadcrumbItem><BreadcrumbLink render={<Link to="/staff/dashboard" />}>FSMO</BreadcrumbLink></BreadcrumbItem><BreadcrumbSeparator /><BreadcrumbItem><BreadcrumbPage>{active}</BreadcrumbPage></BreadcrumbItem></BreadcrumbList></Breadcrumb><NotificationBell className="fsmo-notifications" /></header><main id="staff-content" tabIndex={-1} className="staff-content fsmo-content">{children}</main></div>
}

/** Presentation only: existing route guards, session hooks and server remain authoritative. */
export function OperationalFrame({ identity, active, onSignOut, children }: { identity: Identity; active: string; onSignOut: () => void; children: ReactNode }) {
  const open = useUIStore(s => s.sidebarOpen), setOpen = useUIStore(s => s.setSidebar)
  return <SidebarProvider className="reconstruction-shell" data-sidebar-open={open} open={open} onOpenChange={setOpen} style={{ '--sidebar-width': '260px', '--sidebar-width-icon': '72px' } as CSSProperties}><a className="skip-link" href="#staff-content">Skip to content</a><Sidebar collapsible="icon" className="staff-sidebar fsmo-sidebar"><NavigationBody identity={identity} active={active} onSignOut={onSignOut} /></Sidebar><FrameContent active={active}>{children}</FrameContent></SidebarProvider>
}
