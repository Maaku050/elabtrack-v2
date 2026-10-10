# Human-review UI polish progress

**Current owner acceptance and publication authority, 2026-10-11 — DEC-085:** Final Staff/Admin desktop and Student/Faculty mobile visual review passed; the owner reported no remaining UI issues and authorized auditing, committing and pushing the completed Phases 8–14 application and shared UI to `release/elabtrack-v2-presentation-ready`. Main and production deployment remain untouched. Normal schema 000008 and existing records are preserved. Official FSMO terms/current acceptance, production Student domains, verified activation/Brevo and operational deployment readiness remain separate gates. Earlier pending-review/no-commit statements below are historical task snapshots. [Release audit and publication handoff](RELEASE_PUBLICATION_REPORT.md).

Authority: owner’s final human-review request and both human-review DOCX documents, 2026-10-10. This is a focused refinement of DEC-083, not another full reconstruction or a business-policy change. Owner visual acceptance remains pending. No commit, push, deployment or normal migration/reset.

Both documents were read completely and all 18 embedded images inspected individually: desktop D-01–04, R-01–02, I-01–04; mobile H-01, E-01–03, A-01–04. Private original/extracted references remain ignored. New BEFORE screenshots use the current isolated fictional presentation; earlier reconstruction evidence is retained.

| Requirement | Implementation | Verification status |
|---|---|---|
| HR-01 | Shared themed native scrollbars; Dashboard composition retained | PASS / COMPLETE |
| HR-02 | Compact operational queue, supported status pills/server totals, borrower identity, grounded columns, Review/Open, bounded table/pager | PASS / COMPLETE |
| HR-03 | Existing wizard retained; bounded independent catalog/cart scroll, accessible footer, compact controls | PASS / COMPLETE |
| MR-01 | Accepted Home source/composition retained | PASS / ACCEPTED AND PRESERVED |
| MR-02A | Existing Base UI portalled responsive popover, Escape/outside/focus | PASS / COMPLETE |
| MR-02B | Compact adjacent minus/input/plus/Add; mobile 44px targets | PASS / COMPLETE |
| MR-03A | One red Account heading Sign Out; existing logout hook | PASS / COMPLETE |
| MR-03B | Existing spacing scale: 24px above terms card | PASS / COMPLETE |

No backend, API contract, schema, inventory arithmetic or policy change is needed. Status totals use seven bounded existing list queries (`per_page=1`) cached for 30 seconds under the existing authenticated borrowing query prefix. Search and borrower-ID constraints intersect each count; status selection resets server paging. Failed counts are unknown rather than zero. Existing borrowing-command invalidation refreshes these display counts; they never determine authority/eligibility.

Baseline: normal development 21 tables and accepted Phase 7 demo 22 tables have fresh read-only private digests. All real browser accounts/data are isolated fictional presentation fixtures. Existing demo records are retained; no reseed/reset.

Final verification completed 2026-10-11: 334 frontend tests/32 files, TypeScript/build/lint, Go fmt/vet/tests, 8 presentation + 8 shared operational + 6 actual workflow Chromium groups, 305 new AFTER/workflow screenshots, 11 presentation integrity categories, 21/22-table preservation and source proof PASS. Final secret/Git audit is recorded in the report/machine evidence. Owner visual acceptance remains pending. Normal presentation bundle restored; no commit/push/deployment.

Final Git/secret closeout: main/c2741c5 unchanged, 78 modified tracked / 954 untracked public / zero staged. 2,917 public files compared with 378 private values from 48 configurations: zero secret matches. Exact cumulative inventory and 25-path task delta are in the report. All 332 PNG artifacts/gallery links pass validation. Engineering verified; owner visual acceptance pending.
