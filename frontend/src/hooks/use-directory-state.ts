import { useSearchParams } from 'react-router-dom'

/** URL owns directory view state; Query continues to own all returned data. */
export function useDirectoryState() {
  const [params, setParams] = useSearchParams()
  const value = (key: string) => params.get(key) ?? ''
  const page = (key = 'page') => { const parsed = Number(value(key)); return Number.isSafeInteger(parsed) && parsed > 0 && parsed <= 1000000 ? parsed : 1 }
  const update = (changes: Record<string, string>, resetPage = true, replace = false) => setParams(previous => {
    const next = new URLSearchParams(previous)
    if (resetPage) next.delete('page')
    for (const [key, val] of Object.entries(changes)) { if (val && !((key === 'page' || key === 'category_page') && val === '1')) next.set(key, val); else next.delete(key) }
    return next
  }, { replace })
  return { value, page, update }
}
