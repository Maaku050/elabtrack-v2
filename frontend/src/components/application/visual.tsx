import type { ComponentProps, ReactNode } from 'react'
import { Check, CircleAlert, Clock, Moon, Sun, X, type LucideIcon } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { useUIStore } from '@/stores/ui-store'
import seal from '@/assets/brand/fsmo-seal-reference.png'

export type Tone = 'primary' | 'success' | 'warning' | 'danger' | 'neutral'

export function AppButton({ className = '', ...props }: ComponentProps<typeof Button>) {
  return <Button className={`app-button ${className}`} {...props} />
}

/** Replaceable supplied reconstruction; not an authenticated official master. */
export function AppBrand({ compact = false, sealSrc = seal }: { compact?: boolean; sealSrc?: string }) {
  return <div className={`app-brand ${compact ? 'app-brand-staff' : ''}`}>
    <img src={sealSrc} alt="FSMO seal" width="44" height="44" />
    <div className="app-brand-copy"><strong>{compact ? 'FSMO' : <>eLab<span>Track</span></>}</strong>
      <small>{compact ? 'Food Service Management Organization' : 'Food Service Management Organization'}</small>
    </div>
  </div>
}

export function ThemeControl({ showLabel = false }: { showLabel?: boolean }) {
  const theme = useUIStore(s => s.theme)
  const toggle = useUIStore(s => s.toggleTheme)
  const Icon = theme === 'light' ? Moon : Sun
  return <AppButton variant="ghost" className="theme-control" title={`Switch to ${theme === 'light' ? 'dark' : 'light'} theme`} aria-label={`Switch to ${theme === 'light' ? 'dark' : 'light'} theme`} onClick={toggle}>
    <Icon aria-hidden="true" size={19} />{showLabel && <span>{theme === 'light' ? 'Dark mode' : 'Light mode'}</span>}
  </AppButton>
}

export function SurfaceCard({ className = '', ...props }: ComponentProps<typeof Card>) {
  return <Card className={`surface-card ${className}`} {...props} />
}

export function PageHeading({ title, description, action }: { title: ReactNode; description: string; action?: ReactNode }) {
  return <div className="page-heading"><div><h1>{title}</h1><p>{description}</p></div>{action}</div>
}

export function SectionHeading({ title, action }: { title: string; action?: ReactNode }) {
  return <div className="section-heading"><h2>{title}</h2>{action}</div>
}

export function MetricCard({ label, value, icon: Icon, tone = 'primary', action }: { label: string; value: string | number; icon: LucideIcon; tone?: Tone; action?: ReactNode }) {
  return <SurfaceCard className="metric-card"><span className="icon-bubble" data-tone={tone}><Icon size={22} aria-hidden="true" /></span>
    <div className="metric-copy"><p>{label}</p><strong>{value}</strong></div>{action && <div className="metric-action">{action}</div>}
  </SurfaceCard>
}

export function ToneBadge({ children, tone = 'neutral', icon: Icon = CircleAlert }: { children: ReactNode; tone?: Tone; icon?: LucideIcon }) {
  return <Badge className="tone-badge" data-tone={tone}><Icon size={14} aria-hidden="true" />{children}</Badge>
}

export function BorrowingStatusBadge({ status }: { status: 'PENDING' | 'CHECKED_OUT' | 'COMPLETED' | 'DENIED' | 'CANCELLED' | 'EXPIRED' }) {
  const states = {
    PENDING: ['Pending', 'warning', Clock], CHECKED_OUT: ['Active', 'success', Check],
    COMPLETED: ['Completed', 'neutral', Check], DENIED: ['Denied', 'danger', X],
    CANCELLED: ['Cancelled', 'neutral', X], EXPIRED: ['Expired', 'neutral', Clock],
  } as const
  const [label, tone, icon] = states[status]
  return <ToneBadge tone={tone} icon={icon}>{label}</ToneBadge>
}

/** Presentation only: caller supplies authoritative quantities/amounts later. */
export function AccountabilityBadge({ physical, replacements }: { physical: number; replacements: number }) {
  return <ToneBadge tone={physical > 0 || replacements > 0 ? 'warning' : 'success'} icon={physical > 0 || replacements > 0 ? CircleAlert : Check}>
    {physical} physically outstanding · {replacements} replacements remaining
  </ToneBadge>
}

export function FineStatusBadge({ outstandingLabel, compact = false }: { outstandingLabel?: string; compact?: boolean }) {
  return <ToneBadge tone={outstandingLabel ? 'danger' : 'success'} icon={outstandingLabel ? CircleAlert : Check}>
    {outstandingLabel ? (compact ? `${outstandingLabel} outstanding` : `Fine outstanding: ${outstandingLabel}`) : 'No overdue fine'}
  </ToneBadge>
}
