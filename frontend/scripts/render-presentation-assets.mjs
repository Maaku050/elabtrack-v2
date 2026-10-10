// Render original repository vector artwork to bounded canonical local PNG assets.
// No network, personal photographs or edits of existing images.
import fs from 'node:fs/promises'
import path from 'node:path'
import crypto from 'node:crypto'
import { browser } from './browser-cdp.mjs'
const root=path.resolve('../integration/presentation/assets'),manifest=[]
const b=await browser()
try{
 for(const directory of ['equipment','avatars'])for(const name of (await fs.readdir(path.join(root,directory))).filter(n=>n.endsWith('.svg')).sort()){
  const svg=await fs.readFile(path.join(root,directory,name)),url='data:image/svg+xml;base64,'+svg.toString('base64')
  const encoded=await b.evaluate(`(async()=>{const img=new Image();img.src=${JSON.stringify(url)};await img.decode();const canvas=document.createElement('canvas');canvas.width=256;canvas.height=256;canvas.getContext('2d').drawImage(img,0,0,256,256);return canvas.toDataURL('image/png').split(',')[1]})()`)
  const data=Buffer.from(encoded,'base64'),png=name.replace(/\.svg$/,'.png');if(data.length>512*1024)throw Error('Asset too large');await fs.writeFile(path.join(root,directory,png),data)
  manifest.push({file:`${directory}/${png}`,width:256,height:256,bytes:data.length,sha256:crypto.createHash('sha256').update(data).digest('hex'),source_sha256:crypto.createHash('sha256').update(svg).digest('hex')})
 }
 await fs.writeFile(path.join(root,'MANIFEST.json'),JSON.stringify({provenance:'Original project-owned vector illustrations, locally rendered by Chromium. Fictional presentation catalog and avatars, no external sources or real persons.',files:manifest},null,2)+'\n');console.log(JSON.stringify({equipment:manifest.filter(v=>v.file.startsWith('equipment/')).length,avatars:manifest.filter(v=>v.file.startsWith('avatars/')).length,external_requests:0}))
}finally{await b.close()}
