# eLabTrack V2 frontend

React + TypeScript + Vite + Tailwind v4, retained shadcn/Base UI primitives and Lucide. Phase 4A adds real product login and protected navigation to the approved navy/violet visual system. All 62 UI primitives and `components.json` are preserved.

```bash
npm ci
npm run dev
npm run lint
npm run test:run
npm run build
```

Use the declared Node engine. Preserve existing ignored configuration; create `.env` from `.env.example` only if missing. Host Vite uses `VITE_API_URL=http://localhost:8080/api/v1`; Compose/nginx builds use `/api/v1`. No backend secret belongs in `VITE_*`.

`/` redirects to `/login` or the authenticated role's home. Login uses the existing centralized transport's `POST /api/v1/auth/login`, memory-only access JWT and HttpOnly refresh cookie. No public registration or password-recovery endpoint/UI exists. Accounts are provisioned by FSMO; secure onboarding/recovery and account management require Phase 4B/5 authorization. Local development has no built-in accounts or published passwords.

Borrower home is `/borrower/home`, with `/borrower/equipment`, `/borrower/borrowings`, `/borrower/account` and secondary `/borrower/notifications`. Staff home is `/staff/dashboard`, with `/staff/requests`, `/staff/inventory`, `/staff/borrowers`, `/staff/account` and `/staff/notifications`. Admin inherits Staff routes plus `/admin/reports` and `/admin/administration`. Admin does not impersonate a Borrower workspace. Exact route allowlists constrain post-login return destinations.

The existing SessionBootstrap restores once per document. Unknown/restoring authority hides protected content; expected anonymous 401 shows login without a persistent service warning. Network/server restoration failure offers deliberate Retry session. Every protected route checks current `/auth/me` before revealing content, including route changes, window focus and reconnect. Frontend guards are UX; every API resolves current PostgreSQL role/status. Demotion clears role-dependent private queries and increments the session generation to fence stale work.

Sign Out clears memory/private queries immediately, coordinates cookie revocation with the existing Web Lock and informs peers through strict non-secret BroadcastChannel events. Pending/failure feedback survives route changes; an unconfirmed server logout offers Retry sign out. The session cookie may remain after a failed network request until successful retry, expiry or replacement. Access JWTs are not globally revoked immediately; the existing Phase 1 semantics are preserved.

Authenticated workspaces show truthful feature placeholders and safe read-only current-account metadata. There are no live dashboard metrics, equipment records, borrowings or reports. Four approved synthetic visual previews remain development-only at `/__preview/borrower/home`, `/__preview/borrower/equipment`, `/__preview/staff/dashboard`, `/__preview/staff/requests`; the Admin preview is the existing navigation variant. Production excludes their routes and fixture module. `/status` remains the explicit unauthenticated health check; unknown routes show 404.

Pages use feature hooks → TanStack Query → feature APIs → the one transport. Query owns server state, Zustand owns global session/UI feedback, React owns local form state. React Hook Form + Zod validate the login form. There is no credential persistence or second auth provider. Existing theme persistence stores only preference.

See [Phase 4A report](../docs/project/PHASE4A_REPORT.md), [API contracts](../docs/API_CONTRACTS.md), [isolated reproduction guide](../integration/PHASE4A.md) and [rendered verification](../docs/ux/verification/phase4a/README.md). Browser checks use random-password accounts confined to a disposable PostgreSQL database, never public registration.
