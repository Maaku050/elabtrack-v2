import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Package, Clock, TriangleAlert, Coins, Search, ArrowRight, ChevronRight, ClipboardPlus, ListChecks, UserRound, Megaphone, ShoppingCart, SlidersHorizontal, ArrowDownUp } from 'lucide-react'
import { Checkbox } from '@/components/ui/checkbox'
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription } from '@/components/ui/sheet'
import { EmptyState } from '@/components/feedback/states'
import { AppButton, PageHeading, MetricCard, SectionHeading, BorrowingStatusBadge } from '@/components/application/visual'
import { EquipmentCard, EquipmentThumbnail } from '@/components/application/equipment'
import { SearchField } from '@/components/application/data-table'
import { PreviewDialog } from '@/components/application/feedback'
import { equipment, categories } from './fixtures'
import { BorrowerPreviewFrame } from './preview-context'

export function BorrowerHomePreview() {
  const [disclosure, setDisclosure] = useState<string | null>(null)
  const activities = [
    { item: equipment[1], description: '2 sets borrowed', status: 'CHECKED_OUT', time: 'Today, 9:30 AM' },
    { item: equipment[3], description: '1 pc requested', status: 'PENDING', time: '08 Oct, 8:15 AM' },
    { item: equipment[5], description: '6 pcs returned in good condition', status: 'COMPLETED', time: 'Yesterday' },
  ] as const
  return <BorrowerPreviewFrame active="Home"><PageHeading title={<>Good morning, <span className="primary-text">Alex!</span></>} description="Here’s what’s happening with your borrowings." />
    <div className="borrower-metrics">
      <MetricCard label="Active borrowings" value={2} icon={Package} />
      <MetricCard label="Pending requests" value={1} icon={Clock} tone="warning" />
      <MetricCard label="Overdue" value={0} icon={TriangleAlert} tone="danger" />
      <MetricCard label="Outstanding fine" value="PHP 0" icon={Coins} tone="warning" />
    </div>
    <AppButton className="browse-equipment" nativeButton={false} render={<Link to="/__preview/borrower/equipment" />}><Search size={20} aria-hidden="true" />Browse Equipment<ArrowRight size={20} aria-hidden="true" /></AppButton>
    <section className="home-section" aria-label="Recent activity"><SectionHeading title="Recent Activity" action={<AppButton variant="ghost" onClick={() => setDisclosure('Activity')}>View All <ChevronRight size={16} aria-hidden="true" /></AppButton>} />
      <div className="activity-list">{activities.map(({ item, description, status, time }) => <AppButton key={item.id} variant="ghost" className="activity-item" onClick={() => setDisclosure(item.name)}>
        <EquipmentThumbnail src={item.image} name={item.name} /><div className="activity-copy"><strong>{item.name}</strong><p>{description}</p></div><div className="activity-status"><BorrowingStatusBadge status={status} /><small>{time}</small></div><ChevronRight className="activity-chevron" size={16} aria-hidden="true" />
      </AppButton>)}</div>
    </section>
    <section className="home-section" aria-label="Quick actions"><SectionHeading title="Quick Actions" /><div className="quick-actions">
      <AppButton variant="outline" nativeButton={false} render={<Link to="/__preview/borrower/equipment" />}><span className="icon-bubble" data-tone="primary"><ClipboardPlus size={24} aria-hidden="true" /></span><strong>Start Request</strong><small>Browse equipment</small></AppButton>
      <AppButton variant="outline" onClick={() => setDisclosure('My Borrowings')}><span className="icon-bubble" data-tone="primary"><ListChecks size={24} aria-hidden="true" /></span><strong>My Borrowings</strong><small>Track your items</small></AppButton>
      <AppButton variant="outline" onClick={() => setDisclosure('Account')}><span className="icon-bubble" data-tone="primary"><UserRound size={24} aria-hidden="true" /></span><strong>Account</strong><small>View your profile</small></AppButton>
    </div></section>
    <aside className="reservation-notice"><Megaphone size={23} aria-hidden="true" /><p>Pending requests reserve equipment for <strong>24 hours</strong>. Visit FSMO for approval and physical release.</p></aside>
    <PreviewDialog title={`${disclosure ?? 'Activity'} preview`} open={disclosure !== null} onOpenChange={open => { if (!open) setDisclosure(null) }}>Read-only synthetic activity. Borrowing history and account functionality belong to later phases. Staff/Admin alone records authoritative returns.</PreviewDialog>
  </BorrowerPreviewFrame>
}

export function EquipmentCatalogPreview() {
  const [query, setQuery] = useState('')
  const [category, setCategory] = useState('All')
  const [availableOnly, setAvailableOnly] = useState(false)
  const [sort, setSort] = useState(false)
  const [filtersOpen, setFiltersOpen] = useState(false)
  const [summaryOpen, setSummaryOpen] = useState(false)
  const [selected, setSelected] = useState<Record<string, number>>({ spoon: 2, chafing: 1 })
  const selectedItems = equipment.filter(item => selected[item.id] > 0)
  const units = selectedItems.reduce((sum, item) => sum + selected[item.id], 0)
  const filtered = equipment.filter(item => (category === 'All' || item.category === category) && (!availableOnly || item.available > 0) && `${item.name} ${item.category} ${item.tags.join(' ')}`.toLowerCase().includes(query.toLowerCase()))
  const visible = sort ? [...filtered].sort((a, b) => a.name.localeCompare(b.name)) : filtered
  function clearFilters() { setQuery(''); setCategory('All'); setAvailableOnly(false) }
  const summary = <><p>Local selection only · not reserved.</p><ul className="preview-selection-list">{selectedItems.map(item => <li key={item.id}><span>{item.name}</span><strong>{selected[item.id]} units</strong></li>)}</ul><p>{selectedItems.length} equipment types · {units} units. Request submission is not implemented in this phase.</p></>
  return <BorrowerPreviewFrame active="Equipment" headerAction={<AppButton variant="ghost" className="cart-header-button" aria-label={`Selected equipment preview, ${selectedItems.length} types`} onClick={() => setSummaryOpen(true)}><ShoppingCart size={25} aria-hidden="true" /><span className="count-bubble">{selectedItems.length}</span></AppButton>}>
    <PageHeading title={<>Equipment <span className="primary-text">Catalog</span></>} description="Browse available items and build your request." />
    <div className="catalog-toolbar"><SearchField label="Search equipment" placeholder="Search equipment…" value={query} onChange={setQuery} />
      <AppButton variant="outline" onClick={() => setFiltersOpen(true)} aria-label="Equipment filters"><SlidersHorizontal size={18} aria-hidden="true" /><span>Filters</span></AppButton>
      <AppButton variant="outline" onClick={() => setSort(!sort)} aria-pressed={sort} aria-label="Sort equipment alphabetically"><ArrowDownUp size={18} aria-hidden="true" /><span>Sort</span></AppButton>
    </div>
    <div className="category-chips" role="group" aria-label="Equipment categories">{categories.map(label => <AppButton key={label} variant={category === label ? 'default' : 'secondary'} aria-pressed={category === label} onClick={() => setCategory(label)}>{label}</AppButton>)}</div>
    <div className="catalog-grid">{visible.map(item => <EquipmentCard key={item.id} item={item} selected={selected[item.id] ?? 0} onSelect={quantity => setSelected(current => ({ ...current, [item.id]: quantity }))} />)}</div>
    {visible.length === 0 && <EmptyState title="No matching equipment" description="Try another search or clear the filters. Your local selection is preserved." action={<AppButton variant="outline" onClick={clearFilters}>Clear filters</AppButton>} />}
    <div className="cart-dock"><div className="cart-dock-icon"><ShoppingCart size={27} aria-hidden="true" /><span className="count-bubble">{selectedItems.length}</span></div>
      <div><strong>Request Cart</strong><p aria-live="polite">{selectedItems.length} types · {units} units selected</p></div><AppButton disabled={selectedItems.length === 0} onClick={() => setSummaryOpen(true)}>Review Request<ChevronRight size={18} aria-hidden="true" /></AppButton>
    </div>
    <Sheet open={filtersOpen} onOpenChange={setFiltersOpen}><SheetContent className="app-filter-sheet"><SheetHeader><SheetTitle>Equipment filters</SheetTitle><SheetDescription>Filter synthetic catalog fixtures. Selection stays local.</SheetDescription></SheetHeader>
      <label className="filter-checkbox"><Checkbox checked={availableOnly} onCheckedChange={checked => setAvailableOnly(checked === true)} />Available equipment only</label>
      <div className="filter-actions"><AppButton variant="outline" onClick={clearFilters}>Clear filters</AppButton><AppButton onClick={() => setFiltersOpen(false)}>Done</AppButton></div>
    </SheetContent></Sheet>
    <PreviewDialog title="Selected equipment preview" open={summaryOpen} onOpenChange={setSummaryOpen}>{summary}</PreviewDialog>
  </BorrowerPreviewFrame>
}
