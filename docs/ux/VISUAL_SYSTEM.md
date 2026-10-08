# Approved visual system — Phase 3B

Rebaselined 2026-10-08. **Phase 3A.2 approved; PHASE 3B COMPLETE AND OWNER VISUALLY APPROVED.** Real implementation screenshots, both themes, responsiveness and disclosed minor adaptations are accepted (DEC-064). The navy/indigo/violet PNGs in [approved](approved/README.md) supersede the earlier evergreen proposal. The old `mockups/` SVG package and [Phase 3A.2 report](../project/PHASE3A2_REPORT.md) remain historical design evidence; their green primary tokens are not implementation targets. This updates appearance only. [Phase 2.5 rules](../domain/BUSINESS_RULES.md) govern behavior; [fidelity contract](VISUAL_FIDELITY_CONTRACT.md) governs comparison and deviations.

## Palette evidence and semantic tokens

Chromium decoded the approved PNGs before implementation. Representative literal samples: B01 (350,55) RGB15/0/182; B02 (40,425) RGB34/3/186 and (720,923) RGB35/7/182; S01 (20,380) RGB7/31/84 and (72,107) RGB53/80/218; S01 (1000,72) RGB244/247/251. Large exact-color regions in remaining S05-05 include `#061741`, `#384CE6`, `#F5F7FD`, `#E4E8F3` and white. The core AI render has antialiasing/gradients, so tokens represent that family rather than treating every pixel as a new token.

Actual source: `frontend/src/styles/index.css`; composition geometry: `application.css`. Existing Tailwind v4/shadcn aliases remain compatible. Screen components consume semantic variables rather than their own palettes.

| Meaning | Light | Dark |
|---|---|---|
| Page / text | #F5F7FD / #10132E | #0D1428 / #EEF1FC |
| Surface | #FFFFFF | #17213A |
| Primary / label | #2B10BB / #FFFFFF | #B1A0FF / #17113C |
| Primary hover | #210796 | #C8BCFF |
| Muted / metadata | #F3F5FB / #606680 | #202C46 / #B7C2DC |
| Secondary / text | #F0EFFC / #4C4D74 | #26314F / #D7DFF5 |
| Accent | #EEEBFF | #302A55 |
| Decorative border | #E4E8F3 | #34415E |
| Essential input boundary / focus | #8993AD / #4C38DB | #7986A8 / #B5A3FF |
| Success / backing | #167346 / #EAF9F0 | #82DFAD / #17392E |
| Warning / backing | #925700 / #FFF5DD | #FFCF79 / #3C301C |
| Danger / backing | #B42332 / #FFF0F1 | #FF9CA8 / #422332 |
| Primary backing | #EFECFF | #2E2851 |
| Gold accent | #FFC529 | #FFC529 |
| Staff sidebar gradient | #061E52 → #061741 | Same institutional navy |
| Sidebar active / text | #384CE6 / #FFFFFF | Same readable pairing |
| Sidebar muted / divider / focus | #C5D0EB / #2C4270 / #C4B5FD | Same |
| Chart issued/pending/denied | #2514D3 / #94A4FF / #DFE4F2 | #A9B3FF / #7689DE / #4B5A87 |

Chart categories use the approved blue/lavender spectrum; semantic green/red/amber belongs to actionable states and metrics. Color is never the only status indicator. Dark is an engineering translation of the approved light references, **accepted by the owner through real implementation screenshots; no separate dark reference PNG was supplied**. Theme toggles reuse `useUIStore`; no replacement theme provider, persistence or credentials.

## Type, spacing and composition

Existing system sans font is retained; no font download/dependency. Borrower titles24px, Staff27px; section18px, operational panel16px; body/control14px; counts22–27px with tabular numerals; ordinary metadata12–13px. Borrower catalog tags10px and institutional caption7.5px are subordinate labels reflecting the dense source, not interactive instructions. Actual runtime text remains selectable, with semantic headings and visible control names. Font rasterization differs from the generated artwork; exact font identity is not supplied.

Spacing variables4/8/12/16/20/24/32px. Borrower mobile edges14px (core x≈33/941×390), tablet24px, max container1148px including edges; compact header64px; four labeled bottom destinations72px plus safe area. Content clears navigation by24px; the catalog summary sits above it. Cards12px radius, controls8px, pill filters, status backing7px; quiet one-pixel borders and shadow0 2 9 at4% navy light/16% black dark. Primary controls and effective checkbox targets44px minimum. Small typography and source controls are enlarged where needed; vertical scrolling replaces illegible compression.

Staff sidebar238px at1440,210px at1024,76px when manually collapsed; header64px; workspace24px (20 at1024); nav48px. Below1024 the menu opens an existing Sheet. Metrics4 columns at1024/1440; three operational panels at1440, two then one spanning panel at1024, stacked smaller. Pending rows retain identity-first density with named keyboard-focusable table scrolling at1024; desktop actions fit the viewport after shortening the fine badge. No whole-page horizontal scrolling.

Borrower catalog one column at320/390, two at768, three at1280. Mobile thumbnail/text/control relationships follow B02; at320 controls sit beneath item copy so44px quantities fit. Wide cards retain picture and facts, with an aligned action row. The four borrower destinations remain in one bottom-navigation composition at all preview widths; no unrelated desktop navigation. Home has2×2 mobile metrics and four wide metrics, then Browse/Activity/Quick Actions/reservation notice.

## State and accessibility foundation

Button/Card/Badge/Input/Checkbox/Sheet/Dialog/Table/Sidebar/Skeleton/Alert/Empty are retained Base UI/shadcn primitives. Scoped compositions adjust appearance without modifying the62 primitive files. Existing EmptyState/ErrorState/LoadingState remain available; new LoadingRegion and ConflictAlert provide geometry and persistent conflict presentation. Skeletons respect the existing global reduced-motion rule. Normal links/buttons have visible focus; Sheet/Dialog retain focus trap, Escape and restoration.

QuantitySelector is a local presentation draft: positive selections, zero removal, integer/bound errors with `aria-describedby`, disabled bound buttons. This is not inventory validation; later backend validation remains authoritative. Available/Low stock/Unavailable is an explicit supplied presentation state, not an invented stock threshold. Lifecycle labels are Pending/Active/Completed/Denied/Cancelled/Expired; overdue/physical/replacement/fine are separate meanings. No Approved/Released enum or borrower-return control.

No return-photo, evidence-upload, transmission, retention, storage or attachment feature. Catalog imagery is independent. Seal is a replaceable supplied AI reconstruction, not verified official master; [asset provenance](../../frontend/src/assets/brand/README.md). Synthetic preview thumbnails are catalog-only samples.

Actual contrast calculations, screenshots, responsive checks and limits are in [verification](verification/phase3b/README.md) and [Phase 3B report](../project/PHASE3B_REPORT.md). No WCAG certification or production feature validation is implied.
