import assert from 'node:assert/strict'
import { readFile, writeFile } from 'node:fs/promises'
import { execFileSync } from 'node:child_process'
import { randomUUID, createHash } from 'node:crypto'
import { pathToFileURL, fileURLToPath } from 'node:url'
import path from 'node:path'

const root=path.resolve(path.dirname(fileURLToPath(import.meta.url)),'..')
const {chromium}=await import(pathToFileURL(process.env.PLAYWRIGHT_MODULE).href)
const base='http://localhost:15173'
const email=`phase1g-restart-${randomUUID()}@example.invalid`, password=randomUUID()+'-Synthetic-Only'
const secrets=JSON.parse(await readFile('/tmp/elabtrack-phase1g-secrets.json','utf8'));secrets.push(password)
const sql=query=>execFileSync('psql',['-h','127.0.0.1','-p','15432','-U','postgres','-d','elabtrack_v2_integration','-At'],{input:query,env:{...process.env,PGPASSWORD:process.env.DB_PASSWORD},encoding:'utf8'}).trim()
const fixtureResponse=await fetch(base+'/api/v1/auth/register',{method:'POST',headers:{Origin:base,'Content-Type':'application/json'},body:JSON.stringify({email,password,name:'Synthetic restart account'})})
assert.equal(fixtureResponse.status,201)
const fixture=(await fixtureResponse.json()).data
const id=fixture.user.id;assert.match(id,/^[0-9a-f-]{36}$/)
secrets.push(fixture.access_token,fixtureResponse.headers.get('set-cookie').split(';')[0].split('=')[1])
const browser=await chromium.launch({headless:true})
try{
  const context=await browser.newContext(), page=await context.newPage()
  page.on('request',req=>{const bearer=req.headers().authorization?.replace(/^Bearer /,'');if(bearer)secrets.push(bearer)})
  await page.goto(base+'/phase1g.local/index.html')
  await page.waitForFunction(()=>window.phase1g?.state().status==='unauthenticated')
  await page.evaluate(({email,password})=>window.phase1g.login(email,password),{email,password})
  const cookie=(await context.cookies()).find(c=>c.name==='elabtrack_v2_refresh');secrets.push(cookie.value,createHash('sha256').update(cookie.value).digest('hex'))
  const before=sql("SELECT string_agg(version||':'||applied_at::text,',' ORDER BY version) FROM schema_migrations")
  execFileSync(process.env.COMPOSE_BIN,['-p','elabtrack_v2_phase1g','--env-file','backend/.env.phase1g','-f','docker-compose.yml','-f','integration/compose.override.yml','--profile','full','restart'],{cwd:root,stdio:'pipe',timeout:60000})
  // Compose restart returns once processes start, before PostgreSQL readiness.
  const deadline=Date.now()+15000
  let ready=false
  while(Date.now()<deadline){
    try{ready=(await fetch(base+'/api/v1/ready')).status===200}catch{}
    if(ready)break
    await new Promise(resolve=>setTimeout(resolve,250))
  }
  assert(ready,'API readiness after isolated restart')
  assert.equal(sql("SELECT string_agg(version||':'||applied_at::text,',' ORDER BY version) FROM schema_migrations"),before)
  assert.equal(sql(`SELECT count(*) FROM users WHERE id='${id}'`),'1')
  assert((await context.cookies()).find(c=>c.name==='elabtrack_v2_refresh').value===cookie.value,'cookie persists across restart')
  assert.equal((await page.evaluate(()=>window.phase1g.me())).id,id)
  await page.reload()
  await page.waitForFunction(()=>window.phase1g?.state().status==='authenticated')
  const replacement=(await context.cookies()).find(c=>c.name==='elabtrack_v2_refresh');assert(replacement.value!==cookie.value,'reload rotates cookie');secrets.push(replacement.value,createHash('sha256').update(replacement.value).digest('hex'))
  await page.evaluate(()=>window.phase1g.logout())
  const document=await page.goto(base+'/status')
  assert.equal(document.status(),200)
  await page.getByRole('heading',{name:'eLabTrack V2 service connection'}).waitFor()
  await writeFile('/tmp/elabtrack-phase1g-restart.json',JSON.stringify({postgresVolumePreserved:true,migrationVersionsAndTimestampsUnchanged:true,oldAccessWorks:true,cookieSurvives:true,reloadRotatesAndRestores:true,logoutWorks:true,nginxSPAAvailable:true},null,2))
  console.log('PASS full Compose restart preserves DB, explicit migration state and browser session; reload restores and logout succeeds')
}finally{
  sql(`DELETE FROM users WHERE id='${id}'`)
  await writeFile('/tmp/elabtrack-phase1g-secrets.json',JSON.stringify(secrets),{mode:0o600})
  await browser.close()
}
