import type { MouseEvent } from 'react'
import { useLocation } from 'react-router-dom'
import { Pagination, PaginationContent, PaginationEllipsis, PaginationItem, PaginationLink, PaginationNext, PaginationPrevious } from '@/components/ui/pagination'

import { pageRange } from '@/lib/page-range'

export function ServerPagination({ page, totalPages, total, perPage = 25, count, hasNext = false, pending = false, label, pageKey = 'page', onPage }: {
  page: number; totalPages?: number; total?: number; perPage?: number; count: number; hasNext?: boolean; pending?: boolean; label: string; pageKey?: string; onPage: (page: number) => void
}) {
  const location = useLocation()
  const last = totalPages === undefined ? undefined : Math.max(1, totalPages)
  const previousDisabled = page <= 1 || pending, nextDisabled = pending || (last === undefined ? !hasNext : page >= last)
  function link(target: number, disabled = false) {
    const params = new URLSearchParams(location.search); if (target === 1) params.delete(pageKey); else params.set(pageKey, String(target))
    return { role: 'link', href: disabled ? undefined : `${location.pathname}${params.size ? '?' + params.toString() : ''}`, 'aria-disabled': disabled || undefined, tabIndex: disabled ? -1 : undefined, onClick: (event: MouseEvent<HTMLAnchorElement>) => {
      if (disabled) { event.preventDefault(); return }
      if (event.button === 0 && !event.metaKey && !event.ctrlKey && !event.altKey && !event.shiftKey) { event.preventDefault(); onPage(target) }
    } }
  }
  const start = count ? (page - 1) * perPage + 1 : 0, end = count ? start + count - 1 : 0
  return <div className="directory-pagination"><p role="status">{total === undefined ? `Page ${page}` : `${start}–${end} of ${total} ${label.toLowerCase()}`}<span>{last === undefined ? 'Total pages unavailable' : page > last ? `Page ${page} (outside current results)` : `Page ${page} of ${last}`}</span></p><Pagination aria-label={`${label} pagination`}><PaginationContent><PaginationItem><PaginationPrevious {...link(page - 1, previousDisabled)} /></PaginationItem>{last !== undefined && pageRange(Math.min(page, last), last).map((value, i) => <PaginationItem className="directory-page-number" key={`${value}-${i}`}>{value === 'ellipsis' ? <PaginationEllipsis /> : <PaginationLink aria-label={`Go to page ${value}`} isActive={value === page} {...link(value, pending)}>{value}</PaginationLink>}</PaginationItem>)}<PaginationItem><PaginationNext {...link(page + 1, nextDisabled)} /></PaginationItem></PaginationContent></Pagination></div>
}
