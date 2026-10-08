import { createBrowserRouter, type RouteObject } from 'react-router-dom'
import { EntryRoute, ProtectedRoute, RouteFailure } from '@/features/auth/components/access-boundary'
import { borrowerPaths, staffPaths, adminPaths } from '@/features/auth/navigation'
import { StatusPage } from '@/features/foundation/pages/status-page'
import { NotFoundPage } from '@/components/common/not-found-page'
import { LoadingState } from '@/components/feedback/states'

// Vite removes this branch and its lazy fixture module from production builds.
const previewRoutes: RouteObject[] = import.meta.env.DEV ? [
  '/__preview/borrower/home', '/__preview/borrower/equipment',
  '/__preview/staff/dashboard', '/__preview/staff/requests', '/__preview/admin/dashboard',
].map(path => ({ path, hydrateFallbackElement: <LoadingState label="Loading visual preview…" />, lazy: () => import('@/features/visual-preview/routes') })) : []

export const routes: RouteObject[] = [{ errorElement: <RouteFailure />, hydrateFallbackElement: <LoadingState label="Opening your workspace…" />, children: [
  ...previewRoutes,
  { path: '/', element: <EntryRoute /> },
  { path: '/login', lazy: async () => ({ Component: (await import('@/features/auth/pages/login-page')).LoginPage }) },
  { element: <ProtectedRoute roles={['BORROWER']} />, children: borrowerPaths.map(path => ({ path, lazy: async () => ({ Component: (await import('@/features/workspace/workspace-page')).WorkspacePage }) })) },
  { element: <ProtectedRoute roles={['STAFF', 'ADMIN']} />, children: staffPaths.map(path => ({ path, lazy: async () => ({ Component: (await import('@/features/workspace/workspace-page')).WorkspacePage }) })) },
  { element: <ProtectedRoute roles={['ADMIN']} />, children: adminPaths.map(path => ({ path, lazy: async () => ({ Component: (await import('@/features/workspace/workspace-page')).WorkspacePage }) })) },
  { path: '/status', element: <StatusPage /> },
  { path: '*', element: <NotFoundPage /> },
] }]

export const router = createBrowserRouter(routes)
