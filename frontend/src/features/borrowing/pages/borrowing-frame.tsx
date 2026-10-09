import type { ReactNode } from 'react'
import { BorrowerShell } from '@/components/application/shells'
import { borrowerDestinations } from '@/components/application/borrower-navigation'
import { ManagementPage } from '@/components/application/management'
import { ThemeControl } from '@/components/application/visual'
import { useAuthStore } from '@/stores/auth-store'
export function BorrowingFrame({ title, description, actions, back, children }: { title: string; description: string; actions?: ReactNode; back?: { to: string; label: string }; children: ReactNode }) {
 const borrower = useAuthStore(s => s.user?.role === 'BORROWER')
 const page = <ManagementPage title={title} description={description} domain="Borrowing & reservations" actions={<>{actions}{borrower && <ThemeControl />}</>} back={back}>{children}</ManagementPage>
 return borrower ? <BorrowerShell active="My Borrowings" destinations={borrowerDestinations}>{page}</BorrowerShell> : page
}
