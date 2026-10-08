import type { UserRole } from '@/types/common'

export const borrowerPaths = ['/borrower/home', '/borrower/equipment', '/borrower/borrowings', '/borrower/account', '/borrower/notifications'] as const
export const staffPaths = ['/staff/dashboard', '/staff/requests', '/staff/inventory', '/staff/borrowers', '/staff/account', '/staff/notifications'] as const
export const adminPaths = ['/admin/reports', '/admin/administration'] as const

export function canVisit(role: UserRole, pathname: string): boolean {
  const uuid = '[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}'
  if (role === 'BORROWER' && new RegExp('^/borrower/equipment/' + uuid + '$', 'i').test(pathname)) return true
  if ((role === 'STAFF' || role === 'ADMIN') && (['/staff/inventory/new', '/staff/inventory/categories'].includes(pathname) || new RegExp('^/staff/inventory/' + uuid + '(/edit|/adjust)?$', 'i').test(pathname))) return true
  if (role === 'ADMIN' && new RegExp('^/admin/inventory/' + uuid + '/reconcile$', 'i').test(pathname)) return true
  if ((role === 'STAFF' || role === 'ADMIN') && new RegExp('^/staff/borrowers/' + uuid + '$', 'i').test(pathname)) return true
  if (role === 'ADMIN' && (['/staff/borrowers/new', '/admin/borrowers/bulk', '/admin/administration/accounts', '/admin/administration/accounts/new'].includes(pathname) || new RegExp('^/admin/administration/accounts/' + uuid + '$', 'i').test(pathname))) return true
  return (role === 'BORROWER' && (pathname === '/borrower/terms' || borrowerPaths.some(p => p === pathname))) ||
    ((role === 'STAFF' || role === 'ADMIN') && staffPaths.some(p => p === pathname)) ||
    (role === 'ADMIN' && adminPaths.some(p => p === pathname))
}
export function homeForRole(role: UserRole): string {
  return role === 'BORROWER' ? '/borrower/home' : '/staff/dashboard'
}
// Exact internal route allowlist: no protocol, origin, traversal or preview target.
export function loginDestination(role: UserRole, requested: unknown): string {
  return typeof requested === 'string' && canVisit(role, requested) ? requested : homeForRole(role)
}
