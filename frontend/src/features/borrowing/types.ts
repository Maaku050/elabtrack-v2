export type BorrowingStatus = 'PENDING' | 'CHECKED_OUT' | 'DENIED' | 'CANCELLED' | 'EXPIRED'
export interface BorrowingItem { id: string; equipment_id: string; name: string; quantity: number; reserved_quantity: number; issued_quantity: number }
export interface BorrowingEvent { id: string; actor_id: string | null; kind: string; reason: string; occurred_at: string }
export interface Borrowing { id: string; reference: string; borrower_id: string; borrower_name: string; borrower_type: string; student_id: string; status: BorrowingStatus; entry_path: 'REQUEST' | 'DIRECT'; acceptance_id: string; created_at: string; expires_at: string | null; checked_out_at: string | null; due_at: string | null; terminal_at: string | null; denial_reason: string; items: BorrowingItem[]; events: BorrowingEvent[]; is_overdue: boolean }
export interface BorrowingFilter { page: number; per_page: number; status?: string; search?: string; borrower_id?: string }
export interface BorrowingPage { items: Borrowing[]; page: number; per_page: number; total: number }
export interface RequestInput { items: { equipment_id: string; quantity: number }[]; confirm: boolean; borrower_id?: string; due_at?: string; physical_handover_confirmed?: boolean }
export interface DecisionInput { confirm: boolean; reason?: string; due_at?: string; physical_handover_confirmed?: boolean }
