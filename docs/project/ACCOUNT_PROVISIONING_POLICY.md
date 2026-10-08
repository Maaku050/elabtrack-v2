# Student/Faculty account provisioning policy

2026-10-09 current implementation overlay. **Batch 1 explicitly authorizes Phase 5 and conditional independent Phase 6; implementation and verification are tracked in [the batch plan](BATCH1_IMPLEMENTATION_PLAN.md) and phase reports.** The owner's final Student/Faculty update, [DEC-070](DECISIONS.md#elab-v2-dec-070--final-studentfaculty-provisioning-policy), supersedes the Staff creation grant in DEC-053/069 and the universal institutional-email/optional-ID assumptions in DEC-068. Subsequent owner decisions [DEC-071–073](DECISIONS.md#elab-v2-dec-071--brevo-transactional-email-provider) select Brevo, defer official terms until presentation without blocking independent development, and allow warned/confirmed Student deactivation despite obligations. All approved Student/Faculty provisioning rules below remain unchanged; other borrowing, return, replacement, fine and session rules remain unchanged. Historical phase reports and approved visual assets retain their recorded evidence.

## APPROVED PRODUCT RULES

### Authority and individual creation

Only a current active **ADMIN** can create any new account. STAFF may assist Student/Faculty users operationally but cannot create accounts. Backend enforcement is mandatory on individual creation, bulk preparation, confirmation and replay; hiding a button is insufficient. Privileged STAFF/ADMIN accounts belong to a separate restricted Admin workflow, never the borrower form or Student spreadsheet. No public signup or shared administrative credential.

| Requirement | STUDENT | FACULTY |
|---|---|---|
| Authorization role | BORROWER | BORROWER |
| Borrower type attribute | STUDENT | FACULTY |
| Creation | Admin individually or Admin Student Excel import | Admin individually only |
| Email | Valid, unique institutional address from approved SKSU domain configuration | Valid, unique, accessible address; institutional domain is not required |
| Student ID | Required, unique string matching official student identification | Not required |
| Password | Borrower chooses a separate eLabTrack password through secure activation | Same |
| First-use terms | Accept current officially published FSMO terms via Phase 4B | Same |

Student IDs are identifiers, not numeric quantities: `28366` is an illustrative string, not a seeded user or numeric constraint. Preserve leading zeroes; never cast IDs to integers. Retain the internal UUID as the historical identity key and normalized unique email under existing application/database conventions. Reject Student-ID/email conflicts rather than merging by name or whichever identifier is convenient. Email changes never replace the UUID or rewrite history. Exact ID formatting/reconciliation remains a technical dependency; requiredness, uniqueness and string representation are approved rules.

### Login, secure activation and terms

Both categories use **email + a separate eLabTrack password**. No institutional SSO, mailbox-password collection or Admin-assigned borrower password. Preserve Phase 4A authentication. Borrowers choose passwords through secure activation; single-use email activation links remain the recommended delivery mechanism. Student links go to the approved institutional address; Faculty links go to their valid accessible address, including non-institutional email. Syntax, domain membership, roster presence or Admin creation alone is not proof of mailbox ownership.

Activation requires unpredictable purpose/account-bound tokens, hash-only persistence, expiry, atomic single use, current-account checks, rate limits and safe invalidation/reissue. It must not reactivate disabled accounts or grant roles. Never expose raw passwords/tokens in spreadsheets, exports, logs, audit or results. **Brevo is selected for account activation and future password recovery (DEC-071).** Integrate only through a secure backend adapter when account-management implementation is separately authorized. Keep its API key out of browser code, commits and logs. Live delivery requires a configured API key, verified sender and successful delivery testing. Vendor selection is resolved; integration/delivery, exact lifecycle/ownership/recovery and existing-account handling remain unimplemented dependencies. A test adapter or send intent does not prove real delivery.

Both Student and Faculty must accept current **officially published** FSMO terms using existing Phase 4B versioning/acceptance. Official V2 terms will be finalized after presenting the application to FSMO (DEC-072); no official text has been approved. This external approval does **not block independent account-management or inventory development** when separately authorized. It gates official publication/consent and live borrowing, which must still enforce published terms and documented acceptance. Never publish placeholder/reference terms as official or record acceptance of unapproved content. Missing official terms cannot authorize new borrowing; account access, existing obligations and logout retain recorded Phase 4B behavior.

### Student-only Excel creation

Admin-controlled standard Excel template: **STUDENTS ONLY**. Recommended columns: `studentId`, `name`, `email`, `course`, `contactNumber`. Student ID, name and email are required identity inputs; course/contact requiredness and actual roster availability remain to be finalized without weakening Student ID or institutional email requirements. Read studentId as text and preserve leading zeroes through spreadsheet handling; reject lossy/ambiguous identities rather than inventing digits.

Exclude password, admin role, staff role and faculty role columns. The standard template does not assign authorization or borrower category: the server always creates `role = BORROWER`, `borrower_type = STUDENT`. Faculty accounts cannot be created through it. Reject supplied credentials, privilege/category assignment and conflicting existing Faculty/Staff/Admin identities; never silently coerce Faculty into Student.

Upload → validate complete file → show complete validation preview → Admin reviews/selects eligible rows → explicit confirmation → server reauthorization/revalidation → durable per-record results. File selection/validation alone never creates accounts. Preview covers required fields, approved SKSU domain, Student-ID/email duplicates within the file/database, ambiguous/conflicting identities, invalid rows and planned results. No silent skipping, automatic overwrite or reactivation. Safe retries and durable audit bind results to confirmed rows. Implemented limits, formula rejection, exact matching and atomic selected-row transaction semantics are in [API_CONTRACTS](../API_CONTRACTS.md#phase-5-implemented-account-management-contracts).

### Separate Student-only bulk deactivation

**Admin-only and Student-only.** Reusing the Student import roster is allowed for this separate workflow. Match stable Student ID and normalized email to an existing internal UUID. Reject unmatched, ambiguous or conflicting identities; never match by name alone or prefer one conflicting identifier. No missing-ID fallback may bypass the required Student identity rule.

Show matched Students, unmatched records, conflicting identifiers, already-inactive accounts and authoritative outstanding borrowing/fine/replacement obligations before explicit selection/confirmation. Recheck current target `role = BORROWER` **and** `borrower_type = STUDENT` at confirmation/execution: role alone does not distinguish Faculty. Never deactivate Faculty, Staff or Admin through this workflow, including category/role changes after preview. Leave unselected/roster-absent accounts unchanged. Never hard-delete accounts/history, automatically reactivate, clear fines or erase obligations.

**APPROVED Student deactivation rule (DEC-073; OPEN-029 RESOLVED):** Admin may deactivate a Student despite outstanding fines, active/overdue borrowing, unreturned equipment or replacement obligations. Applies individually and in the Student-only bulk workflow for graduation, withdrawal, transfer or other authorized deactivation. Display appropriate obligation warnings and require confirmation; positive outstanding obligations must never block deactivation or require an obligation-based override. Unavailable projections must not be represented as zero; truthful failure/warning behavior is a technical contract dependency, not an unresolved obligation policy. Deactivation never clears/waives fines, records payments, marks equipment returned, resolves replacements, closes unresolved borrowing, changes overdue calculations or deletes history. Existing due times/stock/obligation balances/fine evidence remain unchanged; overdue continues under approved rules. FSMO resolves matters separately through authorized Staff/Admin return/replacement commands and Admin-only auditable fine clearance.

### Individual account-management UI handoff

The authorized implementation adapts the approved layout. Admin's borrower form offers **STUDENT / FACULTY**, with role fixed to BORROWER. Student requires Student ID and approved institutional email; Faculty requires a valid unique accessible email without institutional-domain restriction or Student-ID requirement. Admin never supplies borrower passwords. Keep STAFF/ADMIN creation in the separate restricted administrative workflow. Staff retains authorized operational directory/detail assistance without creation actions; direct API calls must be denied server-side.

The [visual map](../ux/APPROVED_VISUAL_IMPLEMENTATION_MAP.md), [flows](../ux/STAFF_ADMIN_FLOWS.md), [architecture](../ux/UX_ARCHITECTURE.md) and [wireframes](../ux/WIREFRAMES.md) describe required future behavior/field adaptations. Approved PNG/SVG mockups and manifest are unchanged. Any older pictured Staff creation action or generic bulk category cannot override DEC-070. No UI is implemented here.

## REMAINING TECHNICAL DEPENDENCIES

**Current implementation evidence:** users have stable UUID, email/name, bcrypt password hash, BORROWER/STAFF/ADMIN role, is_active and timestamps. ParseEmail trims/parses/lowercases; PostgreSQL CITEXT enforces unique email. Production public registration is absent. There is no product provisioning/import/deactivation API, Student-ID/borrower-type storage, domain allowlist, activation token/state/provider, password-change or recovery flow. Internal creation remains unmounted. Phase 4B implements terms publication/acceptance, with no official terms published. Admin-only creation and category-specific validation above are future backend contracts, not implemented permissions.

- OPEN-001/009: activation lifecycle, token expiry/reissue, mailbox ownership evidence, existing-account treatment, email changes, password change/recovery and revocation.
- OPEN-028: incorporate approved SKSU Student domains; preserve exact textual IDs; define normalization, changed-email reconciliation, spreadsheet encoding/cell handling and course/contact optionality. No Faculty allowlist or optional Student-ID decision remains.
- OPEN-017: implement the selected Brevo backend adapter, configure its API key/verified sender/trusted link origin and verify actual delivery/failure/retry handling for activation and future recovery. Vendor choice is resolved; test adapters never prove live delivery.
- Reviewed Phase 5 API/error contracts, bounded file/row/query limits, preview expiry/content binding, batch atomicity, concurrent duplicate protection and durable audit/idempotency. Domain/application owns use cases/ports; infrastructure implements parser/mail/SQL adapters; HTTP adapts requests. No new dependencies or schema are authorized here.

Engineering recommendation: bind confirmation to reviewed operation/rows/payload; reauthorize Admin on execution/replay; lock actor/targets in stable order and recheck role/category/identity/status/obligations. Persist effects/results/audit/receipt atomically under the selected transaction strategy. Same-command retries return original committed results; changed selections cannot silently reuse a receipt. Audit actor, operation, target UUIDs, time and outcome/reason without credentials or unnecessary roster contents.

Future acceptance tests: Staff/Borrower denial on every creation path; Student/Faculty conditional validation; leading-zero/duplicate Student IDs; non-institutional Faculty email/activation; forbidden Faculty/privileged/password spreadsheet inputs; complete preview/no pre-confirmation mutation; concurrent identity conflicts; target role/category changes after preview; Student-only deactivation; obligation warnings/confirmation without an obligation veto; unchanged borrowing/stock/due/overdue/fine/replacement/history evidence; retained history; safe replay/audit; secure activation; both categories' current official terms acceptance/missing-content behavior. These are planning criteria, not executed provisioning tests.

## REMAINING INSTITUTIONAL APPROVALS

| Approval/input | Gate |
|---|---|
| Exact approved SKSU Student domain configuration and official ID/roster format examples, without personal records | OPEN-028; no invented domain/ID regex or Faculty domain requirement |
| Operational activation/ownership/recovery responsibility | OPEN-001/009; borrower-owned password and Brevo selection approved; safe operational procedures still need review |
| Official FSMO terms text/version/publication responsibility after application presentation | OPEN-016; gates official publication/consent and live borrowing, not independent account-management/inventory development |

Brevo provider selection and warned/confirmed Student deactivation regardless of obligations are resolved in DEC-071/073; OPEN-029 is no longer an active gate. Admin-only creation, Student-only bulk workflows, required unique textual Student IDs, Faculty individual creation/accessible email/no Student-ID requirement, BORROWER role and separate borrower-owned passwords are **resolved product rules**, not approval gates. This update does not start Phase 5, change Phase 4A authentication, create/alter migrations, change approved mockups, commit, push or deploy.

## Phase 5 readiness

**Product policy is ready for a scoped Phase5 implementation plan; Phase5 is NOT STARTED and still requires explicit implementation authorization.** No provider-selection or outstanding-obligation deactivation policy remains open. Independent account-management development and, when separately authorized under the roadmap, inventory development may proceed without finalized official terms. Remaining technical work covers storage/conditional validation, import/matching/preview/transactions/audit, secure activation/recovery and the backend Brevo adapter. External inputs are actual SKSU domain/official roster details, deployment-only Brevo credentials/verified sender/live delivery tests and FSMO official terms approval after presentation. Official terms must gate live borrowing/consent, never be bypassed or fabricated for development. Operational recovery responsibilities and later notification cadence retain their genuine scope.

## Batch 1 technical handoff

[Implemented Phase5 contracts](../API_CONTRACTS.md#phase-5-implemented-account-management-contracts) define hash-only single-use activation, bounded Excel previews, deterministic locks, idempotency and account audits. Technical defaults are not final institutional activation/recovery procedure: TTL24h, resend1m, preview30m; backend configuration validates bounds and trusted link origin. Existing unclassified Borrowers retain UUID/email/password/history and are excluded from Student rosters. Individual identity/type/role changes and new Admin creation/promotion are unavailable. Student onboarding fails closed without configured domains; exact existing Student ID/email matches can still be deactivated when domain configuration changes. Bulk creation is atomic for selected valid rows and leaves accounts pending; send activation individually from details. Preview/audit retention operations await institutional retention policy.

LIVE EMAIL DELIVERY NOT VERIFIED. Approved Student domain input, verified Brevo sender/key/recipient testing, institutional activation/recovery ownership and official terms after presentation remain external dependencies. Neither accounts nor inventory development publishes or accepts unapproved terms. Future borrowing remains guarded by official publication and durable current-version acceptance.
