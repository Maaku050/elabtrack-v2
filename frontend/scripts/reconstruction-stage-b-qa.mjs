// Stage B acceptance: real Chromium, existing HTTP routes and isolated PostgreSQL.
import fs from 'node:fs/promises'
import path from 'node:path'
import assert from 'node:assert/strict'
import { randomUUID } from 'node:crypto'
import { browser } from './browser-cdp.mjs'
const base='http://localhost:15175',api='http://localhost:18085/api/v1'
const output=path.resolve('../docs/ux/verification/reconstruction-stage-b')
const group=process.argv[2] ?? 'shared'
const fixtures=JSON.parse(await fs.readFile('/tmp/elabtrack-batch1-fixtures.json','utf8'))
const checks=[],screens=[],errors=[],documents=[],requests=[]
const b=await browser(), wait=ms=>new Promise(r=>setTimeout(r,ms))
// Retain awaited page promises while CDP waits; Chromium may collect an
// otherwise unreferenced dynamic-import promise during long screenshot runs.
const evaluate=b.evaluate
b.evaluate=expression=>evaluate(`globalThis.__stageBEvaluation=(async()=>(${expression}))();globalThis.__stageBEvaluation`)
let access, intercept, paused=[]
b.on(m=>{
 if(m.method==='Runtime.exceptionThrown')errors.push(m.params.exceptionDetails.text)
 if(m.method==='Network.requestWillBeSent') { if(m.params.type==='Document')documents.push(m.params.request.url); if(m.params.request.url.startsWith(api))requests.push({method:m.params.request.method,path:new URL(m.params.request.url).pathname,query:new URL(m.params.request.url).search}) }
 if(m.method==='Fetch.requestPaused') {
  const p=m.params
  if(intercept==='hold') paused.push(p.requestId)
  else if(intercept==='error'||intercept==='terms') void b.command('Fetch.fulfillRequest',{requestId:p.requestId,responseCode:503,responseHeaders:[{name:'Content-Type',value:'application/json'},{name:'Access-Control-Allow-Origin',value:base},{name:'Access-Control-Allow-Credentials',value:'true'}],body:Buffer.from(JSON.stringify({success:false,error:{code:intercept==='terms'?'TERMS_NOT_PUBLISHED':'SERVICE_UNAVAILABLE',message:'Test network failure',requestId:randomUUID()},meta:{request_id:'stage-a-test'}})).toString('base64')})
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
 await until(`location.pathname===${JSON.stringify(new URL(target,base).pathname)}&&(${ready})`); if(intercept!=='hold') await settled()
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
 await fill('#login-email',fixtures[role].email);await fill('#login-password',fixtures[role].password);await click('button[type=submit]');await until(role==='borrower'?"location.pathname.startsWith('/borrower/')":"!!document.querySelector('#staff-content h1')");if(intercept!=='hold') await settled()
}

async function file(selector, filename) {
 const {root}=await b.command('DOM.getDocument'),{nodeId}=await b.command('DOM.querySelector',{nodeId:root.nodeId,selector})
 assert.ok(nodeId,'File input exists');await b.command('DOM.setFileInputFiles',{nodeId,files:[filename]});await wait(200)
}
async function capture(name) {
 for(const width of [320,390,440,768,1024,1366,1440]) for(const color of ['light','dark']) {
  await view(width,960);await theme(color);await b.evaluate('window.scrollTo(0,0)')
  assert.ok(await b.evaluate("[...document.querySelectorAll('.form-field input,.form-field select,.form-field textarea')].filter(e=>e.getClientRects().length).every(e=>{const r=e.getBoundingClientRect();return r.width>0&&r.left>=0&&r.right<=innerWidth+1})"),name+' field bounds '+width)
  await screenshot(name+'-'+width+'-'+color)
 }
 for(const width of [390,1366]) {await view(width,560);await screenshot(name+'-'+width+'-short-dark')}
 await view();await theme('light')
}
let items
async function saveItems(){await fs.writeFile('/tmp/elabtrack-stage-b-items.json',JSON.stringify(items),{mode:0o600})}
try {
 await fs.mkdir(output,{recursive:true});await b.command('Network.enable');await view();access=await loginToken('admin')
 try {items=JSON.parse(await fs.readFile('/tmp/elabtrack-stage-b-items.json','utf8'))} catch {
  const category=await ok('/equipment-categories','POST',{name:'TEST Stage B culinary tools',is_active:true,expected_version:0})
  const equipment=await ok('/equipment','POST',{name:'TEST Stage B stockpot',description:'Disposable culinary catalog fixture',category_id:category.id,opening_quantity:15,reason:'Synthetic opening count',expected_version:0})
  const faculty=await ok('/borrowers','POST',{name:'TEST Stage B Faculty',email:randomUUID()+'@faculty.example.invalid',borrower_type:'FACULTY'})
  items={category,equipment,faculty};await saveItems()
 }
 await uiLogin('admin');await route('/staff/inventory/new',"!!document.querySelector('#equipment-name')")
 if(group==='shared') {
  await file('input[type=file]','/tmp/elabtrack-batch1-catalog.png');await until("!!document.querySelector('.file-dropzone .catalog-image-preview')")
  assert.ok(await b.evaluate("document.querySelector('.file-dropzone').contains(document.querySelector('.catalog-image-preview'))"))
  await button('Remove selected file');assert.equal(await b.evaluate("!!document.querySelector('.catalog-image-preview')"),false)
  await file('input[type=file]','/tmp/elabtrack-batch1-catalog.png');await until("!!document.querySelector('.file-dropzone .catalog-image-preview')")
  await screenshot('shared-uploader-light');pass('Shared uploader embeds image preview, replacement and unsaved removal through real file input')
 }
 if(group==='borrowers') {
  await route('/staff/borrowers/new',"!!document.querySelector('#account-name')")
  await button('Create Account');await until("document.body.textContent.includes('Enter a full name.')")
  await fill('#account-name','TEST Stage B Student');await fill('#account-email','student-'+randomUUID()+'@students.example.invalid');await fill('#account-student_id','0028366-'+randomUUID().slice(0,8));await fill('#account-course','TEST Culinary Arts');await fill('#account-contact_number','TEST contact')
  assert.equal(await b.evaluate("!!document.querySelector('input[type=password]')"),false)
  await capture('borrower-create-student');await button('Create Account');await until("location.pathname.startsWith('/staff/borrowers/')&&!location.pathname.endsWith('/new')");await settled()
  items.student=await ok('/borrowers/'+await b.evaluate("location.pathname.split('/').pop()"));await saveItems()
  assert.equal(items.student.borrower_type,'STUDENT');assert.ok(items.student.student_id.startsWith('00'));assert.equal(items.student.activation_required,true);assert.equal(items.student.delivery_status,'UNCONFIGURED')
  await capture('borrower-details');pass('Student creation preserves textual ID/domain, pending activation and no administrator password')
  await fill('#profile-name','TEST Stage B Student updated');await button('Save Profile');await until("document.body.textContent.includes('Profile updated.')");await settled()
  assert.equal((await ok('/borrowers/'+items.student.id)).name,'TEST Stage B Student updated')
  await button('Deactivate Account');await until("!!document.querySelector('[role=alertdialog]')");await screenshot('borrower-deactivation-review-light')
  await b.command('Input.dispatchKeyEvent',{type:'keyDown',key:'Tab',code:'Tab',windowsVirtualKeyCode:9});await b.command('Input.dispatchKeyEvent',{type:'keyUp',key:'Tab',code:'Tab',windowsVirtualKeyCode:9})
  assert.ok(await b.evaluate("document.querySelector('[role=alertdialog]').contains(document.activeElement)"))
  await button('Confirm Deactivation');await until("document.body.textContent.includes('Account deactivated; obligations and history preserved.')");assert.equal((await ok('/borrowers/'+items.student.id)).is_active,false)
  await button('Reactivate Account');await until("!!document.querySelector('[role=alertdialog]')");await button('Confirm Reactivation');await until("document.body.textContent.includes('Account reactivated.')")
  pass('Profile editing, confirmed deactivation/reactivation, audit display and alert-dialog keyboard containment')
  await route('/staff/borrowers/new',"!!document.querySelector('#borrower-type')");await select('#borrower-type','FACULTY');await until("!document.querySelector('#account-student_id')")
  await fill('#account-name','TEST Individual Faculty');await fill('#account-email','faculty-'+randomUUID()+'@external.example.invalid');await capture('borrower-create-faculty');await button('Create Account');await until("location.pathname.startsWith('/staff/borrowers/')&&!location.pathname.endsWith('/new')");await settled()
  items.individualFaculty=await ok('/borrowers/'+await b.evaluate("location.pathname.split('/').pop()"));await saveItems();assert.equal(items.individualFaculty.borrower_type,'FACULTY');assert.equal(items.individualFaculty.student_id,'')
  pass('Faculty individual creation permits external email and omits Student ID')
  await route('/activate',"document.body.textContent.includes('Open the activation link')");await capture('activation-missing-link');pass('Activation UI retains safe missing-link state; no token or password exposed')
 }
 if(group==='bulk') {
  await route('/admin/borrowers/bulk',"document.body.textContent.includes('Student roster')");await capture('student-bulk')
  await button('Download Student Template');await settled();assert.ok(requests.some(r=>r.path.endsWith('/student-template')))
  await file('input[type=file]','/tmp/elabtrack-batch1-students.xlsx');await button('Validate Roster');await until("!!document.querySelector('.directory-table-region tbody')");await settled()
  assert.equal(await b.evaluate("document.querySelectorAll('.directory-table-region tbody tr').length"),4)
  assert.equal(await b.evaluate("document.querySelectorAll('.directory-table-region input:disabled').length"),3)
  await capture('bulk-validation');await button('Select All Valid Rows');await button('Review & Confirm Creation');await until("!!document.querySelector('[role=alertdialog]')");await capture('bulk-creation-confirmation')
  const prior=requests.filter(r=>r.method==='POST'&&r.path.endsWith('/confirm')).length
  await b.evaluate("(()=>{const e=[...document.querySelectorAll('button')].find(e=>e.textContent.trim()==='Confirm Creation');e.click();e.click()})()")
  await until("document.body.textContent.includes('Operation results')");await settled();assert.equal(requests.filter(r=>r.method==='POST'&&r.path.endsWith('/confirm')).length-prior,1)
  const accountId=await b.evaluate("document.querySelector('.directory-table-region tbody a').getAttribute('href').split('/').pop()")
  items.bulkStudent=await ok('/borrowers/'+accountId);await saveItems();assert.equal(items.bulkStudent.borrower_type,'STUDENT');assert.equal(items.bulkStudent.activation_required,true)
  await capture('bulk-creation-results');pass('Complete Excel preview preserves duplicate/domain errors; selected Student creation confirms once with activation pending')
  await select('#bulk-operation','DEACTIVATE');await button('Validate Roster');await until("document.body.textContent.includes('Validation preview')&&!!document.querySelector('.directory-table-region tbody')");await settled()
  await button('Select All Valid Rows');await button('Review & Confirm Deactivation');await until("!!document.querySelector('[role=alertdialog]')");await capture('bulk-deactivation-confirmation');await button('Confirm Deactivation');await until("document.body.textContent.includes('Operation results')");await settled()
  assert.equal((await ok('/borrowers/'+accountId)).is_active,false);assert.equal((await ok('/borrowers/'+items.faculty.id)).is_active,true)
  await screenshot('bulk-deactivation-results-light');pass('Roster reuse deactivates only matched selected Students; Faculty and history remain intact')
 }
 if(group==='categories') {
  if(!items.categoryPaging){for(let i=0;i<100;i++)await ok('/equipment-categories','POST',{name:'TEST A Category '+String(i).padStart(3,'0'),is_active:true,expected_version:0});items.categoryPaging=true;await saveItems()}
  await b.evaluate("(async()=>{const loaded=performance.getEntriesByType('resource').find(e=>new URL(e.name).pathname==='/src/app/query-client.ts');const {queryClient}=await import(loaded.name);await queryClient.invalidateQueries({queryKey:['inventory']})})()");await settled()
  await route('/staff/inventory/categories',"!!document.querySelector('.directory-table-region tbody')");await capture('equipment-categories')
  assert.equal(await b.evaluate("document.querySelectorAll('.directory-table-region tbody tr').length"),100)
  await click('[aria-label="Go to next page"]');await until("new URLSearchParams(location.search).get('page')==='2'");await settled();assert.equal(await b.evaluate("document.querySelector('[aria-label=\"Go to next page\"]').getAttribute('aria-disabled')"),'true')
  await click('[aria-label="Go to previous page"]');await settled();await button('Add Category');await until("!!document.querySelector('#category-name')")
  await button('Save Category');await until("document.body.textContent.includes('Enter a category name.')");const name='AAA TEST UI category '+randomUUID().slice(0,6);await fill('#category-name',name);await capture('category-create-dialog');await button('Save Category');await until("!document.querySelector('[role=dialog]')&&document.body.textContent.includes('Category created.')");await settled()
  items.uiCategory=(await ok('/equipment-categories?page=1')).find(c=>c.name===name);assert.ok(items.uiCategory);await saveItems()
  await click('button[aria-label='+JSON.stringify('Edit '+name)+']');await until("!!document.querySelector('#category-name')");await click('.management-checkbox input');await screenshot('category-edit-dialog-light');await button('Save Category');await until("!document.querySelector('[role=dialog]')&&document.body.textContent.includes('Category updated.')");await settled()
  assert.equal((await ok('/equipment-categories?page=1')).find(c=>c.id===items.uiCategory.id).is_active,false)
  assert.equal((await ok('/equipment/'+items.equipment.id)).category_id,items.category.id)
  pass('Category add/edit/inactive relationship preservation and honest unknown-total server pagination')
 }
 if(group==='equipment') {
  await fill('#equipment-name','TEST Stage B image equipment '+randomUUID().slice(0,6));await fill('#equipment-description','Synthetic equipment description');await fill('#equipment-opening','15');await fill('#equipment-reason','Synthetic opening explanation');await file('input[type=file]','/tmp/elabtrack-batch1-catalog.png');await until("!!document.querySelector('.file-dropzone .catalog-image-preview')");await capture('equipment-add-image');await screenshot('equipment-add-image-area-light',true)
  const prior=requests.filter(r=>r.method==='POST'&&r.path==='/api/v1/equipment').length
  intercept='error';await b.command('Fetch.enable',{patterns:[{urlPattern:api+'/equipment/*/image',requestStage:'Request'}]})
  await button('Save Equipment');await until("document.body.textContent.includes('The image was not confirmed.')");assert.equal(requests.filter(r=>r.method==='POST'&&r.path==='/api/v1/equipment').length-prior,1);await screenshot('equipment-image-retry-light',true)
  intercept=undefined;await b.command('Fetch.disable');await button('Retry Image Only');await until("location.pathname.startsWith('/staff/inventory/')&&!location.pathname.endsWith('/new')");await settled();assert.equal(requests.filter(r=>r.method==='POST'&&r.path==='/api/v1/equipment').length-prior,1)
  items.imageEquipment=await ok('/equipment/'+await b.evaluate("location.pathname.split('/').pop()"));await saveItems();assert.ok(items.imageEquipment.image_id);assert.equal(items.imageEquipment.stock.available,15);pass('Metadata persists through upload failure; retry uploads image only without duplicate equipment')
  await route('/staff/inventory/'+items.imageEquipment.id+'/edit',"!!document.querySelector('#equipment-name')");await until("!!document.querySelector('.file-dropzone img')");await capture('equipment-edit-image');await screenshot('equipment-edit-image-area-light',true)
  await file('input[type=file]','/tmp/elabtrack-batch1-catalog.png');await until("!!document.querySelector('.catalog-image-preview')");await button('Remove selected file');await until("!!document.querySelector('.file-dropzone img')&&!document.querySelector('.catalog-image-preview')")
  await file('input[type=file]','/tmp/elabtrack-batch1-catalog.png');await until("!!document.querySelector('.catalog-image-preview')");await fill('#equipment-name','TEST Stage B replacement verified');await button('Save Equipment');await until("location.pathname==='/staff/inventory/"+items.imageEquipment.id+"'");await settled()
  const edited=await ok('/equipment/'+items.imageEquipment.id);assert.equal(edited.name,'TEST Stage B replacement verified');assert.notEqual(edited.image_id,items.imageEquipment.image_id);assert.equal(edited.stock.available,15);items.imageEquipment=edited;await saveItems();pass('Edit embeds existing image, keeps it after removing unsaved selection, and persists replacement plus metadata')
 }
 if(group==='details') {
  await route('/staff/inventory/'+items.imageEquipment.id,"document.body.textContent.includes('Current physical stock')");await until("!!document.querySelector('.inventory-detail-image img')");assert.equal(await b.evaluate("document.querySelectorAll('input[type=file]').length"),0);await capture('equipment-details');assert.ok(await b.evaluate("[...document.querySelectorAll('a')].some(e=>e.getAttribute('href')?.endsWith('/edit'))"));assert.ok(await b.evaluate("document.body.textContent.includes('Total pages unavailable')"))
  const before=(await ok('/equipment/'+items.imageEquipment.id)).stock
  await button('Make Inactive');await until("!!document.querySelector('[role=alertdialog]')");await screenshot('equipment-status-dialog-light');await button('Confirm Status Change');await until("!document.querySelector('[role=alertdialog]')");await settled();assert.equal((await ok('/equipment/'+items.imageEquipment.id)).status,'INACTIVE');await button('Make Active');await until("!!document.querySelector('[role=alertdialog]')");await button('Confirm Status Change');await until("!document.querySelector('[role=alertdialog]')");await settled();assert.deepEqual((await ok('/equipment/'+items.imageEquipment.id)).stock,before);assert.ok((await ok('/equipment/'+items.imageEquipment.id+'/movements?page=1')).length>0)
  pass('Details displays protected image and Edit action without duplicate uploader; status changes preserve physical stock and movement history')
 }
 if(group==='stock') {
  const equipment=await ok('/equipment','POST',{name:'TEST Stage B stock workflow',description:'Disposable stock fixture',category_id:items.category.id,opening_quantity:15,reason:'Synthetic opening count',expected_version:0});items.stockEquipment=equipment;await saveItems()
  const stock=id=>ok('/equipment/'+id)
  await route('/staff/inventory/'+equipment.id+'/adjust',"!!document.querySelector('#adjust-quantity')");await fill('#adjust-quantity','5');await fill('#adjust-reason','Verified synthetic addition');await capture('stock-adjustment');await button('Review Stock Change');await until("document.body.textContent.includes('Change to confirm')");await capture('stock-review')
  await route('/admin/inventory/'+equipment.id+'/reconcile',"!!document.querySelector('#adjust-quantity')");await fill('#adjust-quantity','20');await fill('#adjust-reason','Verified synthetic available count');await capture('inventory-correction');await button('Review Stock Change');await until("document.body.textContent.includes('Change to confirm')");await capture('correction-review')
  // The frozen correction review above must commit exactly once, even on a double click.
  const countAdjustments = () => requests.filter(n => n.method === 'POST' && n.path.endsWith('/adjustments')).length
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
  for(const role of ['staff','borrower']) {const token=await loginToken(role);const r=await request('/equipment/'+equipment.id+'/adjustments','POST',{kind:'RECONCILE',quantity:12,expected_sequence:current.stock_sequence,reason:'Unauthorized synthetic correction',confirm:true},token);assert.equal(r.status,403)}
  assert.deepEqual([current.stock.reserved,current.stock.checked_out,current.stock.damaged_held],[0,0,0]);pass('Backend correction permission remains Admin-only; unavailable buckets unchanged')
 }
 if(group==='administration') {
  if(!items.staffPaging){for(let i=0;i<26;i++)await ok('/staff-accounts','POST',{name:'TEST Stage B Staff '+String(i).padStart(2,'0'),email:randomUUID()+'@staff.example.invalid'});items.staffPaging=true;await saveItems()}
  await route('/admin/administration',"document.body.textContent.includes('Staff / Admin account directory')");await capture('administration-overview');assert.equal(await b.evaluate("document.querySelector('.directory-search input').placeholder"),'Name or email');await click('[aria-label="Go to next page"]');await until("new URLSearchParams(location.search).get('staff_page')==='2'");await settled();assert.ok(await b.evaluate("!new URLSearchParams(location.search).has('page')"));pass('Administration uses compact overview, actual terms status and embedded server directory with namespaced paging')
  await route('/admin/administration/accounts',"!!document.querySelector('.directory-table-region tbody')");await capture('staff-admin-directory');await click('[aria-label="Go to next page"]');await settled();await fill('.directory-search input','TEST Stage B Staff 00');await until("document.querySelectorAll('tbody tr').length===1&&document.body.textContent.includes('1 matching accounts')");await settled();assert.ok(await b.evaluate("!new URLSearchParams(location.search).has('page')"));assert.equal(await b.evaluate("document.querySelectorAll('tbody tr').length"),1);assert.equal(await b.evaluate("document.querySelector('[aria-label=\"Go to next page\"]').getAttribute('aria-disabled')"),'true');pass('Staff/Admin search describes name/email; paging and filter reset follow real server totals')
  await route('/admin/administration/accounts/new',"!!document.querySelector('#account-name')");assert.equal(await b.evaluate("!!document.querySelector('#borrower-type')||!!document.querySelector('input[type=password]')"),false);await fill('#account-name','TEST Stage B newly created Staff');await fill('#account-email',randomUUID()+'@staff.example.invalid');await capture('staff-create');await button('Create Account');await until("location.pathname.startsWith('/admin/administration/accounts/')&&!location.pathname.endsWith('/new')");await settled();items.uiStaff=await ok('/staff-accounts/'+await b.evaluate("location.pathname.split('/').pop()"));await saveItems();assert.equal(items.uiStaff.role,'STAFF');assert.equal(items.uiStaff.activation_required,true);await capture('staff-details')
  await button('Deactivate Account');await until("!!document.querySelector('[role=alertdialog]')");await button('Confirm Deactivation');await until("document.body.textContent.includes('Account deactivated; obligations and history preserved.')");assert.equal((await ok('/staff-accounts/'+items.uiStaff.id)).is_active,false);pass('Existing fixed-STAFF creation, pending activation, detail audit and confirmed status controls persist')
  await route('/admin/administration/accounts/'+fixtures.admin.id,"document.body.textContent.includes('Existing Admin accounts are read-only here.')");assert.equal(await b.evaluate("[...document.querySelectorAll('button')].some(e=>['Deactivate Account','Save Profile'].includes(e.textContent.trim()))"),false);await screenshot('admin-read-only-details-light');pass('Existing Admin remains read-only; no privileged provisioning or role/password control added')
 }
 if(group==='consistency') {
  await route('/staff/inventory',"!!document.querySelector('.directory-inventory-table')");await view(1366);assert.ok(await b.evaluate("(()=>{const e=document.querySelector('.directory-table-region');return e.scrollWidth<=e.clientWidth+1})()"),'Inventory table fits expanded 1366 desktop with all stock columns');await capture('inventory-final')
  const documentCount=documents.length
  await b.evaluate("(()=>{globalThis.stageBFrame={sidebar:document.querySelector('.staff-sidebar'),header:document.querySelector('.staff-top-bar'),main:document.querySelector('#staff-content')};return true})()")
  for(const target of ['/staff/dashboard','/staff/borrowers','/admin/administration','/staff/inventory/'+items.imageEquipment.id+'/edit','/staff/inventory/'+items.stockEquipment.id+'/adjust','/admin/administration/accounts'])await route(target,"!!document.querySelector('#staff-content h1')")
  assert.ok(await b.evaluate("stageBFrame.sidebar===document.querySelector('.staff-sidebar')&&stageBFrame.header===document.querySelector('.staff-top-bar')&&stageBFrame.main===document.querySelector('#staff-content')"));assert.equal(documents.length,documentCount);pass('Operational shell DOM persists through Stage B routes without document reload; all inventory stock columns fit at 1366 expanded desktop')
  await click('[aria-label="Toggle navigation"]');await until("document.querySelector('.fsmo-navigation-body').dataset.collapsed==='true'");await wait(200);assert.ok(await b.evaluate("(()=>{const a=document.querySelector('.fsmo-brand-link img').getBoundingClientRect(),s=document.querySelector('.staff-sidebar').getBoundingClientRect();return Math.abs(a.x+a.width/2-(s.x+s.width/2))<1})()"));await screenshot('staff-directory-collapsed-light');await theme('dark');await screenshot('staff-directory-collapsed-dark');await theme('light');await click('[aria-label="Toggle navigation"]');await until("document.querySelector('.fsmo-navigation-body').dataset.collapsed==='false'")
  await click('[aria-label="Account actions"]');await until("!!document.querySelector('[role=menu]')");await screenshot('account-menu-light');await b.command('Input.dispatchKeyEvent',{type:'keyDown',key:'Escape',code:'Escape',windowsVirtualKeyCode:27});await until("!document.querySelector('[role=menu]')");assert.equal(await b.evaluate("document.activeElement.getAttribute('aria-label')"),'Account actions')
  for(const width of [320,390])for(const color of ['light','dark']){await view(width,560);await theme(color);await click('[aria-label="Toggle navigation"]');await until("!!document.querySelector('[role=dialog]')");await screenshot('navigation-'+width+'-short-'+color);await b.command('Input.dispatchKeyEvent',{type:'keyDown',key:'Tab',code:'Tab',windowsVirtualKeyCode:9});assert.ok(await b.evaluate("document.querySelector('[role=dialog]').contains(document.activeElement)"));await b.command('Input.dispatchKeyEvent',{type:'keyDown',key:'Escape',code:'Escape',windowsVirtualKeyCode:27});await until("!document.querySelector('[role=dialog]')");assert.equal(await b.evaluate("document.activeElement.getAttribute('aria-label')"),'Toggle navigation')}
  await view();await theme('light');pass('Collapsed seal remains centered; account menu and short mobile Sheet preserve keyboard containment and Escape focus return in both themes')
  // Expected policy absence is controlled at the network boundary; no terms are unpublished in any database.
  intercept='terms';await b.command('Fetch.enable',{patterns:[{urlPattern:api+'/terms/current',requestStage:'Request'}]});const beforeTerms=requests.filter(r=>r.path.endsWith('/terms/current')).length
  await b.evaluate("(async()=>{const loaded=performance.getEntriesByType('resource').find(e=>new URL(e.name).pathname==='/src/app/query-client.ts');const {queryClient}=await import(loaded.name);queryClient.removeQueries({queryKey:['terms']});return true})()");await route('/admin/administration',"document.body.textContent.includes('Official FSMO terms await institutional approval')");await screenshot('administration-terms-unpublished-light');await route('/staff/inventory',"!!document.querySelector('.directory-inventory-table')");await route('/admin/administration',"document.body.textContent.includes('Official FSMO terms await institutional approval')");await wait(1200);assert.equal(requests.filter(r=>r.path.endsWith('/terms/current')).length-beforeTerms,1);intercept=undefined;await b.command('Fetch.disable');pass('TERMS_NOT_PUBLISHED is truthful unavailable policy with no retry or redundant immediate navigation request')
  await route('/staff/inventory/'+items.imageEquipment.id,"!!document.querySelector('.inventory-detail-image')");await capture('equipment-details-final')
  await route('/staff/inventory/'+items.stockEquipment.id+'/adjust',"!!document.querySelector('#adjust-quantity')");await fill('#adjust-reason','Synthetic final presentation');await screenshot('stock-adjustment-fields-desktop-light',true);await view(390);await screenshot('stock-adjustment-projection-mobile-light',true);await theme('dark');await screenshot('stock-adjustment-projection-mobile-dark',true);await view();await theme('light')
  const staffToken=await loginToken('staff');assert.equal((await request('/staff-accounts','GET',undefined,staffToken)).status,403);assert.equal((await request('/borrowers','POST',{name:'TEST denied',email:randomUUID()+'@external.example.invalid',borrower_type:'FACULTY'},staffToken)).status,403)
  await uiLogin('staff');for(const target of ['/admin/administration','/admin/administration/accounts/new','/admin/inventory/'+items.stockEquipment.id+'/reconcile','/staff/borrowers/new'])await route(target,"document.body.textContent.includes('Access denied')");pass('STAFF direct routes and backend account administration/creation remain denied')
  const borrowerToken=await loginToken('borrower'),terms=await request('/terms/current','GET',undefined,borrowerToken);assert.equal(terms.status,200);assert.ok(terms.data.title.startsWith('TEST ONLY'));const accepted=await request('/terms/'+terms.data.id+'/accept','POST',{},borrowerToken);assert.ok(accepted.status<300)
  await uiLogin('borrower');await route('/borrower/equipment',"document.body.textContent.includes('Browse equipment')");await capture('borrower-equipment-catalog');await route('/borrower/equipment/'+items.imageEquipment.id,"document.body.textContent.includes('Current availability')");assert.equal(await b.evaluate("!!document.querySelector('input[type=file]')||[...document.querySelectorAll('a')].some(e=>e.getAttribute('href')?.endsWith('/edit'))"),false);await capture('borrower-equipment-details');pass('Borrower catalog/detail retain terms gate, protected image and read-only availability without stock, edit or upload actions')
 }
 if(group==='edges') {
  const max=await ok('/equipment','POST',{name:'AAA TEST Long inventory identity '+randomUUID(),description:'Synthetic boundary quantity fixture',category_id:items.category.id,opening_quantity:2147483647,reason:'Synthetic maximum count',expected_version:0});const archived=await ok('/equipment/'+max.id+'/status','PATCH',{status:'ARCHIVED',expected_version:max.metadata_version,confirm:true});assert.equal(archived.status,'ARCHIVED');await route('/staff/inventory',"!!document.querySelector('.directory-inventory-table')");await view(1366)
  const bounds=await b.evaluate("(()=>{const row=[...document.querySelectorAll('tbody tr')].find(r=>r.textContent.includes('AAA TEST Long inventory identity'));return {cells:[...row.querySelectorAll('td')].map(e=>({cell:e.getBoundingClientRect().width,contents:[...e.children].map(c=>({w:c.getBoundingClientRect().width,within:c.getBoundingClientRect().right<=e.getBoundingClientRect().right}))})),scroll:document.querySelector('.directory-table-region').scrollWidth,client:document.querySelector('.directory-table-region').clientWidth}})()")
  console.log('Boundary layout',JSON.stringify(bounds));assert.equal(bounds.scroll,bounds.client);assert.ok(bounds.cells[5].contents.every(c=>c.within),'Archived status fits its column');assert.ok(bounds.cells[6].contents.every(c=>c.within),'Open action fits its column');await capture('inventory-final');await view(1366);await screenshot('inventory-maximum-stock-1366-light');await theme('dark');await screenshot('inventory-maximum-stock-1366-dark');pass('Maximum supported stock, long identity and archived status remain readable without desktop table overflow')
  await view();await theme('light');await route('/staff/inventory/new',"!!document.querySelector('#equipment-name')")
  const bytes=await fs.readFile('/tmp/elabtrack-batch1-catalog.png');await b.evaluate(`(()=>{const bytes=Uint8Array.from(atob(${JSON.stringify(bytes.toString('base64'))}),c=>c.charCodeAt(0)),transfer=new DataTransfer();transfer.items.add(new File([bytes],'TEST-drop.png',{type:'image/png'}));const e=document.querySelector('.file-dropzone');e.dispatchEvent(new DragEvent('dragover',{bubbles:true,cancelable:true,dataTransfer:transfer}));e.dispatchEvent(new DragEvent('drop',{bubbles:true,cancelable:true,dataTransfer:transfer}));return true})()`);await until("!!document.querySelector('.catalog-image-preview')");assert.ok(await b.evaluate("document.querySelector('.file-dropzone').contains(document.querySelector('.catalog-image-preview'))"));await screenshot('catalog-image-drag-drop-light',true);await button('Remove selected file')
  await b.evaluate("(()=>{const transfer=new DataTransfer();transfer.items.add(new File(['<svg/>'],'TEST-invalid.svg',{type:'image/svg+xml'}));document.querySelector('.file-dropzone').dispatchEvent(new DragEvent('drop',{bubbles:true,cancelable:true,dataTransfer:transfer}));return true})()");await until("document.body.textContent.includes('Choose a PNG or JPEG')");assert.equal(await b.evaluate("!!document.querySelector('.catalog-image-preview')"),false);pass('Real Chromium DataTransfer drop uses unchanged image validator and rejects unsupported SVG')
  await route('/staff/inventory/categories',"!!document.querySelector('tbody')");await button('Add Category');await until("!!document.querySelector('[role=dialog]')");await until("document.querySelector('[role=dialog]').contains(document.activeElement)");for(let i=0;i<8;i++){await b.command('Input.dispatchKeyEvent',{type:'keyDown',key:'Tab',code:'Tab',windowsVirtualKeyCode:9});await b.command('Input.dispatchKeyEvent',{type:'keyUp',key:'Tab',code:'Tab',windowsVirtualKeyCode:9});await until("document.querySelector('[role=dialog]').contains(document.activeElement)")}await b.command('Input.dispatchKeyEvent',{type:'keyDown',key:'Escape',code:'Escape',windowsVirtualKeyCode:27});await until("!document.querySelector('[role=dialog]')");pass('Category Dialog keyboard focus remains contained and Escape dismisses safely')
 }
 if(group==='final-visual') {
  await route('/staff/inventory/categories',"!!document.querySelector('tbody')");await button('Add Category');await until("!!document.querySelector('[role=dialog]')");await until("(()=>{const r=document.querySelector('[data-slot=dialog-close]').getBoundingClientRect();return r.width>=44&&r.height>=44})()");await capture('category-create-dialog');await button('Cancel')
  for(const correction of [false,true]) {await route((correction?'/admin/inventory/':'/staff/inventory/')+items.stockEquipment.id+(correction?'/reconcile':'/adjust'),"!!document.querySelector('#adjust-quantity')");await fill('#adjust-quantity',correction?'20':'5');await fill('#adjust-reason','Synthetic final visual verification');await button('Review Stock Change');await until("document.body.textContent.includes('Change to confirm')");for(const width of [320,1366])for(const color of ['light','dark']){await view(width,960);await theme(color);await screenshot((correction?'correction':'stock')+'-review-projection-'+width+'-'+color,true)}}
  pass('Category close control has 44px target; current final dialog and stock/correction review/projection captured in both themes and narrow mobile')
 }
 // Group-specific checks are added as each approved implementation group lands.
 assert.deepEqual(errors,[])
 await fs.writeFile(path.join(output,group+'-acceptance.json'),JSON.stringify({result:'PASS',ownerVisualAcceptance:'PENDING',group,checks,screens,errors,apiRequests:requests,documentNavigations:documents.length},null,2)+'\n')
 console.log(JSON.stringify({result:'PASS',group,checks:checks.length,screens:screens.length}))
} catch(error){await screenshot(group+'-failure').catch(()=>{});throw error} finally {await b.close()}
