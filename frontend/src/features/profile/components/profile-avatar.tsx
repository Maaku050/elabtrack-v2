import { useEffect, useRef, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useAuthStore } from '@/stores/auth-store'
import { useProfile } from '../hooks/use-profile'
import { profileApi } from '../api/profile.api'
function Initials({name,className}:{name:string;className?:string}){return <span className={`profile-avatar ${className??''}`} aria-label={`${name} initials`}>{name.trim().split(/\s+/).slice(0,2).map(v=>v[0]).join('').toUpperCase()||'?'}</span>}
export function ProfileAvatar({id,name,className}:{id:string;name:string;className?:string}) {
 const actor=useAuthStore(s=>s.user?.id)
 return actor&&id?<ProtectedAvatar id={id} name={name} className={className} />:<Initials name={name} className={className} />
}
function ProtectedAvatar({id,name,className}:{id:string;name:string;className?:string}) {
 const actor=useAuthStore(s=>s.user),meta=useProfile(id),ref=useRef<HTMLImageElement>(null),[failed,setFailed]=useState('')
 const image=useQuery({queryKey:['profile',actor?.id,actor?.role,id,'image',meta.data?.image_id],queryFn:({signal})=>profileApi.image(meta.data!,signal),enabled:!!meta.data?.image_id,staleTime:300_000,retry:false,meta:{authenticated:true}})
 useEffect(()=>{if(!image.data||!ref.current)return;const url=URL.createObjectURL(image.data);ref.current.src=url;return()=>URL.revokeObjectURL(url)},[image.data,meta.data?.image_id])
 return image.data&&meta.data?.image_id&&failed!==meta.data.image_id&&!image.isError?<img ref={ref} className={`profile-avatar ${className??''}`} alt={`${name} profile image`} onError={()=>setFailed(meta.data!.image_id!)} />:<Initials name={name} className={className} />
}
