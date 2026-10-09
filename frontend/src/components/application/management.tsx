import { cloneElement, type ReactElement, type ReactNode, useId } from 'react'
import { Link } from 'react-router-dom'
import { ArrowLeft, CircleAlert, CircleCheck, Info } from 'lucide-react'
import { Card } from '@/components/ui/card'
import { Field, FieldLabel, FieldDescription, FieldError, FieldSet, FieldLegend, FieldGroup } from '@/components/ui/field'
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert'
import { DirectoryHeader } from './directory'

export function ManagementPage({ title, description, domain = 'Account management', actions, back, children }: { title: string; description: string; domain?: string; actions?: ReactNode; back?: { to: string; label: string }; children: ReactNode }) {
  return <div className="directory-page management-page"><DirectoryHeader title={title} description={description} eyebrow={domain} actions={actions} />{back && <Link className="management-back" to={back.to}><ArrowLeft aria-hidden="true" />{back.label}</Link>}{children}</div>
}

export function ManagementCard({ title, description, actions, children, className = '' }: { title?: string; description?: string; actions?: ReactNode; children: ReactNode; className?: string }) {
  return <Card className={`management-surface ${className}`}>{title && <div className="management-section-heading"><div><h2>{title}</h2>{description && <p>{description}</p>}</div>{actions}</div>}{children}</Card>
}

export function FormSection({ title, description, children, disabled }: { title: string; description?: string; children: ReactNode; disabled?: boolean }) {
  return <FieldSet className="form-section" disabled={disabled}><FieldLegend>{title}</FieldLegend>{description && <FieldDescription>{description}</FieldDescription>}<FieldGroup className="form-grid">{children}</FieldGroup></FieldSet>
}

/** Associate RHF controls with existing Field labels and inline feedback. */
export function FormField({ id, label, required, hint, error, wide = false, children }: { id: string; label: string; required?: boolean; hint?: string; error?: string; wide?: boolean; children: ReactElement<{ id?: string; 'aria-describedby'?: string; 'aria-invalid'?: boolean; 'aria-required'?: boolean }> }) {
  const described = [children.props['aria-describedby'], hint && `${id}-hint`, error && `${id}-error`].filter(Boolean).join(' ') || undefined
  return <Field className={`form-field ${wide ? 'form-field-wide' : ''}`} data-invalid={!!error}><FieldLabel htmlFor={id}>{label}{required && <span className="required-marker" aria-hidden="true">*</span>}</FieldLabel>{cloneElement(children, { id, 'aria-describedby': described, 'aria-invalid': !!error, 'aria-required': required || undefined })}{hint && <FieldDescription id={`${id}-hint`}>{hint}</FieldDescription>}{error && <FieldError id={`${id}-error`}>{error}</FieldError>}</Field>
}

export function FormActions({ children }: { children: ReactNode }) { return <div className="form-actions">{children}</div> }

export function Notice({ tone = 'info', title, children }: { tone?: 'info' | 'success' | 'warning' | 'danger'; title?: string; children: ReactNode }) {
  const Icon = tone === 'success' ? CircleCheck : tone === 'warning' || tone === 'danger' ? CircleAlert : Info
  return <Alert className="management-feedback" data-tone={tone} role={tone === 'danger' ? 'alert' : 'status'}><Icon aria-hidden="true" />{title && <AlertTitle>{title}</AlertTitle>}<AlertDescription>{children}</AlertDescription></Alert>
}

export function WorkflowSteps({ steps, current }: { steps: string[]; current: number }) {
  const id = useId()
  return <ol className="management-steps" aria-label="Workflow progress">{steps.map((step, index) => <li key={`${id}-${step}`} aria-current={index === current ? 'step' : undefined} data-complete={index < current}><span aria-hidden="true">{index < current ? <CircleCheck /> : index + 1}</span>{step}</li>)}</ol>
}
