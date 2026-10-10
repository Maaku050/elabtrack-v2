export type BorrowingStatus = 'PENDING' | 'CHECKED_OUT' | 'COMPLETED' | 'DENIED' | 'CANCELLED' | 'EXPIRED'
export interface BorrowingItem { id: string; equipment_id: string; name: string; quantity: number; reserved_quantity: number; issued_quantity: number; good_quantity?: number; damaged_quantity?: number; lost_quantity?: number; outstanding_quantity?: number }
export interface BorrowingEvent { id: string; actor_id: string | null; kind: string; reason: string; occurred_at: string }
export interface Borrowing { id: string; reference: string; borrower_id: string; borrower_name: string; borrower_type: string; student_id: string; status: BorrowingStatus; entry_path: 'REQUEST' | 'DIRECT'; acceptance_id: string; created_at: string; expires_at: string | null; checked_out_at: string | null; due_at: string | null; terminal_at: string | null; denial_reason: string; items: BorrowingItem[]; event_count?: number; paged_history?: boolean; replacement_units_outstanding?: number; events: BorrowingEvent[]; is_overdue: boolean; completed_at?: string | null; fine_final_minor?: number | null; fine?: Fine; returns?: ReturnDisposition[]; obligations?: Obligation[]; replacements?: ReplacementAcceptance[]; clearances?: FineClearance[] }
export interface BorrowingFilter { page: number; per_page: number; status?: string; search?: string; borrower_id?: string }
export interface BorrowingPage { items: Borrowing[]; page: number; per_page: number; total: number }
export interface RequestInput { items: { equipment_id: string; quantity: number }[]; confirm: boolean; borrower_id?: string; due_at?: string; physical_handover_confirmed?: boolean }
export interface DecisionInput { confirm: boolean; reason?: string; due_at?: string; physical_handover_confirmed?: boolean }
export interface ReturnLine { item_id: string; good: number; damaged: number; lost: number }
export interface ReturnInput { lines: ReturnLine[]; reason: string; confirm: boolean; expected_event_count: number }
export interface ReplacementInput { obligation_id: string; quantity: number; reason: string; equivalent_confirmed: boolean; confirm: boolean; expected_event_count: number }
export interface FineInput { method: 'PAID' | 'WAIVED' | 'OTHER_RESOLUTION'; note: string; expected_outstanding_minor: number; confirm: boolean }
export interface Fine { assessed_minor: number; cleared_minor: number; outstanding_minor: number; is_final: boolean; as_of: string }
export interface Obligation { id: string; item_id: string; equipment_id: string; kind: 'DAMAGED' | 'LOST'; required: number; accepted: number }
export interface ReturnDisposition extends ReturnLine { id: string; actor_id: string; reason: string; occurred_at: string }
export interface ReplacementAcceptance { id: string; obligation_id: string; equipment_id: string; equipment_name: string; equivalent_confirmed: boolean; quantity: number; actor_id: string; reason: string; occurred_at: string }
export interface FineClearance { id: string; actor_id: string; assessed_minor: number; cleared_minor: number; method: string; note: string; occurred_at: string }

export interface HistoryPage<T=Record<string,unknown>> { items:T[]; kind:string; page:number; per_page:number; total:number }
