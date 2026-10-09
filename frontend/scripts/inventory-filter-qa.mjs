// Focused inventory-filter acceptance against the owned disposable API/PostgreSQL harness.
import fs from 'node:fs/promises'
import path from 'node:path'
import assert from 'node:assert/strict'
import { randomUUID } from 'node:crypto'
import { browser } from './browser-cdp.mjs'
const base='http://localhost:15175',api='http://localhost:18085/api/v1'
const output=path.resolve('../docs/ux/verification/inventory-filter')
const fixtures=JSON.parse(await fs.readFile('/tmp/elabtrack-batch1-fixtures.json','utf8'))
const checks=[],screens=[],errors=[],documents=[],requests=[]
const b=await browser(), wait=ms=>new Promise(r=>setTimeout(r,ms))
// Retain awaited page promises while CDP waits; Chromium may collect an
// otherwise unreferenced dynamic-import promise during long screenshot runs.
const evaluate=b.evaluate
b.evaluate=expression=>evaluate(`globalThis.__stageBEvaluation=(async()=>(${expression}))();globalThis.__stageBEvaluation`)
let access
b.on(m=>{
 if(m.method==='Runtime.exceptionThrown')errors.push(m.params.exceptionDetails.text)
 if(m.method==='Network.requestWillBeSent') { if(m.params.type==='Document')documents.push(m.params.request.url); if(m.params.request.url.startsWith(api))requests.push({method:m.params.request.method,path:new URL(m.params.request.url).pathname,query:new URL(m.params.request.url).search}) }

})
async function request(route,method='GET',data,token=access) {
 const r=await fetch(api+route,{method,headers:{Origin:base,Authorization:'Bearer '+token,'Content-Type':'application/json','Idempotency-Key':randomUUID()},body:data===undefined?undefined:JSON.stringify(data)})
 const body=await r.json();return {status:r.status,data:body.data,code:body.error?.code}
}
async function ok(route,method,data) {const r=await request(route,method,data);assert.ok(r.status<300,route+': '+r.status+'/'+r.code);return r.data}
async function loginToken(role) {const r=await fetch(api+'/auth/login',{method:'POST',headers:{Origin:base,'Content-Type':'application/json'},body:JSON.stringify({email:fixtures[role].email,password:fixtures[role].password})});assert.equal(r.status,200);return (await r.json()).data.access_token}
function pass(name){checks.push(name);console.log('PASS',name)}
async function until(expression) {
 for (let i = 0; i < 150; i++) {
  try { if (await b.evaluate(expression)) return } catch (error) { if (!/context|Cannot find|Inspected target|Promise was collected/.test(error.message)) throw error }
  await wait(100)
 }
 throw Error('Timed out: ' + expression)
}
async function settled() { await until("(async()=>{const {queryClient}=await import(performance.getEntriesByType('resource').find(e=>new URL(e.name).pathname==='/src/app/query-client.ts').name);return queryClient.isFetching()===0})()"); await wait(150) }
async function click(selector) {
 await b.evaluate(`(async()=>{const e=document.querySelector(${JSON.stringify(selector)});if(!e)throw Error('Missing control');e.scrollIntoView({block:'center'});await new Promise(resolve=>requestAnimationFrame(()=>requestAnimationFrame(resolve)))})()`)
 await until(`(()=>{const e=document.querySelector(${JSON.stringify(selector)});if(!e)return false;const r=e.getBoundingClientRect(),hit=document.elementFromPoint(r.x+r.width/2,r.y+r.height/2);return hit===e||e.contains(hit)})()`)
 const point = await b.evaluate(`(()=>{const e=document.querySelector(${JSON.stringify(selector)});if(!e)throw Error('Missing control');e.scrollIntoView({block:'center'});const r=e.getBoundingClientRect();return{x:r.x+r.width/2,y:r.y+r.height/2}})()`)
 for (const type of ['mousePressed', 'mouseReleased']) await b.command('Input.dispatchMouseEvent', { type, ...point, button: 'left', clickCount: 1 })
}
async function button(label) {
 await until(`[...document.querySelectorAll('button')].some(e=>e.textContent.trim()===${JSON.stringify(label)}&&!e.disabled&&e.getClientRects().length>0)`)
 const index = await b.evaluate(`[...document.querySelectorAll('button')].findIndex(e=>e.textContent.trim()===${JSON.stringify(label)}&&e.getClientRects().length>0)`)
 assert.ok(index >= 0, 'Missing button: ' + label)
 await b.evaluate(`document.querySelectorAll('button')[${index}].setAttribute('data-qa-click','true')`)
 await click('[data-qa-click]'); await b.evaluate("document.querySelector('[data-qa-click]')?.removeAttribute('data-qa-click')")
}
async function fill(selector, value) {
 await click(selector)
 for (const type of ['keyDown', 'keyUp']) await b.command('Input.dispatchKeyEvent', { type, key: 'a', code: 'KeyA', windowsVirtualKeyCode: 65, modifiers: 2 })
 await b.command('Input.insertText', { text: value })
}
async function select(selector, value) { await b.evaluate(`(()=>{const e=document.querySelector(${JSON.stringify(selector)});e.value=${JSON.stringify(value)};e.dispatchEvent(new Event('change',{bubbles:true}))})()`) }
async function route(target, ready) {
 await b.evaluate(`(async()=>{const loaded=performance.getEntriesByType('resource').find(e=>new URL(e.name).pathname==='/src/app/router.tsx');if(!loaded)throw Error('Loaded router module unavailable');const {router}=await import(loaded.name);void router.navigate(${JSON.stringify(target)});return true})()`)
 await until(`location.pathname===${JSON.stringify(new URL(target,base).pathname)}&&(${ready})`); await settled()
}
async function view(width=1440,height=960) {await b.command('Emulation.setDeviceMetricsOverride',{width,height,deviceScaleFactor:1,mobile:false});await wait(200)}
async function theme(value) {await b.evaluate(`(async()=>{const loaded=performance.getEntriesByType('resource').find(e=>new URL(e.name).pathname==='/src/stores/ui-store.ts');const {useUIStore}=await import(loaded.name);useUIStore.getState().setTheme(${JSON.stringify(value)})})()`)}
async function screenshot(name,full=false) {
 await b.evaluate('(()=>{globalThis.__stageBFonts=document.fonts.ready.then(()=>true);return globalThis.__stageBFonts})()');await wait(200)
 assert.ok(await b.evaluate('document.documentElement.scrollWidth<=innerWidth+1'),name+': viewport overflow')
 const opts={format:'png',captureBeyondViewport:full}
 if(full) {await b.evaluate('window.scrollTo(0,document.documentElement.scrollHeight)');await wait(150)}
 opts.captureBeyondViewport=false
 const shot=await b.command('Page.captureScreenshot',opts);await fs.writeFile(path.join(output,name+'.png'),Buffer.from(shot.data,'base64'));screens.push(name+'.png');if(full)await b.evaluate('window.scrollTo(0,0)')
}
async function uiLogin(role) {
 if(await b.evaluate("!!document.querySelector('button[aria-label=\"Sign Out\"]')")) {await click('button[aria-label="Sign Out"]');await until("location.pathname==='/login'")}
 if(await b.evaluate("location.pathname!=='/login'")) await b.command('Page.navigate',{url:base+'/login'});await until("!!document.querySelector('#login-email')&&!document.querySelector('#login-email').disabled")
 await fill('#login-email',fixtures[role].email);await fill('#login-password',fixtures[role].password);await click('button[type=submit]');await until(role==='borrower'?"location.pathname.startsWith('/borrower/')":"!!document.querySelector('#staff-content h1')");await settled()
}

const checkbox='.directory-availability input'
const prefix='TEST Borrowability '+randomUUID().slice(0,8)
const seed=[]
async function count(total, rows) {
 await until(`document.querySelector('.directory-panel-heading p')?.textContent===${JSON.stringify(total+' matching equipment types')}`)
 await settled()
 assert.equal(await b.evaluate("document.querySelectorAll('.directory-inventory-table tbody tr').length"),rows)
}
try {
 await fs.mkdir(output,{recursive:true});await b.command('Network.enable');await view();access=await loginToken('admin')
 const category=await ok('/equipment-categories','POST',{name:prefix+' utensils',is_active:true,expected_version:0})
 const other=await ok('/equipment-categories','POST',{name:prefix+' cookware',is_active:true,expected_version:0})
 async function create(name,status,quantity,cat) {
  let v=await ok('/equipment','POST',{name:prefix+' '+name,description:'Synthetic isolated inventory-filter fixture',category_id:cat,opening_quantity:quantity,reason:'Synthetic opening count',expected_version:0})
  if(status!=='ACTIVE')v=await ok('/equipment/'+v.id+'/status','PATCH',{status,expected_version:v.metadata_version,confirm:true})
  seed.push(v);return v
 }
 await create('Spoon active positive','ACTIVE',20,category.id)
 await create('Spoon active zero','ACTIVE',0,category.id)
 await create('Spoon inactive positive','INACTIVE',20,category.id)
 await create('Spoon inactive zero','INACTIVE',0,category.id)
 await create('Spoon archived positive','ARCHIVED',20,category.id)
 for(let i=0;i<26;i++)await create('Utensil '+String(i).padStart(2,'0'),'ACTIVE',1,other.id)
 await ok('/equipment-categories/'+category.id,'PATCH',{name:category.name,is_active:false,expected_version:category.version})
 const search='search='+encodeURIComponent(prefix)
 for(const role of ['admin','staff','borrower']) {
  const token=await loginToken(role)
  const p=await request('/equipment?'+search+'&available_only=true&per_page=100','GET',undefined,token)
  assert.equal(p.status,200);assert.equal(p.data.total,27);assert.equal(p.data.totals.available,46);assert.equal(p.data.items.length,27)
  assert.ok(p.data.items.every(v=>v.status==='ACTIVE'&&v.stock.available>0))
 }
 pass('Real API excludes inactive/archived and zero stock for Admin, Staff and Borrower; matching total 27 and stock total 46')
 await uiLogin('admin');await route('/staff/inventory?'+search,"!!document.querySelector('.directory-inventory-table')");await count(31,25)
 await click(checkbox);await count(27,25)
 assert.equal(await b.evaluate("document.querySelector('.directory-availability').textContent.trim()"),'Available for borrowing')
 assert.ok(await b.evaluate("[...document.querySelectorAll('tbody tr')].every(e=>!e.textContent.includes('Inactive')&&!e.textContent.includes('Archived')&&!e.textContent.includes('active zero'))"))
 const page1=await b.evaluate("[...document.querySelectorAll('tbody .directory-record-name')].map(e=>e.textContent)")
 await click('[aria-label="Go to next page"]');await until("new URLSearchParams(location.search).get('page')==='2'");await count(27,2)
 const page2=await b.evaluate("[...document.querySelectorAll('tbody .directory-record-name')].map(e=>e.textContent)")
 assert.equal(new Set([...page1,...page2]).size,27)
 assert.ok(requests.some(r=>r.path==='/api/v1/equipment'&&new URLSearchParams(r.query).get('page')==='2'&&new URLSearchParams(r.query).get('available_only')==='true'))
 await b.command('Page.reload');await until("!!document.querySelector('.directory-availability')");await count(27,2)
 assert.equal(await b.evaluate("new URLSearchParams(location.search).get('page')"),'2')
 assert.deepEqual(await b.evaluate("[...document.querySelectorAll('tbody .directory-record-name')].map(e=>e.textContent)"),page2)
 pass('Real server pagination returns 25 then 2 unique eligible records with total 27; no page-local filter')
 const docsBefore=documents.length
 await select('select[aria-label="Status"]','INACTIVE');await count(0,0)
 assert.equal(await b.evaluate("new URLSearchParams(location.search).get('page')"),null)
 assert.equal(await b.evaluate('document.querySelector('+JSON.stringify(checkbox)+').checked'),true)
 await screenshot('inactive-availability-empty-light');await theme('dark');await screenshot('inactive-availability-empty-dark');await theme('light')
 await click(checkbox);await count(2,2)
 assert.equal(await b.evaluate("document.querySelector('select[aria-label=\"Status\"]').value"),'INACTIVE')
 assert.ok(await b.evaluate("[...document.querySelectorAll('tbody tr')].some(r=>r.textContent.includes('inactive positive')&&r.querySelector('.directory-available-number').textContent==='20')"))
 await select('select[aria-label="Status"]','ARCHIVED');await count(1,1);await click(checkbox);await count(0,0)
 await select('select[aria-label="Status"]','ACTIVE');await count(27,25)
 pass('Status intersections, archived visibility, retained inactive physical stock and reset-to-first-page behave predictably')
 await select('select[aria-label="Category"]',category.id);await count(1,1)
 await fill('input[aria-label="Search equipment"]',prefix+' Spoon inactive');await count(0,0)
 await fill('input[aria-label="Search equipment"]',prefix+' Spoon');await count(1,1)
 await select('select[aria-label="Sort"]','available');await count(1,1)
 pass('Category, literal search, sorting and availability intersect; inactive category does not hide existing active equipment')
 const restoreURL=await b.evaluate('location.href')
 assert.equal(documents.length,docsBefore,'Ordinary filter navigation stays in the persistent SPA')
 await b.command('Page.reload');await until("!!document.querySelector('.directory-availability')");await count(1,1)
 assert.equal(await b.evaluate('location.href'),restoreURL)
 assert.equal(await b.evaluate('document.querySelector('+JSON.stringify(checkbox)+').checked'),true)
 assert.equal(await b.evaluate("document.querySelector('select[aria-label=\"Category\"]').value"),category.id)
 assert.equal(await b.evaluate("document.querySelector('select[aria-label=\"Status\"]').value"),'ACTIVE')
 assert.equal(await b.evaluate("document.querySelector('select[aria-label=\"Sort\"]').value"),'available')
 pass('Full document refresh restores all applied filters and authenticated server results')
 for(const width of [320,390,1366])for(const color of ['light','dark']) {
  await view(width,960);await theme(color);await b.evaluate('window.scrollTo(0,0)');await screenshot('availability-'+width+'-'+color)
 }
 await view();await theme('light');await button('Clear filters');await settled()
 assert.equal(await b.evaluate('location.search'),'')
 assert.equal(await b.evaluate('document.querySelector('+JSON.stringify(checkbox)+').checked'),false)
 assert.equal(await b.evaluate("document.querySelector('select[aria-label=\"Status\"]').value"),'')
 assert.equal(await b.evaluate("document.querySelector('select[aria-label=\"Category\"]').value"),'')
 pass('Clear filters removes search/category/status/availability/sort and resets pagination')
 for(const v of seed) {
  const current=await ok('/equipment/'+v.id)
  assert.deepEqual(current.stock,v.stock);assert.equal(current.status,v.status);assert.equal(current.stock_sequence,v.stock_sequence);assert.equal(current.metadata_version,v.metadata_version)
 }
 pass('Filtering and refresh do not mutate stock buckets, status, sequence or metadata versions')
 await uiLogin('staff');await route('/staff/inventory?'+search+'&status=INACTIVE&available_only=true',"!!document.querySelector('.directory-availability')");await count(0,0)
 assert.equal(await b.evaluate("document.querySelector('select[aria-label=\"Status\"]').value"),'INACTIVE')
 pass('Staff uses the same independent status/availability controls and empty intersection')
 const borrowerToken=await loginToken('borrower'),terms=await ok('/terms/current')
 assert.ok(terms.title.startsWith('TEST ONLY'),'Only synthetic isolated terms may be accepted by the fixture')
 const consent=await request('/terms/'+terms.id+'/accept','POST',{},borrowerToken);assert.equal(consent.status,200)
 await uiLogin('borrower');await route('/borrower/equipment?'+search+'&available_only=true',"!!document.querySelector('.inventory-catalog')")
 assert.equal(await b.evaluate("document.querySelector('.management-checkbox').textContent.trim()"),'Available for borrowing')
 assert.equal(await b.evaluate("document.querySelectorAll('.inventory-catalog-card').length"),25)
 assert.ok(await b.evaluate("[...document.querySelectorAll('.inventory-catalog-card')].every(e=>!e.textContent.includes('inactive')&&!e.textContent.includes('archived')&&!e.textContent.includes('active zero'))"))
 await click('[aria-label="Go to next page"]');await until("new URLSearchParams(location.search).get('page')==='2'");await settled()
 assert.equal(await b.evaluate("document.querySelectorAll('.inventory-catalog-card').length"),2)
 await screenshot('borrower-availability-light');await theme('dark');await screenshot('borrower-availability-dark')
 pass('Borrower catalog uses the same terminology and 25/2 server pages under the existing terms guard')
 assert.deepEqual(errors,[])
 await fs.writeFile(path.join(output,'acceptance.json'),JSON.stringify({result:'PASS',checks,screens,errors,apiRequests:requests,documentNavigations:documents.length,fixtureSummary:{totalEquipment:31,eligibleEquipment:27,eligibleAvailable:46}},null,2)+'\n')
 console.log(JSON.stringify({result:'PASS',checks:checks.length,screens:screens.length}))
} catch(error) {await screenshot('failure').catch(()=>{});throw error} finally {await b.close()}
