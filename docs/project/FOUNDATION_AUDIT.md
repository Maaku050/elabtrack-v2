# Phase 0 foundation audit

Inspected 2026-10-05. Classification describes disposition, not production approval: **KEEP** useful infrastructure, **ADAPT** change/project-hardening needed, **REMOVE** inherited demonstration/workflow removed, **DEFER** wait for an approved later phase. Evidence comes from source, not README claims. No database, Firebase, or deployment inspection occurred.

## Before adaptation: actual supplied starter

Root contained `AGENTS.md`, `.specify/`, `.devin/skills/`, Make wrappers, Windows `scripts/dev.bat`, three env examples, Compose, generic READMEs, architecture/stack/planner docs, and pervasive `*:Zone.Identifier` sidecars. `.git` is a read-only empty directory in this session; `git status` and `git remote -v` both fail. No current numbered feature folder or runtime spec dependency exists; unrelated template/Sarcita specifications were inside bundled workflow references only.

Backend has Go 1.27.1/Fiber v3.5.0, application/domain auth and users, pgx repositories, JWT/bcrypt adapters, env loader, migrator/seeder, HTTP middleware and unit/HTTP tests. Both backend and frontend Dockerfiles and frontend nginx configuration actually exist, despite inherited `docs/STACK.md` claiming Dockerfiles were absent.

Frontend has React/Vite/TypeScript, Tailwind v4, `components.json`, **62** shadcn/Base UI primitives, TanStack Query, Zustand, Hook Form/Zod, Axios, auth/profile/dashboard/showcase routes and **62** component example files (including `example.tsx`). The original root route redirects to an auth-gated demo dashboard. `Providers` hydrates persistent tokens; `/app/components` imports the showcase. npm packages were not installed initially. No equipment/borrowing/report/kiosk implementation is present.

References supplied at `.project-reference/ELABTRACK-SEMIFINAL.docx` and `docs/v1-audit/`; the former was already ignored by `.gitignore`. The V1 source repository is not in this workspace. The audit includes 15 numbered documents, including an undeployed emergency hardening supplement.

## Disposition and evidence

| Area | Class | Source evidence and disposition |
|---|---|---|
| Backend layers | KEEP / ADAPT | `internal/{domain,application,infrastructure,interface/http,bootstrap}` separates use cases/repos/HTTP; composition in `bootstrap/dependencies.go`. Domain/application do not import Fiber/pgx/infrastructure. Retain; health handler directly imports database health adapter and needs a port. |
| Pragmatic DDD | KEEP / DEFER | Generic user value objects/entity/repository and auth entity/repository exist. Product boundaries/invariants belong to approved Phase 2; no speculative aggregate scaffolds. |
| PostgreSQL integration | KEEP / ADAPT | One pgxpool in `infrastructure/database/postgres.go`; parameterized repositories and bounded user list (`Normalize(100)`). `config/database.go` hardcodes sslmode=disable and partially escapes URL credentials; production TLS/config needs review. |
| Transactions | KEEP / ADAPT | `database/transaction.go` injects pgx transactions into repository context. Auth registration/refresh does not use it; application-facing transaction port and atomic auth behavior needed. |
| SQL migrations | KEEP / ADAPT | Paired 000001 users/citext and 000002 refresh_tokens. UUID keys, email/token uniqueness, FK, role/is_active constraints/indexes. Refresh tokens stored raw; role allows only user/admin and is provisional. No eLabTrack tables. |
| Migration runner | ADAPT | SQL runs in a transaction but schema_migrations update occurs separately afterward (also Down); no advisory lock/checksums. CLI migrator creation bootstraps database first. Review tracking atomicity, concurrent runners and production credentials. |
| Seeds | ADAPT | `seeds/development.sql` contains synthetic admin/user with published weak passwords; seeder has a comment against production but no environment guard. Display names changed to Development Admin/User only. No seed executed. |
| Configuration | ADAPT | `config/config.go`, dotenv/env helpers: environment defaults for DB, JWT, origins, timeouts; no production fail-closed secret/origin/TLS validation. App/database/issuer defaults and env examples renamed. |
| Docker/Compose | KEEP / ADAPT | Both Dockerfiles exist; non-root API runtime, SPA nginx fallback/proxy, PostgreSQL health dependency. Resource names/default DB adapted; backend DB_USER/DB_NAME aligned with root Compose credentials; full frontend defaults to same-origin `/api/v1`. API Docker CMD runs migrations automatically; Phase 1 must separate deployment migration responsibility. |
| CORS | KEEP / ADAPT | `middleware/cors.go`: allowlist, credentialed requests, explicit headers/methods/request-id exposure. Config defaults to localhost origins even in production; no production-origin validation. CORS is not authorization. |
| Rate limiting | KEEP / ADAPT | `middleware/rate_limit.go`: global IP limiter using c.IP(), 120/min default, health exempt, normal 429 envelope. No endpoint-specific login/refresh controls or trusted-proxy config in newServer; memory scope is one API process. |
| Security headers | KEEP / ADAPT | Helmet middleware supplies CSP/frame/referrer/HSTS/etc for API responses. Review deployment/TLS compatibility; nginx serves frontend separately without equivalent explicit header policy. |
| Logging/recovery | KEEP / ADAPT | zap application logger; Fiber request logger logs method/path/status/latency/request ID, not bodies/tokens. HTTP logs remain plaintext even with JSON app log setting. Recovery and safe error mapper exist; custom panic log callback has EnableStackTrace=false, so do not assume it runs. Review observability. |
| Validation | KEEP / ADAPT | Fiber Bind Body, validator DTO checks, user value-object validation, parameterized SQL. Validation is backend-authoritative; production config needs separate validation. |
| Response envelope | KEEP / ADAPT | `response.Body` normalizes success/errors, 204 body omitted; `HealthHandler.Health` instead emits raw status/services. Frontend ApiError incorrectly expects message inside error, and toRequestError accepts a different shape; align contracts later. |
| Auth/session | ADAPT | Public register/login/refresh/logout, bearer /me, admin user list. Registration creates generic user and issues tokens; no product provisioning/verification policy. Login/refresh check active status, ordinary bearer middleware checks token only. |
| Access JWT | ADAPT | HS256 issued with issuer/uid/subject/role/expiry. Verify accepts HMAC family and does not require configured issuer, exact HS256, expiry presence, nonzero identity or subject consistency. Default signing secret allowed; role/status can become stale. |
| Refresh token | ADAPT | crypto/rand 32 bytes → hex; DB stores token string with expiry/revoked. Rotation finds/checks then revokes then issues without transaction or compare-and-swap; already-revoked row can still be updated. No family reuse detection. Logout revokes supplied refresh token but leaves access token valid. |
| Frontend architecture | KEEP / ADAPT | `app/`, feature API/hooks/schemas/types, shared client, query keys, stores and strict tsconfig. Auth/user modules remain provisional/unmounted; removed page exports and obsolete navigation targets. Placeholder has independent health feature. |
| Frontend token persistence | ADAPT | `lib/storage.ts` wraps localStorage; auth-store writes both tokens. Namespace changed fst.* → elabtrack_v2.* without changing security behavior. Auth provider hydration removed from placeholder. Old persistence is not accepted V2 session design. |
| Axios refresh | ADAPT | Central bearer injection and in-instance single-flight refresh. No retried-request marker bounds repeated 401 after successful refresh; refresh changes storage without synchronizing store fields. Cross-tab rotation/cache/session policy needs review. |
| shadcn configuration | KEEP | `components.json`: base-nova, rsc=false, tsx=true, zinc CSS variables, Lucide, aliases to components/ui/lib/hooks. Every original UI primitive byte-for-byte retained. No component framework replacement. |
| Styling/state/forms | KEEP / DEFER | Tailwind v4 Vite plugin, tokens/dark/reduced-motion styles; Query server state, appropriate Zustand UI and local React state; Hook Form/Zod available. Full eLabTrack design system and shell deferred to Phase 3. |
| Test infrastructure | KEEP / ADAPT | Go domain auth/user/bcrypt unit tests + bootstrap HTTP middleware tests, not real PostgreSQL integration. Vitest/jsdom/RTL and auth-store/utils tests retained; Phase 0 boundary tests added. No CI workflow in supplied tree. Gates reported separately. |
| Demo/showcase | REMOVE | Removed showcase and dashboard feature folders, demo layout, auth/profile example screens, root *-example.tsx and example.tsx. Auth/user adapters remain reusable starter infrastructure; no demonstration route is registered. |
| Optional workflow | REMOVE | Entire `.specify/`, `.devin/`, `docs/SPECKIT-PLANNER.md`, bundled unrelated historical examples and private/project-owner instructions. Searches show no runtime/build/test dependency. Useful architecture/quality rules rewritten in AGENTS/ARCHITECTURE and project docs. No partial broken workflow retained. |
| Windows sidecars | REMOVE | Removed all supplied `*:Zone.Identifier` files and ignored future sidecars. The optional dev.bat remains; test command aligned with test:run. |
| Reference hygiene | KEEP / ADAPT | Private DOCX remains local, unchanged and ignored; no biographies/sample contact data copied. V1 audit moved to docs/reference/v1-audit without altering content. Historical audits remain evidence, not runtime promises. |
| Product features | DEFER | Phases 2–14 per roadmap; Phase 0 does not create product tables, rules, kiosk, reports, notifications or migration. |

## After adaptation

Frontend now mounts `/`, `/status` and 404 without requiring a session. It demonstrates React/routing/theme and preserved shadcn card/button primitives. The health read is manual, unauthenticated, cookie-free and schema-checked; it adapts the known raw health response. Query retry is disabled for this diagnostic, so users explicitly retry errors. It does not claim real API connectivity without a running database-backed API.

Identity: eLabTrack V2 app/title/frontend package; `elabtrack_v2` DB; `elabtrack-v2` issuer; elabtrack_v2_postgres/backend/frontend and explicitly named elabtrack_v2_pgdata volume. Go module/imports stay inherited because the authoritative remote cannot be inspected. Naming changes require fresh development credentials/tokens; they do not migrate any existing DB or Docker volume.

Inherited specialist routing was inspected before removal. The bundled custom-ui-shadcn reference forbade shadcn use, contradicting the explicit stakeholder choice; it did not govern this adaptation. Useful architecture/security/accessibility checks were incorporated directly, without depending on unavailable plugins, private notes or unrelated projects. This task is a rebaseline, not a full speculative feature proposal; no numbered proposal/ZIP is required.

## Verification limits

Refer to `PHASE0_REPORT.md` for actual gate results. No local toolchain upgrade, database mutation, Docker run, browser end-to-end session, Firebase access, commit, push or deployment occurred. Git/Go/WSL limitations prevent an unconditional Phase 0 exit claim. Static code inspection cannot prove production safety.
