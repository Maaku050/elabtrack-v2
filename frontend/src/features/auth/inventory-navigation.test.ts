import { describe, expect, it } from 'vitest'
import { canVisit, loginDestination } from './navigation'
const id = '12345678-1234-4234-8234-123456789abc'
describe('inventory navigation capabilities', () => {
 it('separates Borrower discovery, Staff maintenance and Admin correction', () => {
  expect(canVisit('BORROWER', `/borrower/equipment/${id}`)).toBe(true)
  expect(canVisit('BORROWER', `/staff/inventory/${id}/edit`)).toBe(false)
  expect(canVisit('STAFF', `/staff/inventory/${id}/adjust`)).toBe(true)
  expect(canVisit('STAFF', `/admin/inventory/${id}/reconcile`)).toBe(false)
  expect(canVisit('ADMIN', `/admin/inventory/${id}/reconcile`)).toBe(true)
 })
 it('rejects external, traversal and extra-suffix login targets', () => {
  for (const target of [`//example.invalid/staff/inventory/${id}`, `/staff/inventory/${id}/edit/extra`, `/staff/inventory/../admin`, `/staff/inventory/${id}%2fedit`, `/__preview/borrower/equipment`]) expect(loginDestination('ADMIN', target)).toBe('/staff/dashboard')
 })
})
