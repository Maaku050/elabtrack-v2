// Stage A acceptance: real Chromium, existing HTTP routes and isolated PostgreSQL.
import fs from 'node:fs/promises'
import path from 'node:path'
import assert from 'node:assert/strict'
import { randomUUID } from 'node:crypto'
import { browser } from './browser-cdp.mjs'
const base='http://localhost:15175',api='http://localhost:18085/api/v1'
const output=path.resolve('../docs/ux/verification/reconstruction-stage-a')
const fixtures=JSON.parse(await fs.readFile('/tmp/elabtrack-batch1-fixtures.json','utf8'))
const checks=[],screens=[],errors=[],documents=[],requests=[]
const b=await browser(), wait=ms=>new Promise(r=>setTimeout(r,ms))
let access, intercept, paused=[]
b.on(m=>{
 if(m.method==='Runtime.exceptionThrown')errors.push(m.params.exceptionDetails.text)
 if(m.method==='Network.requestWillBeSent') { if(m.params.type==='Document')documents.push(m.params.request.url); if(m.params.request.url.startsWith(api))requests.push({method:m.params.request.method,path:new URL(m.params.request.url).pathname,query:new URL(m.params.request.url).search}) }
 if(m.method==='Fetch.requestPaused') {
  const p=m.params
  if(intercept==='hold') paused.push(p.requestId)
  else if(intercept==='error') void b.command('Fetch.fulfillRequest',{requestId:p.requestId,responseCode:503,responseHeaders:[{name:'Content-Type',value:'application/json'},{name:'Access-Control-Allow-Origin',value:base},{name:'Access-Control-Allow-Credentials',value:'true'}],body:Buffer.from(JSON.stringify({success:false,error:{code:'SERVICE_UNAVAILABLE',message:'Test network failure'},meta:{request_id:'stage-a-test'}})).toString('base64')})
  else void b.command('Fetch.continueRequest',{requestId:p.requestId})
 }
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
  try { if (await b.evaluate(expression)) return } catch (error) { if (!/context|Cannot find|Inspected target/.test(error.message)) throw error }
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
 await b.evaluate(`(async()=>{const loaded=performance.getEntriesByType('resource').find(e=>new URL(e.name).pathname==='/src/app/router.tsx');if(!loaded)throw Error('Loaded router module unavailable');const {router}=await import(loaded.name);await router.navigate(${JSON.stringify(target)})})()`)
 await until(`location.pathname===${JSON.stringify(new URL(target,base).pathname)}&&(${ready})`); if(intercept!=='hold') await settled()
}
async function view(width=1440,height=960) {await b.command('Emulation.setDeviceMetricsOverride',{width,height,deviceScaleFactor:1,mobile:false});await wait(200)}
async function theme(value) {await b.evaluate(`(async()=>{const loaded=performance.getEntriesByType('resource').find(e=>new URL(e.name).pathname==='/src/stores/ui-store.ts');const {useUIStore}=await import(loaded.name);useUIStore.getState().setTheme(${JSON.stringify(value)})})()`)}
async function screenshot(name,full=false) {
 await b.evaluate('document.fonts.ready.then(()=>true)');await wait(200)
 assert.ok(await b.evaluate('document.documentElement.scrollWidth<=innerWidth+1'),name+': viewport overflow')
 const opts={format:'png',captureBeyondViewport:full}
 if(full) {await b.evaluate('window.scrollTo(0,document.documentElement.scrollHeight)');await wait(150)}
 opts.captureBeyondViewport=false
 const shot=await b.command('Page.captureScreenshot',opts);await fs.writeFile(path.join(output,name+'.png'),Buffer.from(shot.data,'base64'));screens.push(name+'.png');if(full)await b.evaluate('window.scrollTo(0,0)')
}
async function uiLogin(role) {
 if(await b.evaluate("!!document.querySelector('button[aria-label=\"Sign Out\"]')")) {await click('button[aria-label="Sign Out"]');await until("location.pathname==='/login'")}
 if(await b.evaluate("location.pathname!=='/login'")) await b.command('Page.navigate',{url:base+'/login'});await until("!!document.querySelector('#login-email')&&!document.querySelector('#login-email').disabled")
 await fill('#login-email',fixtures[role].email);await fill('#login-password',fixtures[role].password);await click('button[type=submit]');await until(role==='borrower'?"location.pathname==='/borrower/home'":"!!document.querySelector('#staff-content h1')");if(intercept!=='hold') await settled()
}
try {
 await fs.mkdir(output,{recursive:true});await b.command('Network.enable');await view()
 access=await loginToken('admin')
 let items
 try {items=JSON.parse(await fs.readFile('/tmp/elabtrack-stage-a-items.json','utf8'))} catch {
  const category=await ok('/equipment-categories','POST',{name:'TEST Culinary tools',is_active:true,expected_version:0})
  // Exercise unknown total category navigation without inventing a count in the UI.
  for(let i=0;i<100;i++)await ok('/equipment-categories','POST',{name:`TEST Category ${String(i).padStart(3,'0')}`,is_active:true,expected_version:0})
  const borrowers=[],equipment=[]
  const names=['Mixing bowl','Chef knife','Cutting board','Measuring jug','Saucepan','Whisk','Rolling pin','Colander','Tongs','Ladle','Spatula','Stockpot']
  for(let i=0;i<32;i++) {
   const student=i%4!==0
   borrowers.push(await ok('/borrowers','POST',{name:`TEST ${student?'Student':'Faculty'} ${String(i).padStart(2,'0')}`,email:`stage-a-${randomUUID()}@${student?'students.example.invalid':'faculty.example.invalid'}`,borrower_type:student?'STUDENT':'FACULTY',...(student?{student_id:'00'+String(28366+i),course:'TEST Culinary Arts'}:{})}))
   equipment.push(await ok('/equipment','POST',{name:`TEST ${names[i%names.length]} ${String(i).padStart(2,'0')}`,description:'Disposable Stage A acceptance fixture',category_id:category.id,opening_quantity:i%12,reason:'Synthetic opening stock',expected_version:0}))
  }
  const inactive=borrowers[0];await ok('/borrowers/'+inactive.id+'/status','PATCH',{active:false,confirm:true,expected_updated_at:inactive.updated_at})
  const png=await fs.readFile('/tmp/elabtrack-batch1-catalog.png'),form=new FormData();form.append('image',new Blob([png],{type:'image/png'}),'test.png');form.append('expected_version',String(equipment[0].metadata_version))
  const upload=await fetch(api+'/equipment/'+equipment[0].id+'/image',{method:'POST',headers:{Origin:base,Authorization:'Bearer '+access,'Idempotency-Key':randomUUID()},body:form});assert.equal(upload.status,200)
  items={category,borrowers,equipment};await fs.writeFile('/tmp/elabtrack-stage-a-items.json',JSON.stringify(items),{mode:0o600})
 }
 pass('Real PostgreSQL fixtures: 32 Student/Faculty borrowers, 32 equipment, 101 categories; activation remains UNCONFIGURED')
 await uiLogin('admin');pass('Admin Phase 4A login and existing protected Dashboard')
 await route('/staff/borrowers',"!!document.querySelector('.directory-table-region tbody')")
 const documentCount=documents.length
 await b.evaluate("window.stageAChrome={sidebar:document.querySelector('.staff-sidebar'),header:document.querySelector('.staff-top-bar'),main:document.querySelector('#staff-content')};window.stageASidebarEvents=[];window.stageALastSidebar=stageAChrome.sidebar;new MutationObserver(()=>{const current=document.querySelector('.staff-sidebar');if(current!==window.stageALastSidebar){stageASidebarEvents.push({path:location.href,width:innerWidth,mobile:matchMedia('(max-width:767px)').matches,collapsed:document.querySelector('.fsmo-navigation-body')?.dataset.collapsed});window.stageALastSidebar=current}}).observe(document.querySelector('#root'),{subtree:true,childList:true});true")
 assert.equal(await b.evaluate("document.querySelectorAll('.directory-table-region tbody tr').length"),25)
 await screenshot('borrowers-expanded-light');await theme('dark');await screenshot('borrowers-expanded-dark');await theme('light')
 const expectedSecond=await ok('/borrowers?page=2&per_page=25');await click('[aria-label="Go to next page"]');await until(`new URLSearchParams(location.search).get('page')==='2'&&document.querySelector('.directory-record-name')?.textContent===${JSON.stringify(expectedSecond.items[0].name)}`);await settled()
 assert.equal(await b.evaluate("document.querySelectorAll('.directory-table-region tbody tr').length"),expectedSecond.items.length)
 await screenshot('borrowers-pagination-page2-light',true)
 await b.evaluate('history.back();true');await until("!new URLSearchParams(location.search).has('page')&&document.querySelectorAll('.directory-table-region tbody tr').length===25");await settled()
 pass('Borrower server pagination and browser Back restore URL/view state')
 const facultyPage=await ok('/borrowers?page=1&per_page=25&borrower_type=FACULTY');await select('select[aria-label="Borrower type"]','FACULTY');await settled();await until(`document.querySelectorAll('.directory-type[data-type=FACULTY]').length===${facultyPage.items.length}`)
 assert.equal(await b.evaluate("document.querySelectorAll('.directory-type[data-type=STUDENT]').length"),0)
 const inactivePage=await ok('/borrowers?page=1&per_page=25&borrower_type=FACULTY&status=INACTIVE');await select('select[aria-label="Status"]','INACTIVE');await settled();await until(`document.querySelectorAll('.directory-table-region tbody tr').length===${inactivePage.items.length}`)
 pass('Student/Faculty and inactive filters use real server results')
 await route('/staff/borrowers?search=NO-MATCH-STAGE-A',"document.body.textContent.includes('No matching records')");await screenshot('borrowers-empty-light');assert.equal(await b.evaluate("document.querySelector('[aria-label=\"Go to next page\"]').getAttribute('aria-disabled')"),'true')
 pass('Zero-result empty state and disabled single-page navigation')
 await route('/staff/inventory',"!!document.querySelector('.directory-table-region tbody')")
 await screenshot('inventory-expanded-light');await theme('dark');await screenshot('inventory-expanded-dark');await theme('light')
 const totals=(await ok('/equipment?page=1&per_page=25&sort=name')).totals
 const displayed=await b.evaluate("[...document.querySelectorAll('.directory-stat strong')].map(e=>Number(e.textContent.replaceAll(',','')))")
 assert.deepEqual(displayed,[totals.available,totals.reserved,totals.checked_out,totals.damaged_held])
 await b.evaluate(`(()=>{const e=[...document.querySelectorAll('.directory-equipment')].find(e=>e.textContent.includes(${JSON.stringify(items.equipment[0].name)}));if(!e)throw Error('Image fixture row missing');e.scrollIntoView({block:'center'});return true})()`);await until("!!document.querySelector('.directory-equipment img')");await screenshot('inventory-catalog-image-light');await b.evaluate('window.scrollTo(0,0)')
 pass('All four server stock buckets, result-scoped metrics and actor-protected catalog image')
 await select('select[aria-label="Sort"]','available');await settled();await until("new URLSearchParams(location.search).get('sort')==='available'")
 const backendSorted=await ok('/equipment?page=1&per_page=25&sort=available')
 assert.deepEqual(await b.evaluate("[...document.querySelectorAll('.directory-record-name')].map(e=>e.textContent)"),backendSorted.items.map(e=>e.name))
 await click('.directory-availability input');await settled();assert.ok(await b.evaluate("[...document.querySelectorAll('.directory-available-number')].every(e=>Number(e.textContent)>0)"))
 pass('Inventory sort and available-only filters preserve backend ordering/counts')
 await route('/staff/inventory',"!!document.querySelector('.directory-table-region tbody')")
 await click('[aria-label="Categories pagination"] [aria-label="Go to next page"]');await until("new URLSearchParams(location.search).get('category_page')==='2'&&document.querySelector('.directory-category-pages [role=status]')?.textContent.includes('Page 2')&&document.querySelector('[aria-label=\"Categories pagination\"] [aria-label=\"Go to next page\"]')?.getAttribute('aria-disabled')==='true'");await settled()
 assert.equal(await b.evaluate("document.querySelector('[aria-label=\"Categories pagination\"] [aria-label=\"Go to next page\"]').getAttribute('aria-disabled')"),'true')
 pass('Unknown-total category pages stay honest and stop at short final batch')
 await route('/staff/inventory',"!!document.querySelector('.directory-table-region tbody')")
 await click('[aria-label="Toggle navigation"]');await until("document.querySelector('.fsmo-navigation-body').dataset.collapsed==='true'");await wait(250)
 const center=await b.evaluate("(()=>{const a=document.querySelector('.fsmo-brand-link img').getBoundingClientRect(),s=document.querySelector('.staff-sidebar').getBoundingClientRect();return Math.abs(a.x+a.width/2-(s.x+s.width/2))})()")
 assert.ok(center<1,'FSMO seal centered in icon rail')
 await screenshot('inventory-collapsed-light');await theme('dark');await screenshot('inventory-collapsed-dark');await theme('light')
 await b.evaluate("document.querySelector('.fsmo-navigation-link[aria-label=Inventory]').focus()");await b.command('Input.dispatchKeyEvent',{type:'keyDown',key:'Tab',code:'Tab',windowsVirtualKeyCode:9});await b.command('Input.dispatchKeyEvent',{type:'keyUp',key:'Tab',code:'Tab',windowsVirtualKeyCode:9})
 await until("document.activeElement.getAttribute('aria-label')==='Borrowers'&&[...document.querySelectorAll('[data-slot=tooltip-content][data-open]')].some(e=>e.textContent==='Borrowers'&&e.getClientRects().length)")
 await screenshot('collapsed-tooltip-light')
 await b.command('Input.dispatchKeyEvent',{type:'keyDown',key:'Escape',code:'Escape',windowsVirtualKeyCode:27});await wait(100)
 await route('/staff/dashboard',"document.querySelector('#staff-content h1')?.textContent==='Dashboard'")
 assert.equal(await b.evaluate("document.querySelector('.fsmo-navigation-body').dataset.collapsed"),'true')
 assert.ok(await b.evaluate("stageAChrome.sidebar===document.querySelector('.staff-sidebar')&&stageAChrome.header===document.querySelector('.staff-top-bar')&&stageAChrome.main===document.querySelector('#staff-content')"))
 assert.equal(documents.length,documentCount)
 pass('Persistent shell/DOM, centered collapsed seal, keyboard tooltip and no SPA document reload')
 await b.command('Input.dispatchKeyEvent',{type:'keyDown',key:'b',code:'KeyB',modifiers:2,windowsVirtualKeyCode:66});await b.command('Input.dispatchKeyEvent',{type:'keyUp',key:'b',code:'KeyB',modifiers:2,windowsVirtualKeyCode:66})
 await until("document.querySelector('.fsmo-navigation-body').dataset.collapsed==='false'");pass('Ctrl+B uses the same SidebarProvider/Zustand collapse state')
 await route('/staff/borrowers',"!!document.querySelector('.directory-table-region tbody')")
 await click('[aria-label="Account actions"]');await until("!!document.querySelector('[role=menu]')");await screenshot('account-footer-menu-light');await b.command('Input.dispatchKeyEvent',{type:'keyDown',key:'Escape',code:'Escape',windowsVirtualKeyCode:27});await wait(200)
 assert.equal(await b.evaluate('document.activeElement.getAttribute("aria-label")'),'Account actions')
 pass('Base UI footer menu supports account link, Sign Out and Escape focus return')
 for(const width of [390,768,1024,1366]) {
  await view(width);assert.ok(await b.evaluate("[...document.querySelectorAll('.directory-actions .app-button')].every(e=>e.scrollWidth<=e.clientWidth+1)"),'Action text fits at '+width);await screenshot('borrowers-responsive-'+width+'-light');await theme('dark');await screenshot('borrowers-responsive-'+width+'-dark');await theme('light')
 }
 await view(390);await click('[aria-label="Toggle navigation"]');await until("!!document.querySelector('[role=dialog]')")
 await screenshot('mobile-navigation-light');await theme('dark');await screenshot('mobile-navigation-dark');await theme('light')
 await b.command('Input.dispatchKeyEvent',{type:'keyDown',key:'Escape',code:'Escape',windowsVirtualKeyCode:27});await until("!document.querySelector('[role=dialog]')")
 assert.equal(await b.evaluate('document.activeElement.getAttribute("aria-label")'),'Toggle navigation')
 await click('[aria-label="Toggle navigation"]');await until("!!document.querySelector('[role=dialog]')");await click('[role=dialog] a[aria-label=Inventory]');await until("location.pathname==='/staff/inventory'&&!document.querySelector('[role=dialog]')");await settled()
 await screenshot('inventory-mobile-light');await theme('dark');await screenshot('inventory-mobile-dark');await theme('light')
 pass('390/768/1024/1366 responsive themes, mobile focus return, navigation closes drawer, no document overflow')
 await view();await route('/staff/inventory/'+items.equipment[0].id,"document.querySelector('#staff-content h1')?.textContent==='Equipment Details'&&!!document.querySelector('.inventory-detail-image')")
 await route('/staff/borrowers/'+items.borrowers[1].id,"document.body.textContent.includes('Account information')")
 pass('Existing equipment and borrower detail links remain functional')
 // Network failures are controlled browser interceptions; ordinary checks use real HTTP/PG.
 await route('/staff/borrowers',"!!document.querySelector('.directory-table-region tbody')")
 intercept='hold';await b.command('Fetch.enable',{patterns:[{urlPattern:api+'/borrowers*',requestStage:'Request'}]})
 await route('/staff/borrowers?search=LOADING-STAGE-A',"!!document.querySelector('.directory-loading')")
 await screenshot('borrowers-loading-light');intercept=undefined;for(const id of paused.splice(0))await b.command('Fetch.continueRequest',{requestId:id});await b.command('Fetch.disable');await settled()
 intercept='error';await b.command('Fetch.enable',{patterns:[{urlPattern:api+'/borrowers*',requestStage:'Request'}]})
 await route('/staff/borrowers?search=ERROR-STAGE-A',"!!document.querySelector('.directory-results [role=alert]')")
 await screenshot('borrowers-error-light');intercept=undefined;await b.command('Fetch.disable');await button('Retry');await until("document.body.textContent.includes('No matching records')");await settled()
 pass('Accessible loading/error/Retry states under controlled network delay/failure')
 const staffToken=await loginToken('staff')
 assert.equal((await request('/borrowers','POST',{name:'TEST prohibited',email:'prohibited@faculty.example.invalid',borrower_type:'FACULTY'},staffToken)).status,403)
 assert.equal((await request('/equipment/'+items.equipment[0].id+'/adjustments','POST',{kind:'RECONCILE',quantity:20,expected_sequence:items.equipment[0].stock_sequence,reason:'TEST denied',confirm:true},staffToken)).status,403)
 await uiLogin('staff');await route('/staff/borrowers',"!!document.querySelector('.directory-table-region tbody')")
 assert.equal(await b.evaluate("!!document.querySelector('a[href=\"/staff/borrowers/new\"]')"),false)
 assert.equal(await b.evaluate("!!document.querySelector('.fsmo-navigation-link[aria-label=Administration]')"),false)
 await route('/staff/borrowers/new',"document.body.textContent.includes('Access denied')")
 pass('STAFF navigation and controls respect roles; backend creation/correction and direct routes deny')
 await view();await route('/staff/dashboard',"!!document.querySelector('.staff-sidebar')");await click('button[aria-label="Sign Out"]');await until("location.pathname==='/login'&&!!document.querySelector('#login-email')")
 pass('Existing logout clears protected workspace and returns to unchanged immersive login')
 await uiLogin('borrower');await route('/staff/inventory',"document.body.textContent.includes('Access denied')")
 pass('BORROWER cannot enter operational inventory')
 assert.deepEqual(errors,[])
 const evidence={date:new Date().toISOString(),result:'PASS',ownerVisualAcceptance:'PENDING',scope:'Isolated PostgreSQL, real existing API, Chromium; synthetic fixture data, no live Brevo; delay/error checks explicitly intercepted',checks,screens,errors,apiRequests:requests,documentNavigations:documents.length}
 await fs.writeFile(path.join(output,'acceptance.json'),JSON.stringify(evidence,null,2)+'\n')
 console.log(JSON.stringify({result:'PASS',checks:checks.length,screenshots:screens.length}))
} catch(error) {await fs.writeFile('/tmp/elabtrack-stage-a-browser-diagnostic.json',JSON.stringify(await b.evaluate("({active:document.activeElement?.outerHTML,tooltip:[...document.querySelectorAll('[data-slot=tooltip-content]')].map(e=>e.outerHTML),collapsed:document.querySelector('.fsmo-navigation-body')?.dataset.collapsed,identity:window.stageAChrome?{sidebar:stageAChrome.sidebar===document.querySelector('.staff-sidebar'),header:stageAChrome.header===document.querySelector('.staff-top-bar'),main:stageAChrome.main===document.querySelector('#staff-content')}:null,href:location.href,sidebarEvents:window.stageASidebarEvents})")));await fs.writeFile('/tmp/elabtrack-stage-a-browser-documents.json',JSON.stringify(documents));await fs.writeFile('/tmp/elabtrack-stage-a-browser-failure.txt',String(error));await screenshot('failure').catch(()=>{});throw error} finally {await b.close()}
