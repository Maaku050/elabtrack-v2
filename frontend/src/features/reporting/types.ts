export interface ReportDefinition { key: string; label: string; columns: string[]; equipment_filter: boolean; category_filter: boolean; borrower_filter: boolean; admin_only: boolean }
export interface ReportFilter { due_today?:boolean; page: number; per_page: number; search?: string; status?: string; from?: string; to?: string; equipment_id?: string; category_id?: string; borrower_id?: string }
export interface ReportPage { kind: string; columns: string[]; rows: string[][]; total: number; page: number; per_page: number; as_of: string }
export interface Dashboard { trends?: {day:string;issued:number;pending:number;denied:number}[]; metrics: Record<string, number>; recent: { id: string; borrowing_id: string; reference: string; kind: string; at: string }[]; as_of: string }
