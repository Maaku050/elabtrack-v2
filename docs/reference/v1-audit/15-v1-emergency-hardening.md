# 15. V1 emergency user-management hardening

## Scope and status

This patch contains the critical unauthenticated V1 user-management HTTP API. It is not a V1 redesign or V2 implementation. It changes no Firestore schema, borrowing behavior, inventory logic, or institutional role policy, and it has not been deployed.

## Pre-change behavior recorded

### API contract

The exported `userManagement` Express application supports these operations:

| Route | Existing request shape | Existing success shape |
|---|---|---|
| `POST /createUser` | JSON user fields: `email`, `password`, `name`, `role`, `course`, `contactNumber`, optional `imageBase64` | HTTP 201; `{ status, message, data }` with created profile fields |
| `POST /createBulkUsers` | `{ users: CreateUserRequest[] }` | HTTP 200; `{ status, message, data }` with total/success/failure counts and result lists |
| `POST /deleteBulkUsers` | `{ emails: string[] }` | HTTP 200; `{ status, message, data }` with totals and per-email results |
| `GET /getUser/:uid` | UID path parameter | HTTP 200; `{ status, data }` containing the Firestore profile |
| `PATCH /updateUser/:uid` | Arbitrary profile updates plus optional `imageBase64`; `uid` and `createdAt` are removed | HTTP 200; `{ status, message, data: { uid, updates } }` |
| `DELETE /deleteUser/:uid` | UID path parameter | HTTP 200; `{ status, message, data: { uid } }` |

Validation and operation-specific error responses remain unchanged. The security boundary adds JSON `401` and `403` errors using the existing `{ status: "error", message }` convention.

Before this patch, `cors({ origin: true })` reflected any requesting origin. Every route reached Firebase Admin Auth, Firestore, and Storage without verifying the caller. The single-user, bulk-create, single-delete, and bulk-delete modals call hardcoded HTTPS Function routes. No current frontend caller of the get or update HTTP routes was found; profile editing still writes to Firestore directly.

V1 signs users in through Firebase Authentication. Navigation and client guards obtain role/status from the Firestore `users` profile. Account creation also writes a role custom claim, but client authorization does not use it and claims may be stale. The known roles are `student`, `staff`, and `admin`; staff and admin are treated equivalently by the current admin UI.

## Containment changes

### Authentication

All six routes now pass through one middleware before JSON body parsing or an operation handler. It requires exactly one `Authorization: Bearer <Firebase ID token>` header and calls Firebase Admin Auth `verifyIdToken`. Missing/malformed headers and invalid/expired tokens return HTTP 401. Request-body UID, role, and email values are never used to authenticate or authorize the caller.

### Authorization

After token verification, the middleware reads the caller's current `users/{verifiedUid}` Firestore document. Only an `active` profile whose stored role is `staff` or `admin` reaches a handler. Missing profiles, inactive profiles, students, and unknown roles return HTTP 403. This uses the same current profile source as V1 navigation while making the server the enforcement point.

### Client calls

The four active frontend callers now use a shared authenticated-fetch helper. It obtains the current Firebase user's ID token with `getIdToken()` for each request and sends it in the standard Authorization header. It does not place the token in localStorage or any new persistence layer. Existing methods, URLs, bodies, and success/error parsing remain unchanged.

### CORS

Repository deployment evidence identifies Firebase Hosting project `equipment-tracking-syste-65e94`. The allowlist is limited to its canonical `web.app` and `firebaseapp.com` HTTPS origins, plus HTTP `localhost` and `127.0.0.1` origins for the repository's active Expo/emulator development workflows. Origin-less native/non-browser calls remain supported; they still require authentication and authorization. No custom production domain is evidenced in the repository. If one exists, operators must verify and explicitly add that exact origin before deployment.

### Defensive logging

The authorization boundary logs operation method/path, verified caller UID when available, authorization outcome, completion category, HTTP status, and platform timestamp. It never logs the Authorization header or token. The prior whole-request-body log—which included plaintext passwords—was removed. Several unnecessary name/email/image URL/path messages were also reduced. Firebase operation errors continue to be logged for diagnosis and should be reviewed under the production log-access/retention policy.

## Staff/admin privilege behavior retained

No institutional rule distinguishes staff from admin today. Both can still create student, staff, or admin accounts, bulk-create those roles, update profile role fields through the HTTP route, and delete users including admins. This preserves observed V1 behavior and avoids inventing governance policy.

This remains a **high-risk policy ambiguity**. In particular, a compromised staff account can create another privileged account or delete an administrator. Because staff and admin already receive the same V1 client capabilities, changing that boundary requires an accountable school decision, a target-operation matrix, and recovery/bootstrap rules. The emergency patch prevents unauthenticated and student access but does not resolve privileged separation of duties.

## Focused tests

The authorization middleware has dependency-injected Firebase Admin token verification and profile lookup so tests make no live Firebase calls. Tests cover:

- missing Authorization header → 401;
- malformed Bearer header → 401;
- mocked invalid/expired token → 401;
- authenticated active student → 403;
- authenticated active staff → operation handler reached;
- authenticated active admin → operation handler reached.

An inactive staff/admin profile is also denied by the implemented authorization rule, although the minimum requested matrix focuses on role outcomes.

## Remaining V1 risks and required discovery

### Critical/high risks not fixed here

- **Firestore authorization remains unverifiable.** No rules source is tracked. Live rules must be exported read-only and reviewed for authentication, staff/admin capability, student ownership, allowed fields, query constraints, and denial of direct privilege changes.
- **Storage authorization remains unverifiable.** Live rules and object ACLs must be inspected for profile/equipment upload, read, replace, and delete access. Function-created profile images remain publicly readable under existing behavior.
- **Inventory concurrency remains high risk.** Client-side, non-transactional stock arithmetic can oversubscribe equipment or lose updates. This patch intentionally does not alter borrowing/inventory behavior.
- **Direct client writes remain broad and rules-dependent.** User profiles, equipment, transactions, records, and fine-like state are still modified outside this HTTP boundary.
- **Staff/admin separation is unresolved.** Equivalent privileges include creating and deleting admins, changing roles, managing equipment/transactions, viewing PII, and clearing fines.
- **No rate limiting or privileged step-up authentication exists.** A stolen privileged session can repeatedly invoke destructive operations until token/account revocation takes effect.
- **Multi-resource operations are non-atomic.** Auth, Firestore, and Storage create/delete failures can leave orphaned or partial state.

### Password spreadsheet risk

Bulk account creation still accepts initial plaintext passwords from an Excel file and sends them inside the authenticated TLS request. This compatibility behavior was not redesigned. Passwords may remain in source spreadsheets, downloads, browser/device caches, backups, screen recordings, or operator sharing channels. The removed body log prevents deliberate server logging, but an authorized operating procedure is still required for file creation, transfer, local deletion, password handoff, forced change/reset, and incident response.

### Live read-only inspection required

Do not invent or deploy replacement rules from repository assumptions. An authorized Firebase custodian must provide timestamped read-only evidence for:

- deployed Firestore rules, rule versions/tests, composite indexes, exemptions, and TTL settings;
- deployed Storage rules, buckets, public ACLs, object paths, orphan objects, and referenced/missing files;
- deployed Functions, regions, runtimes, IAM/invoker policy, active routes, scheduler jobs, logs, and environment configuration names;
- installed Extensions and the notification delivery consumer/fields;
- Authentication providers, disabled users, email-verification state, custom-claim variants, and claim/profile mismatches;
- Hosting releases, exact production/custom origins, redirects/headers, and environment separation;
- live data volumes and integrity exceptions needed for inventory, user, transaction, record, fine, and notification risk assessment.

## Validation record

- Functions `npm run build`: passed.
- Functions `npm run lint`: passed with warnings only; no lint errors. Existing unused-catch warnings remain in the legacy function source.
- Focused authorization tests: see the final task validation output; no production Firebase access is used.
- Root `npx tsc --noEmit`: blocked by pre-existing errors in admin return typing, report chart props, bottom-sheet refs, and table refs. No error identified the new authenticated-fetch helper or changed callers.
- No deployment, commit, or push was performed.

## Deployment steps required

1. Confirm the actual production web origin/custom domain and amend the exact allowlist only if repository evidence is incomplete.
2. Confirm every intended staff/admin has one current `users/{uid}` profile with lowercase `staff` or `admin` and `status: "active"`; mismatches will correctly receive 403.
3. Build and test Functions in a nonproduction Firebase project/emulator with mocked or synthetic accounts. Exercise all six routes and browser preflight from each approved origin.
4. Build/export the frontend and verify the four active flows: create one, bulk create, delete one, and bulk delete. Deploying the Function before the updated client will cause old clients to receive 401.
5. Review deployed Function invoker/IAM settings, logs, and rollback plan. Preserve the previous artifact/revision for rapid rollback without reopening unauthenticated access longer than necessary.
6. Deploy the Function and frontend as one coordinated emergency release using the project's authorized process. Do not deploy from this task output automatically.
7. Immediately test unauthenticated, student, staff, admin, inactive, invalid-token, disallowed-origin, and approved-origin cases; monitor 401/403/5xx and privileged-operation logs.
8. Revoke suspicious sessions/accounts and begin the separate authorized Firestore/Storage rules inspection and V1 security review.
