import { useRef } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useForm, useWatch } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { accountFormSchema, type Values } from '../schemas/account.schema'
import { Input } from '@/components/ui/input'
import { AppButton } from '@/components/application/visual'
import { ManagementPage, ManagementCard, FormSection, FormField, FormActions, Notice } from '@/components/application/management'
import { NativeSelect } from '@/components/ui/native-select'
import { apiErrorMessage, normalizeApiError } from '@/lib/api-error'
import { accountsApi } from '../api/accounts.api'
import { useAccountCommand, useProvisioningPolicy } from '../hooks/use-accounts'
import type { AccountRecord } from '../types'
export function AccountCreate({ staff = false }: { staff?: boolean }) {
 const navigate = useNavigate(); const policy = useProvisioningPolicy(); const key = useRef({ payload: '', key: crypto.randomUUID() })
 const form = useForm<Values>({ resolver: zodResolver(accountFormSchema), defaultValues: { name: '', email: '', borrower_type: staff ? 'FACULTY' : 'STUDENT', student_id: '', course: '', contact_number: '' } })
 const type = useWatch({ control: form.control, name: 'borrower_type' }); const base = staff ? '/admin/administration/accounts' : '/staff/borrowers'
 const command = useAccountCommand<Values, AccountRecord>(async value => { const input = staff ? { name: value.name, email: value.email } : { ...value, student_id: value.borrower_type === 'STUDENT' ? value.student_id : '' }; const payload = JSON.stringify(input); if (key.current.payload !== payload) key.current = { payload, key: crypto.randomUUID() }; return accountsApi.create(input, staff, key.current.key) })
 async function submit(value: Values) { if (command.isPending) return; try { const record = await command.mutateAsync(value); navigate(`${base}/${record.id}`, { replace: true, state: { created: true } }) } catch (error) { const fields = normalizeApiError(error).fields; for (const field of ['name', 'email', 'student_id', 'course', 'contact_number'] as const) { if (fields?.[field]) form.setError(field, { type: 'server', message: fields[field] }) } form.setFocus('email') } }
 const student = !staff && type === 'STUDENT'; const missingDomains = student && !policy.isPending && (!policy.data?.student_domains?.length || policy.isError)
 function field(name: Exclude<keyof Values, 'borrower_type'>, label: string, required = false, hint?: string) {
  return <FormField id={`account-${name}`} label={label} required={required} hint={hint} error={form.formState.errors[name]?.message}><Input type={name === 'email' ? 'email' : 'text'} autoComplete={name === 'email' ? 'email' : name === 'name' ? 'name' : 'off'} {...form.register(name)} /></FormField>
 }
 return <ManagementPage title={staff ? 'Add Staff Account' : 'Add Borrower'} description="The account holder creates their own eLabTrack password through secure activation." back={{to:base,label:`Back to ${staff ? 'accounts' : 'borrowers'}`}}><ManagementCard>
 {command.error && <Notice tone="danger">{apiErrorMessage(command.error)}</Notice>}{missingDomains && <Notice tone="warning">Student onboarding is unavailable until approved institutional email domains are configured. Faculty can be created individually.</Notice>}
 <form onSubmit={form.handleSubmit(submit)} noValidate aria-busy={command.isPending}>
 <FormSection title="Identity" disabled={command.isPending}>{!staff && <FormField id="borrower-type" label="Borrower type" required><NativeSelect {...form.register('borrower_type')}><option value="STUDENT">Student</option><option value="FACULTY">Faculty</option></NativeSelect></FormField>}{field('name','Full name',true)}</FormSection>
 <FormSection title="Contact information" disabled={command.isPending}>{field('email',student ? 'Institutional email' : 'Email',true,student && policy.data?.student_domains?.length ? `Configured domains: ${policy.data.student_domains.join(', ')}` : undefined)}{!staff && field('contact_number','Contact number (optional)')}</FormSection>
 {student && <FormSection title="Institutional details" disabled={command.isPending}>{field('student_id','Student ID',true,'Use the number printed on the official identification, including any leading zeros.')}{field('course','Course / program (optional)')}</FormSection>}
 <Notice title="Account access">No password is assigned by an administrator. {staff ? 'The account will have STAFF authority.' : 'Student and Faculty accounts have BORROWER authority and must review current officially published FSMO terms before borrowing.'}</Notice>
 <FormActions><AppButton type="submit" disabled={missingDomains || command.isPending}>{command.isPending ? 'Creating…' : 'Create Account'}</AppButton><AppButton variant="outline" nativeButton={false} render={<Link to={base} />}>Cancel</AppButton></FormActions></form>
 </ManagementCard></ManagementPage>
}
export function StaffCreate() { return <AccountCreate staff /> }
