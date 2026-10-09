/** Bounded server page links; endpoints are included, never all pages for large totals. */
export function pageRange(page: number, total: number): (number | 'ellipsis')[] {
  const pages = new Set([1, total, page - 1, page, page + 1])
  if (page <= 3) { pages.add(2); pages.add(3); pages.add(4) }
  if (page >= total - 2) { pages.add(total - 1); pages.add(total - 2); pages.add(total - 3) }
  const sorted = [...pages].filter(p => p >= 1 && p <= total).sort((a, b) => a - b)
  const result: (number | 'ellipsis')[] = []
  sorted.forEach((p, i) => { if (i && p - sorted[i - 1] > 1) result.push('ellipsis'); result.push(p) })
  return result
}
