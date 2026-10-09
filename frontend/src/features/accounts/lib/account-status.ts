export function activationLabel(required: boolean) { return required ? 'Pending activation' : 'Activation not required' }
export function deliveryLabel(status: string) {
 switch (status) {
 case 'UNCONFIGURED': return 'Activation email was not sent. The email service is not configured.'
 case 'PENDING': return 'Activation email has not been submitted yet.'
 case 'ACCEPTED': return 'The email provider accepted the activation message. Recipient delivery is not confirmed.'
 case 'FAILED': return 'Activation email submission failed. An Admin can retry sending the link.'
 case 'UNKNOWN': return 'Email submission outcome is unknown. Recipient delivery is not confirmed.'
 case 'NOT_REQUIRED': return 'No activation email is required.'
 default: return 'Email submission status is unavailable.'
 }
}
