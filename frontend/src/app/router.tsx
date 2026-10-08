import { createBrowserRouter, type RouteObject } from 'react-router-dom'
import { FoundationPage } from '@/features/foundation/pages/foundation-page'
import { StatusPage } from '@/features/foundation/pages/status-page'
import { NotFoundPage } from '@/components/common/not-found-page'
import { LoadingState } from '@/components/feedback/states'

// Vite removes this branch and its lazy fixture module from production builds.
const previewRoutes: RouteObject[] = import.meta.env.DEV ? [
  '/__preview/borrower/home', '/__preview/borrower/equipment',
  '/__preview/staff/dashboard', '/__preview/staff/requests', '/__preview/admin/dashboard',
].map(path => ({ path, hydrateFallbackElement: <LoadingState label="Loading visual preview…" />, lazy: () => import('@/features/visual-preview/routes') })) : []

export const routes: RouteObject[] = [
  ...previewRoutes,
  { path: '/', element: <FoundationPage /> },
  { path: '/status', element: <StatusPage /> },
  { path: '*', element: <NotFoundPage /> },
]

export const router = createBrowserRouter(routes)
