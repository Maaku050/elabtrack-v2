import { describe, it, expect } from 'vitest'
import { accountFormSchema } from './account.schema'
import { canVisit, loginDestination } from '@/features/auth/navigation'
const student = { name: 'Student Example', email: 's@students.example.invalid', borrower_type: 'STUDENT', student_id: '0028366', course: '', contact_number: '' }
describe('account onboarding validation', () => {
 it('retains leading zeroes and requires Student ID', () => { expect(accountFormSchema.parse(student).student_id).toBe('0028366'); expect(accountFormSchema.safeParse({ ...student, student_id: '' }).success).toBe(false) })
 it('supports individual Faculty without Student ID', () => { expect(accountFormSchema.safeParse({ ...student, borrower_type: 'FACULTY', student_id: '', email: 'faculty@external.example.invalid' }).success).toBe(true) })
 it('does not accept an authorization role as a borrower type', () => { expect(accountFormSchema.safeParse({ ...student, borrower_type: 'ADMIN' }).success).toBe(false) })
 it('limits account-management deep links by role and exact UUID path', () => { const id = 'aaaaaaaa-aaaa-4aaa-aaaa-aaaaaaaaaaaa'; expect(canVisit('STAFF', '/staff/borrowers/' + id)).toBe(true); for (const role of ['STAFF', 'BORROWER'] as const) { expect(canVisit(role, '/staff/borrowers/new')).toBe(false); expect(canVisit(role, '/admin/borrowers/bulk')).toBe(false); expect(canVisit(role, '/admin/administration/accounts')).toBe(false) } expect(loginDestination('ADMIN', '/admin/administration/accounts/' + id)).toBe('/admin/administration/accounts/' + id); expect(loginDestination('ADMIN', '/staff/borrowers/' + id + '/../../admin')).toBe('/staff/dashboard') })
})
