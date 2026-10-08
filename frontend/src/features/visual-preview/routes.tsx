import { useLocation } from 'react-router-dom'
import { BorrowerHomePreview, EquipmentCatalogPreview } from './borrower-pages'
import { StaffDashboardPreview, PendingRequestsPreview } from './staff-pages'

export function Component() {
  const { pathname } = useLocation()
  if (pathname.endsWith('/borrower/equipment')) return <EquipmentCatalogPreview />
  if (pathname.endsWith('/borrower/home')) return <BorrowerHomePreview />
  if (pathname.endsWith('/staff/requests')) return <PendingRequestsPreview />
  return <StaffDashboardPreview audience={pathname.includes('/admin/') ? 'admin' : 'staff'} />
}
