import { useState } from 'react'
import { Clock, Package, CalendarDays, TriangleAlert, ArrowRight, ChevronDown, Check, CircleAlert, SlidersHorizontal, X, ChartNoAxesColumn, Activity } from 'lucide-react'
import { Checkbox } from '@/components/ui/checkbox'
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table'
import { EmptyState } from '@/components/feedback/states'
import { AppButton, PageHeading, MetricCard, SurfaceCard, FineStatusBadge } from '@/components/application/visual'
import { StaffShell } from '@/components/application/shells'
import { SearchField, DataTableShell, PaginationControls } from '@/components/application/data-table'
import { EquipmentThumbnail } from '@/components/application/equipment'
import { PreviewDialog } from '@/components/application/feedback'
import { requestTrends, stockSummary, pendingRequests, equipment } from './fixtures'
import { PreviewNotice } from './preview-context'

function RequestTrend() {
  return <figure className="request-trend"><svg viewBox="0 0 360 165" role="img" aria-labelledby="trend-title trend-desc">
    <title id="trend-title">Synthetic request trends over seven days</title><desc id="trend-desc">Daily issued, pending and denied counts. Exact values follow in a text table.</desc>
    {[0, 5, 10, 15].map(value => <g key={value}><line x1="28" x2="355" y1={140 - value * 8} y2={140 - value * 8} className="chart-grid" /><text x="20" y={144 - value * 8} textAnchor="end">{value}</text></g>)}
    {requestTrends.map((point, index) => <g key={point.day}>
      <rect className="chart-issued" x={39 + index * 45} y={140 - point.issued * 8} width="22" height={point.issued * 8} rx="2" />
      <rect className="chart-pending" x={39 + index * 45} y={140 - (point.issued + point.pending) * 8} width="22" height={point.pending * 8} />
      <rect className="chart-denied" x={39 + index * 45} y={140 - (point.issued + point.pending + point.denied) * 8} width="22" height={point.denied * 8} rx="2" />
      <text x={50 + index * 45} y="162" textAnchor="middle">{point.day}</text>
    </g>)}
  </svg>
    <table className="sr-only"><caption>Synthetic daily request counts</caption><thead><tr><th>Day</th><th>Issued</th><th>Pending</th><th>Denied</th></tr></thead><tbody>{requestTrends.map(point => <tr key={point.day}><th>{point.day}</th><td>{point.issued}</td><td>{point.pending}</td><td>{point.denied}</td></tr>)}</tbody></table>
  </figure>
}

function EquipmentStatus() {
  return <div className="equipment-status-summary"><svg viewBox="0 0 140 140" role="img" aria-label="124 total tracked units: 82 available, 28 checked out, 8 damaged held, 6 reserved">
    {stockSummary.map((item, index) => {
      const length = item.count / 124 * 314.16
      const start = stockSummary.slice(0, index).reduce((sum, prior) => sum + prior.count / 124 * 314.16, 0)
      return <circle key={item.label} className="stock-arc" data-tone={item.tone} cx="70" cy="70" r="50" fill="none" strokeWidth="18" strokeDasharray={`${length} ${314.16 - length}`} strokeDashoffset={-start} transform="rotate(-90 70 70)" />
    })}
    <text x="70" y="71" textAnchor="middle" className="donut-total">124</text><text x="70" y="88" textAnchor="middle">Total tracked</text>
  </svg><ul>{stockSummary.map(item => <li key={item.label}><span className="legend-dot" data-tone={item.tone} /><span>{item.label}</span><strong>{item.count}</strong></li>)}</ul></div>
}

export function StaffDashboardPreview({ audience = 'staff' }: { audience?: 'staff' | 'admin' }) {
  const [disclosure, setDisclosure] = useState<string | null>(null)
  const base = audience === 'admin' ? '/__preview/admin' : '/__preview/staff'
  const summaries = [
    { label: 'Pending requests', value: 8, icon: Clock, tone: 'primary' }, { label: 'Active borrowings', value: 12, icon: Package, tone: 'success' },
    { label: 'Due today', value: 5, icon: CalendarDays, tone: 'warning' }, { label: 'Overdue', value: 3, icon: TriangleAlert, tone: 'danger' },
  ] as const
  const activities = [
    { label: 'Request submitted', detail: 'Alex Rivera · 4 units', time: '5 min ago', icon: Clock, tone: 'primary' },
    { label: 'Return recorded', detail: 'Jordan Cruz · 3 good units', time: '18 min ago', icon: Check, tone: 'success' },
    { label: 'Borrowing overdue', detail: 'Casey Santos · original due passed', time: '35 min ago', icon: CircleAlert, tone: 'danger' },
    { label: 'Equipment issued', detail: 'Taylor Reyes · 2 units', time: '1 hour ago', icon: Package, tone: 'primary' },
    { label: 'Request denied', detail: 'Morgan Diaz · reason recorded', time: '2 hours ago', icon: X, tone: 'danger' },
  ] as const
  return <StaffShell audience={audience} active="Dashboard" dashboardHref={`${base}/dashboard`} requestsHref="/__preview/staff/requests" onUnavailable={setDisclosure}>
    <PageHeading title="Dashboard" description="Monitor requests, borrowings, and equipment at a glance." action={<AppButton variant="outline" onClick={() => setDisclosure('Date range')}><CalendarDays size={17} aria-hidden="true" />08 Oct 2026<ChevronDown size={15} aria-hidden="true" /></AppButton>} />
    <div className="staff-metrics">{summaries.map(item => <MetricCard key={item.label} {...item} action={<AppButton variant="ghost" aria-label={`View ${item.label.toLowerCase()} preview`} onClick={() => setDisclosure(item.label)}>View<ArrowRight size={14} aria-hidden="true" /></AppButton>} />)}</div>
    <div className="dashboard-panels"><SurfaceCard className="dashboard-panel trend-panel"><div className="panel-heading"><h2><ChartNoAxesColumn size={17} aria-hidden="true" />Request Trends</h2><AppButton variant="ghost" onClick={() => setDisclosure('Request date range')}>Last 7 days<ChevronDown size={14} aria-hidden="true" /></AppButton></div><div className="chart-legend"><span>Issued</span><span>Pending</span><span>Denied</span></div><RequestTrend /></SurfaceCard>
      <SurfaceCard className="dashboard-panel"><div className="panel-heading"><h2><Package size={17} aria-hidden="true" />Equipment Status</h2><span>Units</span></div><EquipmentStatus /></SurfaceCard>
      <SurfaceCard className="dashboard-panel"><div className="panel-heading"><h2><Activity size={17} aria-hidden="true" />Recent Activities</h2><AppButton variant="ghost" onClick={() => setDisclosure('Activity')}>View All</AppButton></div>
        <ul className="staff-activities">{activities.map(({ label, detail, time, icon: Icon, tone }) => <li key={label}><span className="icon-bubble" data-tone={tone}><Icon size={16} aria-hidden="true" /></span><div><strong>{label}</strong><p>{detail}</p></div><time>{time}</time></li>)}</ul>
      </SurfaceCard>
    </div><PreviewNotice staff />
    <PreviewDialog title={`${disclosure ?? 'Workspace'} preview`} open={disclosure !== null} onOpenChange={open => { if (!open) setDisclosure(null) }}>This visual shell uses local synthetic data. Dashboard queries, account controls and business workflows are deferred. The displayed audience does not change the authenticated session or grant permissions.</PreviewDialog>
  </StaffShell>
}

export function PendingRequestsPreview() {
  const [query, setQuery] = useState('')
  const [fineOnly, setFineOnly] = useState(false)
  const [page, setPage] = useState(1)
  const [selectedRows, setSelectedRows] = useState<Set<string>>(() => new Set())
  const [disclosure, setDisclosure] = useState<string | null>(null)
  const [review, setReview] = useState<(typeof pendingRequests)[number] | null>(null)
  const filtered = pendingRequests.filter(row => `${row.reference} ${row.name} ${row.program}`.toLowerCase().includes(query.toLowerCase()) && (!fineOnly || row.fine !== null))
  const pages = Math.max(1, Math.ceil(filtered.length / 5)), currentPage = Math.min(page, pages)
  const visible = filtered.slice((currentPage - 1) * 5, currentPage * 5)
  function toggleRow(reference: string, checked: boolean) { setSelectedRows(current => { const next = new Set(current); if (checked) next.add(reference); else next.delete(reference); return next }) }
  const groups = [['Pending', 8], ['Active', 12], ['Completed', 10], ['Denied', 4], ['Cancelled', 3], ['Expired', 5], ['All', 42]] as const
  return <StaffShell active="Requests & Borrowings" dashboardHref="/__preview/staff/dashboard" requestsHref="/__preview/staff/requests" onUnavailable={setDisclosure}>
    <PageHeading title="Pending Requests" description="Review requests and confirm approval while physically handing over equipment." action={<div className="pending-toolbar"><AppButton variant="outline" aria-label="Filter requests with outstanding fines" aria-pressed={fineOnly} onClick={() => { setFineOnly(!fineOnly); setPage(1) }}><SlidersHorizontal size={17} aria-hidden="true" />Filters</AppButton><SearchField label="Search pending requests" placeholder="Search requests…" value={query} onChange={next => { setQuery(next); setPage(1) }} /></div>} />
    <div className="request-status-filters" role="group" aria-label="Borrowing lifecycle filters">{groups.map(([label, count]) => <AppButton key={label} variant={label === 'Pending' ? 'default' : 'secondary'} aria-pressed={label === 'Pending'} onClick={() => { if (label !== 'Pending') setDisclosure(`${label} filter`) }}>{label}<span>{count}</span></AppButton>)}</div>
    {fineOnly && <p className="filter-summary">Showing requests with an outstanding fine. This context does not automatically block borrowing.</p>}
    <DataTableShell label="Pending requests table; scroll horizontally for additional columns" footer={<PaginationControls page={currentPage} pages={pages} summary={filtered.length ? `Showing ${(currentPage - 1) * 5 + 1}–${Math.min(currentPage * 5, filtered.length)} of ${filtered.length} synthetic requests` : 'No matching requests'} onPageChange={setPage} />}>
      <Table className="pending-table"><TableHeader><TableRow>
        <TableHead className="selection-cell"><Checkbox aria-label="Select visible preview rows" checked={visible.length > 0 && visible.every(row => selectedRows.has(row.reference))} onCheckedChange={checked => visible.forEach(row => toggleRow(row.reference, checked === true))} /></TableHead>
        {['Request ID', 'Borrower', 'Submitted', 'Context / program', 'Items', 'Total units', 'Requested due', 'Accountability', 'Actions'].map(label => <TableHead key={label} scope="col">{label}</TableHead>)}
      </TableRow></TableHeader><TableBody>{visible.map(row => <TableRow key={row.reference} data-state={selectedRows.has(row.reference) ? 'selected' : undefined}>
        <TableCell className="selection-cell"><Checkbox aria-label={`Select ${row.reference}`} checked={selectedRows.has(row.reference)} onCheckedChange={checked => toggleRow(row.reference, checked === true)} /></TableCell>
        <TableCell><strong className="primary-text">{row.reference}</strong></TableCell>
        <TableCell><div className="borrower-cell"><span className="account-avatar">{row.initials}</span><div><strong>{row.name}</strong><small>{row.category} · {row.program}</small></div></div></TableCell>
        <TableCell>{row.submitted}<small className="cell-secondary">{row.time}</small></TableCell><TableCell>{row.context}<small className="cell-secondary">{row.program}</small></TableCell>
        <TableCell><div className="request-item-images">{row.itemIds.map(id => { const item = equipment.find(item => item.id === id)!; return <EquipmentThumbnail key={id} src={item.image} name={item.name} /> })}</div></TableCell>
        <TableCell><strong>{row.units}</strong></TableCell><TableCell>{row.due}<small className="cell-secondary">Asia/Manila</small></TableCell>
        <TableCell><FineStatusBadge compact outstandingLabel={row.fine ?? undefined} /></TableCell><TableCell><AppButton aria-label={`Review ${row.reference}`} onClick={() => setReview(row)}>Review<ArrowRight size={14} aria-hidden="true" /></AppButton></TableCell>
      </TableRow>)}</TableBody></Table>
      {visible.length === 0 && <EmptyState title="No matching requests" description="Clear the filters to view the synthetic pending queue." action={<AppButton variant="outline" onClick={() => { setQuery(''); setFineOnly(false); setPage(1) }}>Clear filters</AppButton>} />}
    </DataTableShell><PreviewNotice staff />
    <PreviewDialog title="Request review preview" open={review !== null} onOpenChange={open => { if (!open) setReview(null) }}>{review && <><h2>{review.reference} · {review.name}</h2><p>Pending · {review.units} units reserved in this synthetic example. Requested due: {review.due}, Asia/Manila.</p><p>Submission starts a 24-hour reservation. Staff/Admin approval and physical release are one transition. This preview cannot approve, deny, issue or record a return.</p><FineStatusBadge outstandingLabel={review.fine ?? undefined} /></>}</PreviewDialog>
    <PreviewDialog title={`${disclosure ?? 'Feature'} preview`} open={disclosure !== null} onOpenChange={open => { if (!open) setDisclosure(null) }}>The representative table contains pending fixtures only. Other lifecycle queues and business features are not implemented.</PreviewDialog>
  </StaffShell>
}
