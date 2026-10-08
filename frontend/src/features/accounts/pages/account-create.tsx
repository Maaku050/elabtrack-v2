import { useRef } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useForm, useWatch } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { accountFormSchema, type Values } from '../schemas/account.schema'
import { Input } from '@/components/ui/input'
import { OperationalShell } from '@/components/application/operational-shell'
import { AppButton, PageHeading, SurfaceCard } from '@/components/application/visual'
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
 return <OperationalShell active={staff ? 'Administration' : 'Borrowers'}><PageHeading title={staff ? 'Add Staff Account' : 'Add Borrower'} description="The account holder creates their own eLabTrack password through secure activation." /><SurfaceCard className="management-card management-form"><Link to={base}>← Back to {staff ? 'accounts' : 'borrowers'}</Link><h2>{staff ? 'Staff account' : 'Borrower information'}</h2>{command.error && <p className="auth-error" role="alert">{apiErrorMessage(command.error)}</p>}{missingDomains && <p role="alert" className="management-warning">Student onboarding is unavailable until approved institutional email domains are configured. Faculty can be created individually.</p>}
 <form onSubmit={form.handleSubmit(submit)} noValidate aria-busy={command.isPending}><fieldset disabled={command.isPending}>
 {!staff && <div className="auth-field"><label htmlFor="borrower-type">Borrower type</label><select id="borrower-type" {...form.register('borrower_type')}><option value="STUDENT">Student</option><option value="FACULTY">Faculty</option></select></div>}
 {(['name', 'email', ...(student ? ['student_id', 'course'] : []), ...(!staff ? ['contact_number'] : [])] as const).map(field => { const name = field as Exclude<keyof Values, 'borrower_type'>; const labels = { name: 'Full name', email: student ? 'Institutional email' : 'Email', student_id: 'Student ID', course: 'Course / program (optional)', contact_number: 'Contact number (optional)' }; const error = form.formState.errors[name]; return <div key={name} className="auth-field"><label htmlFor={`account-${name}`}>{labels[name]}</label><Input id={`account-${name}`} type={name === 'email' ? 'email' : 'text'} autoComplete={name === 'email' ? 'email' : name === 'name' ? 'name' : 'off'} aria-invalid={!!error} aria-describedby={error ? `error-${name}` : undefined} {...form.register(name)} />{error && <p className="field-error" id={`error-${name}`} role="alert">{error.message}</p>}</div> })}
 {student && policy.data?.student_domains?.length && <p className="management-note">Configured domains: {policy.data.student_domains.join(', ')}</p>}
 <p className="management-note">No password is assigned by an administrator. {staff ? 'The account will have STAFF authority.' : 'Student and Faculty accounts have BORROWER authority and must review current officially published FSMO terms before borrowing.'}</p><div className="management-actions"><AppButton type="submit" disabled={missingDomains || command.isPending}>{command.isPending ? 'Creating…' : 'Create Account'}</AppButton><AppButton  variant="outline" nativeButton={false} render={<Link to={base} />}> Cancel</AppButton></div></fieldset></form>
 </SurfaceCard></OperationalShell>
}
export function StaffCreate() { return <AccountCreate staff /> }
