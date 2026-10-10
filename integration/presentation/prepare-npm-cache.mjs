// Export only locked public tarballs for a fresh offline npm ci; no npmrc or
// credentials are copied. Requires the documented, installed npm toolchain.
import fs from 'node:fs/promises'
import path from 'node:path'
import {createRequire} from 'node:module'
import {execFileSync} from 'node:child_process'
const [lockPath,sourcePath,destinationPath]=process.argv.slice(2)
if(!lockPath||!sourcePath||!destinationPath)throw Error('Usage: node prepare-npm-cache.mjs PACKAGE_LOCK SOURCE_NPM_CACHE NEW_CACHE')
const source=path.resolve(sourcePath,'_cacache'),destination=path.resolve(destinationPath,'_cacache')
if(source===destination)throw Error('Source and destination cache must differ')
const globalRoot=execFileSync('npm',['root','-g'],{encoding:'utf8'}).trim()
const cache=createRequire(path.join(globalRoot,'npm/package.json'))('cacache')
const lock=JSON.parse(await fs.readFile(lockPath,'utf8'))
const entries=[...new Map(Object.values(lock.packages).filter(v=>v.resolved).map(v=>[v.resolved,v])).values()]
for(const entry of entries){const url=new URL(entry.resolved);if(url.origin!=='https://registry.npmjs.org'||url.username||url.password)throw Error('Only public locked npm registry tarballs may be exported')}
await fs.mkdir(destination,{recursive:true,mode:0o700})
let copied=0,missing=0,bytes=0
for(const entry of entries){
 const key='make-fetch-happen:request-cache:'+entry.resolved
 let info=await cache.get.info(source,key),data
 if(info){({data}=await cache.get(source,key))}
 else{try{data=await cache.get.byDigest(source,entry.integrity);info={integrity:entry.integrity,metadata:{}}}catch{missing++;continue}}
 const headers={};for(const name of ['content-type','content-length','cache-control','etag','last-modified','vary'])if(info.metadata?.resHeaders?.[name])headers[name]=info.metadata.resHeaders[name]
 const integrity=await cache.put(destination,key,data,{metadata:{time:Date.now(),url:entry.resolved,reqHeaders:{},resHeaders:headers,options:{compress:true}}})
 if(String(integrity)!==String(info.integrity))throw Error('Exported tarball integrity differs')
 copied++;bytes+=data.length
}
if(!copied)throw Error('No prepared locked tarballs found; prepare dependencies while online')
console.log(JSON.stringify({copied_public_locked_tarballs:copied,missing_platform_optional_or_unprepared_tarballs:missing,bytes,credentials_copied:false}))
