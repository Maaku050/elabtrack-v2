# eLabTrack V2 frontend foundation

React + Vite + TypeScript + Tailwind v4; shadcn/ui `base-nova` primitives with Base UI and Lucide. `components.json` and all 62 existing UI primitives are preserved. A future eLabTrack design system will build on them during Phase 3.

```bash
npm ci
npm run dev
npm run lint
npm run test:run
npm run build
```

Use Node matching `package.json` engines. Copy `.env.example` to `.env` for host development; `VITE_API_URL=http://localhost:8080/api/v1`. Never expose backend secrets in `VITE_*`. Compose builds default to `/api/v1` through nginx.

Routes: `/` is the Phase 0 placeholder and theme toggle, `/status` is a user-triggered unauthenticated health read, and unknown routes show 404. There is no product dashboard, signup/login screen, equipment, borrowing, report, or kiosk workflow. API connectivity requires a running backend and database.

Query owns server state; Zustand owns suitable global client state; React owns local state. Feature hooks/APIs and shared transport should handle real workflows later. The health adapter uses the inherited health endpoint's raw response instead of the standard envelope. Generic auth/user adapters and auth-store tests remain unmounted and require Phase 1 review; both tokens currently persist in localStorage. See [foundation audit](../docs/project/FOUNDATION_AUDIT.md).
