import { Link } from 'react-router-dom'
import { Package, ListChecks, Clock, Hand, ArrowRight, PhilippinePeso, FileCheck2, ChevronRight, UserRound } from 'lucide-react'
import { AppButton } from '@/components/application/visual'
import { LoadingState } from '@/components/feedback/states'
import { BorrowerShell } from '@/components/application/shells'
import { borrowerDestinations } from '@/components/application/borrower-navigation'
import { useAuthStore } from '@/stores/auth-store'
import { apiErrorMessage } from '@/lib/api-error'
import { manilaTimeCompact } from '@/features/borrowing/lib/time'
import { pesos } from '@/features/borrowing/lib/money'
import { useBorrowings } from '@/features/borrowing/hooks/use-borrowing'
import { BorrowingCard } from '@/features/borrowing/components/borrowing-card'
import { useDashboard } from '../hooks/use-reporting'
import { Notice } from '@/components/application/management'
export function BorrowerHome(){
 const user=useAuthStore(s=>s.user)!,query=useDashboard(),recent=useBorrowings({page:1,per_page:3}),data=query.data
 const personal=[['active_loans','Active Borrowings',Package,'blue','/borrower/borrowings?status=CHECKED_OUT'],['pending_requests','Pending Requests',ListChecks,'amber','/borrower/borrowings?status=PENDING'],['overdue_loans','Overdue Borrowings',Clock,'red','/borrower/borrowings?status=CHECKED_OUT'],['outstanding_fines_minor','Outstanding Fines',PhilippinePeso,'purple','/borrower/borrowings']] as const
 return <BorrowerShell active="Home" destinations={borrowerDestinations}><div className="borrower-home"><header className="borrower-greeting"><h1>Hello, {user.name} <Hand className="greeting-hand" aria-hidden="true"/></h1><span>Ready for your next culinary activity?</span></header>{query.isPending&&<LoadingState/>}{query.error&&<Notice tone="danger">{apiErrorMessage(query.error)}<AppButton onClick={()=>{void query.refetch()}}>Retry metrics</AppButton></Notice>}{data&&<div className="borrower-personal-kpis">{personal.map(([key,label,Icon,tone,href])=><Link key={key} to={href} data-tone={tone}><span className="kpi-icon"><Icon aria-hidden="true"/></span><div><span>{label}</span><strong>{key.endsWith('_minor')?data.metrics[key]===undefined?'—':pesos(data.metrics[key]):(data.metrics[key]===undefined?'—':data.metrics[key].toLocaleString())}</strong></div></Link>)}</div>}<Link to="/borrower/equipment" className="borrower-browse-cta"><Package aria-hidden="true"/><span>Browse Equipment</span><ArrowRight aria-hidden="true"/></Link>
 <section className="borrower-home-section"><div className="section-heading"><h2>Recent Activity</h2><Link to="/borrower/borrowings">View all<ChevronRight aria-hidden="true"/></Link></div><p className="sr-only">Recent borrowing activities</p>{recent.isPending?<LoadingState/>:recent.error?<Notice tone="danger">{apiErrorMessage(recent.error)}<AppButton onClick={()=>{void recent.refetch()}}>Retry activity</AppButton></Notice>:recent.data?.items.length?<div className="borrower-recent-records">{recent.data.items.map(v=><BorrowingCard key={v.id} record={v} compact/>)}</div>:<p className="home-empty">Your borrowing activity will appear here.</p>}</section>
 <section className="borrower-home-section"><h2>Quick Actions</h2><div className="borrower-quick-actions">{[[FileCheck2,'Start Request','/borrower/equipment','Choose equipment'],[ListChecks,'My Borrowings','/borrower/borrowings','Status and due dates'],[UserRound,'Account','/borrower/account','Your profile']].map(([Icon,label,href,caption])=>{const I=Icon as typeof ListChecks;return <Link key={String(href)} to={String(href)}><I aria-hidden="true"/><span>{String(label)}</span><small>{String(caption)}</small></Link>})}</div></section><Notice tone="warning"><strong>Borrowing information</strong><p>Submitting a request reserves available equipment for up to 24 hours. Visit FSMO for Staff/Admin to record physical handover.</p></Notice><p className="borrower-local-context">Your account only · {data?manilaTimeCompact(data.as_of):'Current recorded data'} · Asia/Manila</p></div></BorrowerShell>
}
