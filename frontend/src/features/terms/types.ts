export interface TermsVersion {
 id: string
 version: string
 title: string
 body: string
 content_hash: string
 published_at: string
}
export interface TermsAcceptance { id: string; terms_version_id: string; accepted_at: string }
export interface TermsStatus {
 state: 'unpublished' | 'required' | 'updated' | 'accepted'
 current_terms: TermsVersion | null
 acceptance: TermsAcceptance | null
 has_previous_acceptance: boolean
 acceptance_required: boolean
 can_initiate_borrowing: boolean
}
