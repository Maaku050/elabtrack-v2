// Real Chromium + explicitly isolated PostgreSQL/HTTP. Never owner credentials.
import fs from 'node:fs/promises'
import path from 'node:path'
import assert from 'node:assert/strict'
import { execFileSync } from 'node:child_process'
import { randomUUID } from 'node:crypto'
import { browser } from './browser-cdp.mjs'

const mode = process.argv.includes('--before') ? 'before' : 'after'
const checksOnly = process.argv.includes('--checks-only')
const base = 'http://localhost:15175', api = 'http://localhost:18085/api/v1'
const output = path.resolve('../docs/ux/verification/phase56-refinement', mode)
const fixtures = JSON.parse(await fs.readFile('/tmp/elabtrack-batch1-fixtures.json', 'utf8'))
const imageFile = '/tmp/elabtrack-batch1-catalog.png'
const png = await fs.readFile(imageFile)
const loginResponse = await fetch(api + '/auth/login', { method: 'POST', headers: { Origin: base, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: fixtures.admin.email, password: fixtures.admin.password }) })
assert.equal(loginResponse.status, 200)
const access = (await loginResponse.json()).data.access_token
const checks = [], screens = [], errors = [], network = []
const b = await browser()
const wait = ms => new Promise(resolve => setTimeout(resolve, ms))
b.on(message => {
 if (message.method === 'Runtime.exceptionThrown') errors.push(message.params.exceptionDetails.text)
 if (message.method === 'Network.requestWillBeSent' && message.params.request.url.startsWith(api)) network.push({ method: message.params.request.method, path: new URL(message.params.request.url).pathname })
})
async function request(route, method = 'GET', data, key = randomUUID()) {
 const headers = { Origin: base, Authorization: 'Bearer ' + access, 'Idempotency-Key': key }
 if (!(data instanceof FormData)) headers['Content-Type'] = 'application/json'
 const response = await fetch(api + route, { method, headers, body: data instanceof FormData ? data : data === undefined ? undefined : JSON.stringify(data) })
 const body = await response.json()
 return { status: response.status, data: body.data, code: body.error?.code, message: body.error?.message }
}
async function ok(route, method, data) { const result = await request(route, method, data); assert.ok(result.status < 300, `Fixture API ${route}: ${result.status}/${result.code}`); return result.data }
async function until(expression) {
 for (let i = 0; i < 150; i++) {
  try { if (await b.evaluate(expression)) return } catch (error) { if (!/context|Cannot find|Inspected target/.test(error.message)) throw error }
  await wait(100)
 }
 throw Error('Timed out: ' + expression)
}
async function settled() { await until("(async()=>{const {queryClient}=await import('/src/app/query-client.ts');return queryClient.isFetching()===0})()"); await wait(150) }
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
 await until(`location.pathname===${JSON.stringify(target)}&&(${ready})`); await settled()
}
async function file(selector, filename) {
 const { root } = await b.command('DOM.getDocument')
 const { nodeId } = await b.command('DOM.querySelector', { nodeId: root.nodeId, selector })
 assert.ok(nodeId, 'Missing file picker')
 await b.command('DOM.setFileInputFiles', { nodeId, files: [filename] })
 await wait(150)
}
async function shot(name, width = 1440, theme = 'light') {
 await b.command('Emulation.setDeviceMetricsOverride', { width, height: 1000, deviceScaleFactor: 1, mobile: false })
 await b.evaluate(`(async()=>{const {useUIStore}=await import('/src/stores/ui-store.ts');useUIStore.getState().setTheme(${JSON.stringify(theme)})})()`)
 await b.evaluate('document.fonts.ready'); await b.evaluate('window.scrollTo(0,0)'); await wait(200)
 assert.ok(await b.evaluate('document.documentElement.scrollWidth<=innerWidth+1'), name + ': viewport overflow')
 const { cssContentSize } = await b.command('Page.getLayoutMetrics')
 const screenshot = await b.command('Page.captureScreenshot', { format: 'png', captureBeyondViewport: true, clip: { x: 0, y: 0, width: cssContentSize.width, height: cssContentSize.height, scale: 1 } })
 const filename = `${name}-${width}-${theme}.png`
 await fs.writeFile(path.join(output, filename), Buffer.from(screenshot.data, 'base64')); screens.push(filename)
}
async function capture(name) {
 if (checksOnly) return
 for (const width of mode === 'before' ? [1440] : [390, 768, 1024, 1366, 1440]) for (const theme of ['light', 'dark']) await shot(name, width, theme)
 await b.command('Emulation.setDeviceMetricsOverride', { width: 1440, height: 1000, deviceScaleFactor: 1, mobile: false })
}
async function stock(id) { return ok('/equipment/' + id) }
function pass(name) { checks.push(name); console.log('PASS', name) }
try {
 await fs.mkdir(output, { recursive: true }); await b.command('Network.enable'); await b.command('Emulation.setDeviceMetricsOverride', { width: 1440, height: 1000, deviceScaleFactor: 1, mobile: false })
 const category = await ok('/equipment-categories', 'POST', { name: 'TEST refinement ' + randomUUID().slice(0, 8), is_active: true, expected_version: 0 })
 let equipment = await ok('/equipment', 'POST', { name: 'TEST culinary equipment ' + randomUUID().slice(0, 8), description: 'Disposable manual acceptance test equipment.', category_id: category.id, expected_version: 0, opening_quantity: 15, reason: 'Verified synthetic opening count' })
 const upload = new FormData(); upload.append('image', new Blob([png], { type: 'image/png' }), 'test-catalog.png'); upload.append('expected_version', String(equipment.metadata_version))
 equipment = await ok('/equipment/' + equipment.id + '/image', 'POST', upload)
 const faculty = await ok('/borrowers', 'POST', { name: 'TEST Faculty Borrower', email: 'faculty-' + randomUUID() + '@students.example.invalid', borrower_type: 'FACULTY' })
 assert.equal(faculty.delivery_status, 'UNCONFIGURED')
 await fs.writeFile('/tmp/elabtrack-phase56-items.json', JSON.stringify({ equipment, category, faculty }), { mode: 0o600 })
 if (mode === 'after') execFileSync('/home/marvin/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.27.1.linux-amd64/bin/go', ['run', '../integration/phase56-workbook-fixtures.go'], { cwd: path.resolve('../backend'), env: { ...process.env, GOCACHE: '/tmp/elabtrack-phase1b-go-cache' }, stdio: 'pipe' })
 await b.command('Page.navigate', { url: base + '/login' })
 await until("!!document.querySelector('#login-email')&&!document.querySelector('#login-email').disabled")
 await fill('#login-email', fixtures.admin.email); await fill('#login-password', fixtures.admin.password); await click('button[type=submit]')
 await until("document.querySelector('#staff-content h1')?.textContent==='Dashboard'"); await settled()
 await route('/staff/borrowers/' + faculty.id, "document.body.textContent.includes('Account information')"); await capture('borrower-details')
 await route('/admin/borrowers/bulk', "document.body.textContent.includes('Student roster')"); await capture('student-bulk')
 await route('/staff/inventory/categories', "!!document.querySelector('#category-name')"); await capture('categories')
 await route('/staff/inventory', "!!document.querySelector('.management-table')"); await capture('inventory')
 await route('/staff/inventory/new', "!!document.querySelector('#equipment-name')"); await capture('equipment-create')
 await route('/staff/inventory/' + equipment.id + '/edit', "!!document.querySelector('#equipment-name')"); await capture('equipment-edit')
 await route('/staff/inventory/' + equipment.id + '/adjust', "!!document.querySelector('#adjust-quantity')")
 await fill('#adjust-quantity', '5'); await fill('#adjust-reason', 'Verified synthetic addition'); await capture('stock-adjustment')
 await button('Review Stock Change'); await until("!![...document.querySelectorAll('button')].find(e=>e.textContent.trim()==='Confirm Stock Change')"); await capture('stock-review')
 await route('/admin/inventory/' + equipment.id + '/reconcile', "!!document.querySelector('#adjust-quantity')")
 await fill('#adjust-quantity', '20'); await fill('#adjust-reason', 'Verified synthetic available count'); await capture('inventory-correction')
 await button('Review Stock Change'); await until("!![...document.querySelectorAll('button')].find(e=>e.textContent.trim()==='Confirm Stock Change')"); await capture('correction-review')

 if (mode === 'after') {
  // The frozen correction review above must commit exactly once, even on a double click.
  const countAdjustments = () => network.filter(n => n.method === 'POST' && n.path.endsWith('/adjustments')).length
  const start = countAdjustments()
  await b.evaluate("(()=>{const e=[...document.querySelectorAll('button')].find(e=>e.textContent.trim()==='Confirm Stock Change');e.click();e.click()})()")
  await until(`location.pathname==='/staff/inventory/${equipment.id}'`);await settled()
  let current=await stock(equipment.id);assert.equal(current.stock.available,20);assert.equal(current.stock.total_tracked,20);assert.equal(countAdjustments()-start,1)
  await until("document.body.textContent.includes('Stock change recorded.')")
  pass('Admin correction increase persists; synchronous double confirmation submits once and shows success')
  await route('/staff/inventory/'+equipment.id+'/adjust', "!!document.querySelector('#adjust-quantity')")
  await select('#adjust-kind','REMOVE');await fill('#adjust-quantity','5');await fill('#adjust-reason','Verified removal in browser')
  assert.ok(await b.evaluate("document.querySelector('.stock-projection').textContent.includes('-5')"));await button('Review Stock Change');await button('Confirm Stock Change');await until(`location.pathname==='/staff/inventory/${equipment.id}'`);await settled()
  current=await stock(equipment.id);assert.equal(current.stock.available,15);pass('Stock removal projection and persisted stock decrease')
  await route('/staff/inventory/'+equipment.id+'/adjust', "!!document.querySelector('#adjust-quantity')")
  await fill('#adjust-quantity','5');await button('Review Stock Change');await until("document.body.textContent.includes('Explain this change.')")
  await fill('#adjust-reason','Verified browser addition');assert.ok(await b.evaluate("document.querySelector('.stock-projection').textContent.includes('+5')"))
  await button('Review Stock Change');await button('Confirm Stock Change');await until(`location.pathname==='/staff/inventory/${equipment.id}'`);await settled();current=await stock(equipment.id);assert.equal(current.stock.available,20)
  pass('Required explanation and stock addition projection/persistence')
  await route('/staff/inventory/'+equipment.id+'/adjust', "!!document.querySelector('#adjust-quantity')");await select('#adjust-kind','REMOVE');await fill('#adjust-quantity','21');await fill('#adjust-reason','Invalid removal');await button('Review Stock Change');await until("document.body.textContent.includes('outside the allowed range')")
  assert.ok(await b.evaluate("!document.body.textContent.includes('Confirm Stock Change')"));await fill('#adjust-quantity','1.5');await button('Review Stock Change');await until("document.body.textContent.includes('whole number')");pass('Invalid and insufficient quantities cannot enter review')
  await select('#adjust-kind','ADD');await fill('#adjust-quantity','5');await fill('#adjust-reason','Stale browser review');await button('Review Stock Change');current=await stock(equipment.id)
  await ok('/equipment/'+equipment.id+'/adjustments','POST',{kind:'ADD',quantity:2,reason:'Concurrent synthetic change',expected_sequence:current.stock_sequence,confirm:true})
  await button('Confirm Stock Change');await until("document.body.textContent.includes('Stock changed after your review')");assert.ok(await b.evaluate("[...document.querySelectorAll('button')].find(e=>e.textContent.trim()==='Confirm Stock Change').disabled"));assert.equal((await stock(equipment.id)).stock.available,22)
  await button('Back to Edit');assert.ok(await b.evaluate("document.querySelector('.stock-projection').textContent.includes('22')"));pass('Concurrent write rejects stale review and forces updated basis')
  await route('/admin/inventory/'+equipment.id+'/reconcile', "!!document.querySelector('#adjust-quantity')");await fill('#adjust-quantity','12');await fill('#adjust-reason','Verified lower available count')
  assert.ok(await b.evaluate("document.querySelector('.stock-projection').textContent.includes('-10')&&document.body.textContent.includes('decreases available physical stock')"));await button('Review Stock Change');await button('Confirm Stock Change');await until(`location.pathname==='/staff/inventory/${equipment.id}'`);await settled();current=await stock(equipment.id);assert.equal(current.stock.available,12);assert.equal(current.stock.total_tracked,12)
  const movements=await ok('/equipment/'+equipment.id+'/movements');assert.ok(movements.some(m=>m.kind==='RECONCILE'&&m.delta.available===-10));pass('Correction decrease warning, persisted counts and ledger movement history')
  // Browser upload + genuine drag/drop for a malformed roster; server errors remain actionable.
  await route('/admin/borrowers/bulk', "!!document.querySelector('.file-dropzone')")
  await file('input[type=file]','/tmp/elabtrack-phase56-header.xlsx');await button('Validate Roster');await until("document.body.textContent.includes('Header does not match')");pass('Browser workbook structure rejection identifies the header location')
  const validBytes=await fs.readFile('/tmp/elabtrack-phase56-valid.xlsx')
  await b.evaluate(`(()=>{const zone=document.querySelector('.file-dropzone'),dt=new DataTransfer();dt.items.add(new File([new Uint8Array(${JSON.stringify([...validBytes])})],'students.xlsx',{type:'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'}));zone.dispatchEvent(new DragEvent('drop',{bubbles:true,cancelable:true,dataTransfer:dt}))})()`)
  await until("document.body.textContent.includes('students.xlsx')&&!document.querySelector('[role=alert]')");await button('Validate Roster');await until("document.body.textContent.includes('1 valid')");await button('Select All Valid Rows');await button('Review & Confirm Creation');await until("!!document.querySelector('[role=alertdialog]')");await button('Confirm Creation');await until("document.body.textContent.includes('Operation results')");assert.ok(await b.evaluate("document.body.textContent.includes('activation pending')"));pass('Actual drag/drop, validation preview, selection and confirmed Student creation')
  await until("!document.querySelector('.bulk-layout select').disabled");await select('.bulk-layout select','DEACTIVATE');await until("document.querySelector('.bulk-layout select').value==='DEACTIVATE'&&!document.querySelector('.bulk-layout select').disabled");await button('Validate Roster');await until("document.body.textContent.includes('1 valid')");await button('Select All Valid Rows');await button('Review & Confirm Deactivation');await until("!!document.querySelector('[role=alertdialog]')");await button('Confirm Deactivation');await until("document.body.textContent.includes('Operation results')");assert.ok(await b.evaluate("document.body.textContent.includes('deactivated')"));pass('Student roster reuse supports reviewed bulk deactivation with history-preservation warning')
  // Each workbook kind exercises the real parser + identity preview, without guessing the owner's workbook.
  for (const [kind,word,status] of [['numeric','text cell',200],['missing','valid Student name',200],['domain','configured institutional',200],['conflict','identity',200],['duplicate','Duplicate',200],['formula','Row 2, Name: Formulas',400]]) {
   const form=new FormData();form.append('roster',new Blob([await fs.readFile('/tmp/elabtrack-phase56-'+kind+'.xlsx')]),kind+'.xlsx');form.append('operation','CREATE');const result=await request('/borrowers/rosters','POST',form);assert.equal(result.status,status,kind)
   const text=status===400?result.message:result.data.rows.map(row=>row.error).join(' ');assert.ok(text.toLowerCase().includes(word.toLowerCase()),kind+' diagnostic: '+text)
  }
  pass('Real workbook numeric/missing/domain/conflict/duplicate row previews and located formula rejection')
  await route('/staff/borrowers/'+faculty.id,"document.body.textContent.includes('Account information')");assert.ok(await b.evaluate("document.body.textContent.includes('Pending activation')&&document.body.textContent.includes('Activation email was not sent.')&&document.body.textContent.includes('Accountability unavailable')"));pass('Account administrative/activation status and unconfigured email are truthful')
  assert.equal((await request('/borrowers','POST',{name:'TEST invalid Student',email:'other@example.invalid',borrower_type:'STUDENT',student_id:'TEST-'+randomUUID()})).status,400)
  const otherFaculty=await ok('/borrowers','POST',{name:'TEST external Faculty',email:'faculty-'+randomUUID()+'@example.invalid',borrower_type:'FACULTY'});assert.equal(otherFaculty.delivery_status,'UNCONFIGURED');pass('Student domain enforced; Faculty accepts valid external email without Student ID')
  // Metadata saved once, simulated image network failure, then image-only retry.
  await route('/staff/inventory/new',"!!document.querySelector('#equipment-name')");await fill('#equipment-name','TEST image partial '+randomUUID().slice(0,8));await file('input[type=file]',imageFile);await until("!!document.querySelector('.catalog-image-preview')")
  let failImage=true
  const imageIntercept=async message=>{if(message.method!=='Fetch.requestPaused')return;const p=message.params;try{if(p.request.url.endsWith('/image')&&p.request.method==='POST'&&failImage){failImage=false;await b.command('Fetch.failRequest',{requestId:p.requestId,errorReason:'Failed'})}else await b.command('Fetch.continueRequest',{requestId:p.requestId})}catch{errors.push('Image interception failed')}}
  b.on(imageIntercept);await b.command('Fetch.enable',{patterns:[{urlPattern:'*/image',requestStage:'Request'}]})
  const creates=()=>network.filter(n=>n.method==='POST'&&n.path==='/api/v1/equipment').length,createStart=creates()
  await button('Save Equipment');await until("document.body.textContent.includes('Equipment information was saved.')");assert.equal(creates()-createStart,1)
  await b.command('Fetch.disable');await button('Retry Image Only');await until("/^\\/staff\\/inventory\\/[^/]+$/.test(location.pathname)");await settled();assert.equal(creates()-createStart,1)
  const createdID=await b.evaluate("location.pathname.split('/').at(-1)");const created=await stock(createdID);assert.ok(created.image_id);pass('Failed image preserves saved equipment; retry uploads image only without duplicate metadata')
  await route('/staff/inventory/new',"!!document.querySelector('#equipment-name')")
  await b.evaluate("(()=>{const dt=new DataTransfer();dt.items.add(new File(['<svg/>'],'fake.svg',{type:'image/svg+xml'}));document.querySelector('.file-dropzone').dispatchEvent(new DragEvent('drop',{bubbles:true,cancelable:true,dataTransfer:dt}))})()")
  await until("document.body.textContent.includes('Choose a PNG or JPEG')");assert.ok(await b.evaluate("!document.querySelector('.catalog-image-preview')"));pass('Image drop validates unsupported file type without a metadata write')
  // Existing image is retained on edit; catalog thumbnail uses authenticated cache.
  await route('/staff/inventory/'+createdID+'/edit',"!!document.querySelector('#equipment-name')");assert.ok(await b.evaluate("document.body.textContent.includes('current image will be kept')&&document.body.textContent.includes('Stored image removal is not supported')"));await file('input[type=file]',imageFile);await until("!!document.querySelector('.catalog-image-preview')");await button('Save Equipment');await until(`location.pathname==='/staff/inventory/${createdID}'`);await settled();assert.notEqual((await stock(createdID)).image_id,created.image_id);pass('Edit shows current image, replaces it through the existing endpoint, and exposes no fake stored removal')
  const imageReads=()=>network.filter(n=>n.method==='GET'&&n.path.includes('/images/')).length
  await route('/staff/inventory',"!!document.querySelector('.inventory-item-identity img')");await settled();const cachedBefore=imageReads();await route('/staff/inventory/'+createdID,"document.body.textContent.includes('Inventory movement history')");await settled();assert.equal(imageReads(),cachedBefore);pass('List/detail reuse cached immutable authenticated image without duplicate download')
  for (const role of ['staff','borrower']) {const response=await fetch(api+'/auth/login',{method:'POST',headers:{Origin:base,'Content-Type':'application/json'},body:JSON.stringify({email:fixtures[role].email,password:fixtures[role].password})});assert.equal(response.status,200);const token=(await response.json()).data.access_token;const denied=await fetch(api+'/equipment/'+equipment.id+'/adjustments',{method:'POST',headers:{Origin:base,'Content-Type':'application/json',Authorization:'Bearer '+token,'Idempotency-Key':randomUUID()},body:JSON.stringify({kind:'RECONCILE',quantity:12,expected_sequence:current.stock_sequence,reason:'Unauthorized test',confirm:true})});assert.equal(denied.status,403)}
  pass('Backend denies Staff and Borrower corrections through direct HTTP')
  const chrome = await b.evaluate("(()=>{window.__phase56Shell=document.querySelector('.staff-sidebar');return !!window.__phase56Shell})()")
  assert.ok(chrome);await route('/staff/inventory',"!!document.querySelector('.management-table')");await route('/staff/borrowers',"document.body.textContent.includes('Borrowers')")
  assert.ok(await b.evaluate("window.__phase56Shell===document.querySelector('.staff-sidebar')"));pass('Persistent application sidebar retains its DOM identity through navigation')
  await b.command('Page.navigate',{url:base+'/staff/inventory'});await until("document.querySelector('#staff-content h1')?.textContent==='Inventory'");await settled();pass('Reload restores the session through the existing refresh-cookie flow')
  await button('Sign Out');await until("location.pathname==='/login'&&!!document.querySelector('#login-email')");assert.ok(await b.evaluate("!document.querySelector('.staff-sidebar')"));pass('Logout removes protected navigation and returns to immersive login')

 }
 pass('before/after screens use real persisted equipment, image, category and unconfigured activation state')
 assert.deepEqual(errors, [])
 if (checksOnly) screens.push(...(await fs.readdir(output)).filter(name=>name.endsWith('.png')).sort())
 await fs.writeFile(path.join(output, 'acceptance.json'), JSON.stringify({ date: new Date().toISOString(), mode, screenshotsReusedFromFinalCapture: checksOnly, ownerAcceptance: 'PENDING', scope: 'Real isolated HTTP/PostgreSQL/Chromium; synthetic accounts and inventory only; no live email', checks, screens, errors }, null, 2) + '\n')
 console.log(JSON.stringify({ result: 'PASS', mode, screenshots: screens.length, checks: checks.length }))
} catch (error) { console.log(JSON.stringify({ failure: error.message, errors, page: await b.evaluate('({path:location.pathname,text:document.body.textContent.slice(-1400)})') })); throw error } finally { await b.close() }
