import { expect, it } from 'vitest'
import { activationLabel, deliveryLabel } from './account-status'
it('keeps activation independent from administrative access and provider acceptance', () => {
 expect(activationLabel(true)).toBe('Pending activation')
 expect(activationLabel(false)).toBe('Activation not required')
 expect(deliveryLabel('UNCONFIGURED')).toContain('was not sent')
 expect(deliveryLabel('PENDING')).toContain('not been submitted')
 expect(deliveryLabel('ACCEPTED')).toContain('delivery is not confirmed')
 expect(deliveryLabel('FAILED')).toContain('failed')
 expect(deliveryLabel('UNKNOWN')).toContain('unknown')
 expect(deliveryLabel('NOT_REQUIRED')).toContain('No activation email')
 expect(deliveryLabel('unexpected')).toContain('unavailable')
})
