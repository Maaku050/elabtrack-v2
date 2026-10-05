# eLabTrack V2 project charter

Status: confirmed current direction, 2026-10-05. Authority: the stakeholder's Phase 0 request. Engineering foundation work does not authorize product feature implementation.

## Purpose and reason

eLabTrack V2 is a ground-up rebuild, redesign, optimization, and security/maintainability improvement of the existing eLabTrack laboratory equipment borrowing operation. It exists to make the system reliable, supportable, testable, and accountable, then evaluate it in real FSMO operation before proposing broader adoption.

**eLabTrack V2 is NOT the campus-wide system.** Its operational scope is the existing Food Service Management System Organization (FSMO). No campus/college/department/laboratory hierarchy, tenancy, or cross-department borrowing is approved. Clean module boundaries should avoid needless hardcoding without building speculative future features.

## Relationship to V1 and future work

V1 provides a factual capability and technical-debt baseline, not architecture to copy. The capstone provides intended workflows and stakeholder expectations; discrepancies require explicit decisions. V2 replaces client-owned rules, direct database writes, UI-only authorization, race-prone stock arithmetic, global collection downloads, non-atomic completion/fine/history writes, destructive history deletion, duplicated logic, and weak testing with server authority and reliable persistence boundaries.

Results and lessons from a V2 pilot may inform a **separate** campus-wide equipment system proposal. That proposal has no implementation authorization in V2. Historical examples and undeployed V1 fixes do not prove current runtime behavior.

## Users and capabilities being modernized

Intended audiences are FSMO borrowers, laboratory operational staff, and responsible administrators. V1 represents students/staff/admin; the capstone also names faculty borrowers. Borrower eligibility and staff/admin authority remain open decisions, rather than accepted roles or permissions.

The capability family for later planning includes account access, borrower administration, aggregate equipment inventory/discovery, requests, approvals/denials, direct staff checkout, active/due/overdue borrowing, partial/final returns, damage/loss accountability, fines, history, email confirmations/reminders, dashboards, administrative reporting, and kiosk-oriented access. These are inputs to approved feature phases, not current implementations or an exhaustive committed backlog.

## Explicit non-goals

- Phase 0 implements no equipment, borrowing, returns, fines, reports, notifications, kiosk workflow, or migration.
- No business-domain tables before approved Phase 2 design.
- No tenancy, campus hierarchy, cross-department workflows, RFID/barcode/QR physical tracking, or native mobile applications.
- No microservices, Redis, Kafka, RabbitMQ, Kubernetes, GraphQL, CQRS, event sourcing, or distributed infrastructure without later explicit justification.
- No V1 Firebase access, data modification, deployment, commit, push, or automatic remote changes during Phase 0.

## Technical direction

One responsive React/TypeScript/Vite SPA and later kiosk-oriented interface → versioned REST Go/Fiber v3 API → PostgreSQL using pgx/pgxpool. Tailwind CSS v4, official shadcn/ui primitives, Lucide, TanStack Query, appropriate Zustand client state, React local state, React Hook Form, and Zod are accepted. Backend uses Clean Architecture and pragmatic DDD. SQL migrations, explicit constraints, and transactions govern persistence. Docker/Compose and Linux-compatible development are supported.

## Quality goals

Server-enforced validation/authorization; consistent API responses; safe session/secrets handling; transactional integrity under concurrent requests; preserved history and immutable audit evidence when domain workflows are implemented; bounded server queries; modular code with focused tests; accessible responsive interfaces; explicit configuration; reproducible development; observable errors without private data. Performance and pilot success thresholds must be measured and agreed in later phases rather than fabricated now.

Phase 0 exits when the source baseline, decision/open-question registers, roadmap and security backlog are reviewable, identity/demo cleanup is complete, private references stay private, and required validation evidence is available. Environment blockers prevent an unconditional exit claim; see `PHASE0_REPORT.md`.
