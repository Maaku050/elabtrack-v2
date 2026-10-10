export interface Notification { id: string; borrowing_id: string; kind: string; title: string; body: string; scope: 'BORROWER' | 'OPERATIONS'; created_at: string; read_at: string | null }
export interface NotificationFilter { page: number; per_page: number; read?: '' | 'read' | 'unread' }
export interface NotificationPage { items: Notification[]; page: number; per_page: number; total: number; unread: number }
