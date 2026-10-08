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
  { path: '/activate', lazy: async () => ({ Component: (await import('@/features/accounts/pages/activation-page')).ActivationPage }) },
  { path: '/login', lazy: async () => ({ Component: (await import('@/features/auth/pages/login-page')).LoginPage }) },
  { element: <ProtectedRoute roles={['BORROWER']} />, children: [
    { path: '/borrower/terms', lazy: async () => ({ Component: (await import('@/features/terms/pages/terms-page')).TermsPage }) },
    { lazy: async () => ({ Component: (await import('@/features/terms/components/terms-gate')).TermsGate }), children: [{path:'/borrower/home',lazy:async()=>({Component:(await import('@/features/workspace/workspace-page')).WorkspacePage})},{path:'/borrower/equipment',lazy:async()=>({Component:(await import('@/features/inventory/pages/equipment-list')).BorrowerCatalog})},{path:'/borrower/equipment/:id',lazy:async()=>({Component:(await import('@/features/inventory/pages/equipment-detail')).BorrowerEquipmentDetail})}] },
    ...borrowerPaths.filter(path => path !== '/borrower/home' && path !== '/borrower/equipment').map(path => ({ path, lazy: async () => ({ Component: (await import('@/features/workspace/workspace-page')).WorkspacePage }) })),
  ] },
  { element: <ProtectedRoute roles={['STAFF', 'ADMIN']} />, children: [
 {path:'/staff/inventory',lazy:async()=>({Component:(await import('@/features/inventory/pages/equipment-list')).EquipmentList})},
 {path:'/staff/inventory/new',lazy:async()=>({Component:(await import('@/features/inventory/pages/equipment-form')).EquipmentForm})},
 {path:'/staff/inventory/categories',lazy:async()=>({Component:(await import('@/features/inventory/pages/category-management')).CategoryManagement})},
 {path:'/staff/inventory/:id',lazy:async()=>({Component:(await import('@/features/inventory/pages/equipment-detail')).EquipmentDetail})},
 {path:'/staff/inventory/:id/edit',lazy:async()=>({Component:(await import('@/features/inventory/pages/equipment-form')).EquipmentForm})},
 {path:'/staff/inventory/:id/adjust',lazy:async()=>({Component:(await import('@/features/inventory/pages/stock-adjustment')).StockAdjustment})},
 {path:'/staff/borrowers',lazy: async()=>({Component:(await import('@/features/accounts/pages/account-directory')).AccountDirectory})},
 {path:'/staff/borrowers/:id',lazy: async()=>({Component:(await import('@/features/accounts/pages/account-detail')).AccountDetail})},
 ...staffPaths.filter(path=>path!='/staff/borrowers'&&path!='/staff/inventory').map(path => ({ path, lazy: async () => ({ Component: (await import('@/features/workspace/workspace-page')).WorkspacePage }) }))] },
  { element: <ProtectedRoute roles={['ADMIN']} />, children: [
 {path:'/admin/inventory/:id/reconcile',lazy:async()=>({Component:(await import('@/features/inventory/pages/stock-adjustment')).InventoryReconciliation})},
 {path:'/staff/borrowers/new',lazy:async()=>({Component:(await import('@/features/accounts/pages/account-create')).AccountCreate})},
 {path:'/admin/borrowers/bulk',lazy:async()=>({Component:(await import('@/features/accounts/pages/student-bulk')).StudentBulk})},
 {path:'/admin/administration',lazy:async()=>({Component:(await import('@/features/accounts/pages/administration-page')).AdministrationPage})},
 {path:'/admin/administration/accounts',lazy:async()=>({Component:(await import('@/features/accounts/pages/account-directory')).StaffDirectory})},
 {path:'/admin/administration/accounts/new',lazy:async()=>({Component:(await import('@/features/accounts/pages/account-create')).StaffCreate})},
 {path:'/admin/administration/accounts/:id',lazy:async()=>({Component:(await import('@/features/accounts/pages/account-detail')).StaffDetail})},
 ...adminPaths.filter(path=>path!='/admin/administration').map(path => ({ path, lazy: async () => ({ Component: (await import('@/features/workspace/workspace-page')).WorkspacePage }) }))] },
  { path: '/status', element: <StatusPage /> },
  { path: '*', element: <NotFoundPage /> },
] }]

export const router = createBrowserRouter(routes)
