export interface Stock { available: number; reserved: number; checked_out: number; damaged_held: number; total_tracked: number }
export interface Equipment { id: string; name: string; description: string; category_id: string | null; category_name: string; status: 'ACTIVE' | 'INACTIVE' | 'ARCHIVED'; stock: Stock; stock_sequence: number; metadata_version: number; image_id: string | null; created_at: string; updated_at: string; archive_safety?: string }
export interface Category { id: string; name: string; is_active: boolean; version: number }
export interface EquipmentFilter { page: number; per_page: number; search?: string; status?: string; category_id?: string; available_only?: boolean; sort?: string }
export interface EquipmentPage { items: Equipment[]; total: number; page: number; per_page: number; totals: Stock }
export interface Metadata { name: string; description: string; category_id: string | null; expected_version: number }
export interface Movement { id: string; actor_id: string | null; sequence: number; kind: string; delta: Stock; after: Stock; reason: string; created_at: string }
export interface Adjustment { kind: 'ADD' | 'REMOVE' | 'RECONCILE'; quantity: number; reason: string; expected_sequence?: number; confirm: boolean }
