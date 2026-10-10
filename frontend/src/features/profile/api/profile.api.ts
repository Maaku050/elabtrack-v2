export interface OwnAccount { name:string; email:string; role:string; borrower_type:string; student_id:string; course:string; contact_number:string }
export interface ProfileImage { account_id: string; image_id: string | null; version: number }
const transport = async () => (await import('@/lib/api-client')).apiClient
export const profileApi = {
 own: async (signal?:AbortSignal):Promise<OwnAccount> => (await transport()).get('/profile-images/account/me',{signal}),
 metadata: async (id: string, signal?: AbortSignal): Promise<ProfileImage> => (await transport()).get(`/profile-images/${id}`,{signal}),
 image: async (v: ProfileImage, signal?: AbortSignal): Promise<Blob> => (await transport()).download(`/profile-images/${v.account_id}/${v.image_id}`,signal),
 save: async (v: ProfileImage, file: File, key: string): Promise<ProfileImage> => {const data = new FormData();data.append('image',file);data.append('expected_version',String(v.version));return (await transport()).post(`/profile-images/${v.account_id}`,data,{headers:{'Idempotency-Key':key,'Content-Type':undefined}})},
 remove: async (v: ProfileImage,key: string): Promise<ProfileImage> => (await transport()).patch(`/profile-images/${v.account_id}`,{expected_version:v.version,confirm:true},{headers:{'Idempotency-Key':key}}),
}
