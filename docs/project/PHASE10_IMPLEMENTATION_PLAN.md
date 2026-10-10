# Phase 10 — truthful dashboards and bounded reports

Prepared under DEC-082; implementation begins only after Phase9 exit evidence is recorded. No migrations are needed for read-only projections. No normal database changes.

Use one PostgreSQL snapshot per dashboard/report response, original custody buckets without joins that multiply stock, exact paginated report counts, deterministic timestamp/ID ordering, and current account/role checks inside application transactions. Borrower dashboard shows only their own obligations/workflows; operational metrics cover the actual FSMO inventory and loans. Financial amounts are integer centavos and original due/completion assessment rules remain authoritative.

Reports: inventory, stock movements, requests, issued loans, active checkouts, overdue loans, returns, damaged/lost incidents, replacement obligations, fine assessment/clearance, and account audit. Staff/Admin may read operational reports; account audit is Admin-only. Bounded date/search/status/category/equipment/borrower filters apply on the server before pagination/count/export. Date bounds are inclusive start/exclusive end; report-specific unsupported filters are rejected. CSV uses identical filtered data, deterministic columns, UTF-8 and RFC4180 escaping, formula-safe cells, and a maximum5000 rows; oversized exports require narrower filters rather than silent truncation.

Reuse the persistent shell, approved metric cards, tables, directory controls and TanStack Query transport. No invented revenue, rental costs, payment checkout, chart data or institutional scheduling. Gate: source reconciliation, pagination/filter/date behavior, export bounds/formula security, role/deactivation checks, frontend quality and real Chromium UI/export checks before Phase11.
