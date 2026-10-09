import { useDeferredValue } from 'react'
import { Link } from 'react-router-dom'
import { Plus, ArrowUpRight, Check, Clock } from 'lucide-react'
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table'
import { DirectoryHeader, DirectoryPanel, DirectorySearch, DirectorySelect, DirectoryTableRegion } from '@/components/application/directory'
import { ServerPagination } from '@/components/application/server-pagination'
import { AppButton, ToneBadge } from '@/components/application/visual'
import { useDirectoryState } from '@/hooks/use-directory-state'
import { useAuthStore } from '@/stores/auth-store'
import { apiErrorMessage } from '@/lib/api-error'
import { useAccounts } from '../hooks/use-accounts'
import { activationLabel } from '../lib/account-status'
import { BorrowerDirectory } from './borrower-directory'

function StaffAccounts({ embedded = false }: { embedded?: boolean }) {
 const admin=useAuthStore(s=>s.user?.role==='ADMIN'),view=useDirectoryState()
 const prefix=embedded?'staff_':'',pageKey=prefix+'page',searchKey=prefix+'search',statusKey=prefix+'status'
 const page=view.page(pageKey),search=view.value(searchKey),deferred=useDeferredValue(search),status=['ACTIVE','INACTIVE','PENDING'].includes(view.value(statusKey))?view.value(statusKey):''
 const query=useAccounts({page,per_page:25,search:deferred,borrower_type:'',status},true),data=query.data,base='/admin/administration/accounts'
 function filter(changes:Record<string,string>,replace=false){view.update({...changes,[pageKey]:''},false,replace)}
 const panel=<DirectoryPanel title="Staff / Admin account directory" summary={data&&!query.isError?`${data.total.toLocaleString()} matching accounts`:'Restricted account administration'} pending={query.isPending} error={query.isError?apiErrorMessage(query.error):undefined} empty={!data?.items.length} onRetry={()=>{void query.refetch()}} toolbar={<><DirectorySearch label="Search Staff / Admin accounts" placeholder="Name or email" value={search} onChange={e=>filter({[searchKey]:e.target.value},true)} /><DirectorySelect label="Status" value={status} onChange={e=>filter({[statusKey]:e.target.value})}><option value="">All statuses</option><option value="ACTIVE">Active</option><option value="INACTIVE">Inactive</option><option value="PENDING">Activation pending</option></DirectorySelect>{(search||status)&&<AppButton variant="ghost" onClick={()=>filter({[searchKey]:'',[statusKey]:''})}>Clear filters</AppButton>}</>} footer={data&&<ServerPagination label="Staff / Admin accounts" pageKey={pageKey} page={data.page} perPage={data.per_page} total={data.total} count={data.items.length} totalPages={Math.ceil(data.total/data.per_page)} pending={query.isFetching} onPage={p=>view.update({[pageKey]:p===1?'':String(p)},false)} />}>
 <DirectoryTableRegion label="Staff / Admin accounts table"><Table><TableHeader><TableRow>{['Account','Role','Status','Activation'].map(label=><TableHead key={label} scope="col">{label}</TableHead>)}<TableHead scope="col"><span className="sr-only">Details</span></TableHead></TableRow></TableHeader><TableBody>{data?.items.map(account=><TableRow key={account.id}><TableCell><div className="directory-person"><span className="directory-person-avatar" aria-hidden="true">{account.name.trim().split(/\s+/).slice(0,2).map(n=>n[0]).join('').toUpperCase()}</span><div><Link className="directory-record-name" to={`${base}/${account.id}`}>{account.name}</Link><span title={account.email}>{account.email}</span></div></div></TableCell><TableCell><ToneBadge tone={account.role==='ADMIN'?'primary':'neutral'}>{account.role}</ToneBadge></TableCell><TableCell><ToneBadge tone={account.is_active?'success':'neutral'} icon={Check}>{account.is_active?'Active':'Inactive'}</ToneBadge></TableCell><TableCell><ToneBadge tone={account.activation_required?'warning':'neutral'} icon={account.activation_required?Clock:Check}>{activationLabel(account.activation_required)}</ToneBadge></TableCell><TableCell><AppButton variant="ghost" nativeButton={false} render={<Link to={`${base}/${account.id}`} aria-label={`Open ${account.name}`} />}>Open<ArrowUpRight aria-hidden="true" /></AppButton></TableCell></TableRow>)}</TableBody></Table></DirectoryTableRegion></DirectoryPanel>
 return embedded?panel:<div className="directory-page"><DirectoryHeader title="Staff & Admin Accounts" description="Manage Staff accounts and review existing Admin accounts. New Admin provisioning requires a separate approved security policy." eyebrow="Administration" actions={admin&&<AppButton nativeButton={false} render={<Link to={`${base}/new`} />}><Plus aria-hidden="true" />Add Staff</AppButton>} />{panel}</div>
}
export function StaffDirectory(){return <StaffAccounts />}
export function EmbeddedStaffDirectory(){return <StaffAccounts embedded />}
export function AccountDirectory({staff=false,embedded=false}:{staff?:boolean;embedded?:boolean}){return staff?<StaffAccounts embedded={embedded}/>:<BorrowerDirectory />}
