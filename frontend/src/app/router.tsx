import { createBrowserRouter, type RouteObject } from 'react-router-dom'
import { FoundationPage } from '@/features/foundation/pages/foundation-page'
import { StatusPage } from '@/features/foundation/pages/status-page'
import { NotFoundPage } from '@/components/common/not-found-page'

export const routes: RouteObject[] = [
  { path: '/', element: <FoundationPage /> },
  { path: '/status', element: <StatusPage /> },
  { path: '*', element: <NotFoundPage /> },
]

export const router = createBrowserRouter(routes)
