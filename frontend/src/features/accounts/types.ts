export interface Obligations { availability: 'UNAVAILABLE' | 'AVAILABLE'; active_borrowings?: number; overdue_borrowings?: number; unreturned_units?: number; fine_minor?: number; replacement_units?: number }
export interface AccountRecord { id: string; name: string; email: string; role: 'BORROWER' | 'STAFF' | 'ADMIN'; is_active: boolean; borrower_type: '' | 'STUDENT' | 'FACULTY'; student_id: string; course: string; contact_number: string; activation_required: boolean; delivery_status: string; created_at: string; updated_at: string; obligations: Obligations }
export interface AccountInput { name: string; email: string; borrower_type?: 'STUDENT' | 'FACULTY'; student_id?: string; course?: string; contact_number?: string }
export interface AccountFilter { page: number; per_page: number; search: string; borrower_type: string; status: string }
export interface AccountPage { items: AccountRecord[]; total: number; page: number; per_page: number }
export interface ProvisioningPolicy { student_domains: string[] | null; activation_ttl_seconds: number; resend_cooldown_seconds: number }
export interface RosterRow { number: number; input: AccountInput; error?: string; target_id?: string; outcome?: string; account_id?: string; obligations: Obligations }
export interface RosterBatch { id: string; operation: 'CREATE' | 'DEACTIVATE'; rows: RosterRow[]; expires_at: string; confirmed: boolean }
export interface AuditEvent { id: string; actor_id: string; account_id: string; action: string; occurred_at: string }
