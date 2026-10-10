// Real forms, PostgreSQL-backed API and Chromium; no fabricated auth responses.
import fs from 'node:fs/promises'
import path from 'node:path'
import { execFileSync } from 'node:child_process'
import { browser } from './browser-cdp.mjs'
const base = 'http://localhost:15174'
const output = path.resolve('../docs/ux/verification/phase4b')
const fixtures = JSON.parse(await fs.readFile('/tmp/elabtrack-phase4b-fixtures.json', 'utf8'))
await fs.mkdir(output, { recursive: true })
const b = await browser(), checks = [], screens = [], errors = [], warnings = [], responses = []
const sessionEvents = []
let holdAccept = false, paused = [], cookiePolicies
b.on(async message => {
  if (message.method === 'Runtime.exceptionThrown') errors.push({ text: message.params.exceptionDetails.text })
  if (message.method === 'Runtime.consoleAPICalled' && ['error', 'warning'].includes(message.params.type)) {
    const destination = message.params.type === 'error' ? errors : warnings
    destination.push({ text: message.params.args.map(arg => arg.value ?? arg.type).join(' ') })
  }
  if (message.method === 'Network.responseReceived' && message.params.response.url.includes('/api/v1/')) responses.push({ path: new URL(message.params.response.url).pathname, status: message.params.response.status })
  if (message.method === 'Fetch.requestPaused') {
    if (holdAccept && message.params.request.url.endsWith('/accept') && message.params.request.method === 'POST') paused.push({ id: message.params.requestId, session: message.sessionId })
    else await b.command('Fetch.continueRequest', { requestId: message.params.requestId }, message.sessionId)
  }
})
const wait = ms => new Promise(resolve => setTimeout(resolve, ms))
async function evaluate(expression, session) {
  if (!session) return b.evaluate(expression)
  const result = await b.command('Runtime.evaluate', { expression, awaitPromise: true, returnByValue: true }, session)
  if (result.exceptionDetails) throw new Error('Secondary tab evaluation failed')
  return result.result.value
}
async function until(expression, session) {
  for (let i = 0; i < 150; i++) {
    try { if (await evaluate(expression, session)) return } catch (error) { if (!/context|Cannot find|Inspected target/.test(error.message)) throw error }
    await wait(100)
  }
  throw new Error(`Browser condition timed out: ${expression}`)
}
async function settle() { await b.evaluate('document.fonts.ready'); await wait(160) }
async function navigate(route, session) {
  console.log('NAVIGATE',route,session?'peer':'main')
  await b.command('Page.bringToFront',{},session)
  const marker = `navigation-${Date.now()}-${Math.random()}`
  await evaluate(`document.documentElement.dataset.phase4b=${JSON.stringify(marker)}`, session)
  await b.command('Page.navigate', { url: base + route }, session)
  await until(`document.documentElement?.dataset.phase4b !== ${JSON.stringify(marker)} && document.readyState==='complete'`, session)
}
async function click(selector) {
  const rect = await b.evaluate(`(()=>{const e=document.querySelector(${JSON.stringify(selector)});if(!e)throw Error('Missing control');e.scrollIntoView({block:'center'});const r=e.getBoundingClientRect();return {x:r.x+r.width/2,y:r.y+r.height/2};})()`)
  await b.command('Input.dispatchMouseEvent', { type: 'mousePressed', ...rect, button: 'left', clickCount: 1 })
  await b.command('Input.dispatchMouseEvent', { type: 'mouseReleased', ...rect, button: 'left', clickCount: 1 })
}
async function button(text) { await b.evaluate(`(()=>{const e=[...document.querySelectorAll('button')].find(e=>e.textContent.trim()===${JSON.stringify(text)});if(!e)throw Error('Missing named button');e.dataset.phase4bButton='true';})()`); await click('[data-phase4b-button="true"]'); await b.evaluate(`document.querySelector('[data-phase4b-button]')?.removeAttribute('data-phase4b-button')`) }
async function key(key, code, number) { for (const type of ['keyDown','keyUp']) await b.command('Input.dispatchKeyEvent', { type, key, code, windowsVirtualKeyCode: number }) }
async function fill(selector, value) {
  await click(selector); await b.command('Input.dispatchKeyEvent', { type: 'keyDown', key: 'a', code: 'KeyA', windowsVirtualKeyCode: 65, modifiers: 2 }); await b.command('Input.dispatchKeyEvent', { type: 'keyUp', key: 'a', code: 'KeyA', windowsVirtualKeyCode: 65, modifiers: 2 }); await b.command('Input.insertText', { text: value })
}
async function login(kind, destination) {
  const fixture = fixtures[kind]
  await until(`!!document.querySelector('#login-email') && !document.querySelector('#login-email').disabled && !![...document.querySelectorAll('button')].find(e=>e.textContent.trim()==='Sign In' && !e.disabled)`)
  await fill('#login-email', fixture.email); await fill('#login-password', fixture.password); await button('Sign In')
  await until(`location.pathname===${JSON.stringify(destination)} `)
}
async function signOut() {
  await button('Sign Out'); await until(`location.pathname==='/login' && !!document.querySelector('#login-email') && !document.body.innerText.includes('Signing out…')`)
}
async function theme(value) {
  if (await b.evaluate(`document.documentElement.classList.contains('dark')`) !== (value === 'dark')) await click(`button[aria-label="Switch to ${value} theme"]`)
  await settle()
}
async function viewport(width) { await b.command('Emulation.setDeviceMetricsOverride', { width, height: 900, deviceScaleFactor: 1, mobile: width < 768 }); await settle() }
async function capture(name, width, mode) {
  await viewport(width); await theme(mode); await b.evaluate('window.scrollTo(0,0)'); await settle()
  const geometry = await b.evaluate(`(()=>{const visible=e=>e.getBoundingClientRect().width>0&&e.getBoundingClientRect().height>0&&getComputedStyle(e).visibility!=='hidden';return {width:innerWidth,pageWidth:document.documentElement.scrollWidth,clientWidth:document.documentElement.clientWidth,theme:document.documentElement.classList.contains('dark')?'dark':'light',smallTargets:[...document.querySelectorAll('button,a,input')].filter(visible).filter(e=>!e.classList.contains('skip-link')).map(e=>({name:e.getAttribute('aria-label')||e.textContent.trim(),width:(e.closest('.terms-consent')||e).getBoundingClientRect().width,height:(e.closest('.terms-consent')||e).getBoundingClientRect().height})).filter(e=>e.width<43.9||e.height<43.9),h1:document.querySelector('h1')?.textContent};})()`)
  const result = await b.command('Page.captureScreenshot', { format: 'png', captureBeyondViewport: false })
  const file = `${name}-${width}-${mode}.png`; await fs.writeFile(path.join(output, file), Buffer.from(result.data, 'base64'))
  console.log('CAPTURE', file)
  screens.push({ file, ...geometry, pass: geometry.pageWidth <= geometry.clientWidth && geometry.smallTargets.length === 0 && geometry.theme === mode })
  if (!screens.at(-1).pass) throw new Error(`Responsive geometry failed: ${file}: ${JSON.stringify(geometry)}`)
}
async function check(name, expression, session) { const pass = !!await evaluate(expression, session); console.log(pass ? 'PASS' : 'FAIL', name); checks.push({ name, pass }); if (!pass) throw new Error(`Check failed: ${name}`) }
function state(action, kind, role) { execFileSync('python3', ['../integration/phase4b-fixture-state.py', action, kind, ...(role ? [role] : [])], { stdio: 'pipe' }) }
const store = `(await import('/src/stores/auth-store.ts')).useAuthStore.getState()`
const api='http://localhost:18084/api/v1'
let adminAccess, borrowerAccess
const runSuffix=Date.now().toString(36)
async function directLogin(kind) {
 const response=await fetch(api+'/auth/login',{method:'POST',headers:{'Content-Type':'application/json',Origin:base},body:JSON.stringify({email:fixtures[kind].email,password:fixtures[kind].password})})
 if(!response.ok)throw Error('Isolated direct login failed')
 return (await response.json()).data.access_token
}
async function publish(number) {
 const currentResponse=await fetch(api+'/terms/current',{headers:{Authorization:'Bearer '+adminAccess}})
 const current=currentResponse.ok?(await currentResponse.json()).data.id:null
 const body='SYNTHETIC TEST DOCUMENT — NOT OFFICIAL FSMO POLICY.\n\n'+Array.from({length:24},(_,i)=>`Test paragraph ${i+1}. This is isolated browser verification content, not institutional obligations. It verifies complete-document reading, page scrolling and exact-version consent. No actual borrower is accepting this text.`).join('\n\n')+'\n\nEND OF SYNTHETIC TEST DOCUMENT.'
 const response=await fetch(api+'/terms/versions',{method:'POST',headers:{Authorization:'Bearer '+adminAccess,Origin:base,'Content-Type':'application/json'},body:JSON.stringify({version:`TEST-BROWSER-${runSuffix}-${number}`,title:'SYNTHETIC TEST TERMS — NOT OFFICIAL',body,expected_current_version_id:current})})
 if(response.status!==201)throw Error('Synthetic publication failed '+response.status)
 return (await response.json()).data
}
async function ownStatus(){const r=await fetch(api+'/terms/status',{headers:{Authorization:'Bearer '+borrowerAccess}});if(!r.ok)throw Error('Status verification failed');return (await r.json()).data}
async function peerTab(route) {
 const {targetId}=await b.command('Target.createTarget',{url:'about:blank'}),{sessionId}=await b.command('Target.attachToTarget',{targetId,flatten:true})
 await b.command('Page.enable',{},sessionId);await b.command('Runtime.enable',{},sessionId)
 await navigate(route,sessionId);return {targetId,sessionId}
}
const peers=[]
try {
 await b.command('Network.enable')
 await b.command('Fetch.enable',{patterns:[{urlPattern:'*/api/v1/terms/*/accept',requestStage:'Request'}]})
 adminAccess=await directLogin('ADMIN')
 await navigate('/borrower/equipment');await until(`location.pathname==='/login'&&!!document.querySelector('#login-email')`)
 await check('Anonymous cannot enter terms or borrower workspace',`!document.querySelector('.borrower-shell')`)
 await login('BORROWER','/borrower/terms');await until(`document.body.innerText.includes('Official terms are not available yet')`)
 borrowerAccess=await b.evaluate(`(async()=>(${store}).accessToken)()`)
 await check('Missing official terms never imply consent',`!document.querySelector('[role="checkbox"]')&&document.body.innerText.includes('No terms acceptance has been recorded')`)
 for(const width of [320,390])for(const mode of ['light','dark'])await capture('terms-unpublished',width,mode)
 for(const route of ['/borrower/account','/borrower/borrowings']){await navigate(route);await until(`location.pathname===${JSON.stringify(route)}&&!!document.querySelector('.borrower-shell')`)}
 await check('Account and existing obligations remain accessible without consent',`location.pathname==='/borrower/borrowings'&&!!document.querySelector('.borrower-shell')`)
 await signOut();await navigate('/borrower/equipment');await until(`location.pathname==='/login'&&!!document.querySelector('#login-email')`)
 await publish(1)
 await login('BORROWER','/borrower/terms');await until(`document.querySelector('h1')?.textContent==='Review borrowing terms'&&!!document.querySelector('[role="checkbox"]')`)
 borrowerAccess=await b.evaluate(`(async()=>(${store}).accessToken)()`)
 await check('First-use checkbox starts unchecked; viewing is not acceptance',`document.querySelector('[role="checkbox"]').getAttribute('aria-checked')==='false'&&[...document.querySelectorAll('button')].find(e=>e.textContent==='Accept and continue').disabled`)
 if((await ownStatus()).state!=='required')throw Error('Reading fabricated acceptance')
 for(const width of [320,390,768,1280])for(const mode of ['light','dark'])await capture('terms-required',width,mode)
 await check('Full long text uses page scrolling and remains plain text',`document.querySelector('article').textContent.includes('END OF SYNTHETIC TEST DOCUMENT.')&&getComputedStyle(document.querySelector('article')).overflowY==='visible'&&document.documentElement.scrollHeight>innerHeight`)
 for(const mode of ['light','dark']){await viewport(390);await theme(mode);await b.evaluate('document.querySelector(".terms-consent").scrollIntoView({block:"center"})');await settle();const r=await b.command('Page.captureScreenshot',{format:'png'});const file='terms-reading-end-390-'+mode+'.png';await fs.writeFile(path.join(output,file),Buffer.from(r.data,'base64'));screens.push({file,width:390,theme:mode,pass:true})}
 const peer=await peerTab('/borrower/terms');peers.push(peer);await until(`!!document.querySelector('[role="checkbox"]')`,peer.sessionId)
 await b.evaluate(`document.querySelector('[role="checkbox"]').focus()`);await key(' ','Space',32)
 await check('Keyboard Space selects the explicit consent checkbox',`document.querySelector('[role="checkbox"]').getAttribute('aria-checked')==='true'`)
 holdAccept=true;await button('Accept and continue');await until(`document.body.innerText.includes('Recording acceptance')`)
 await check('Pending acceptance disables repeat submission',`[...document.querySelectorAll('button')].find(e=>e.textContent==='Recording acceptance…').disabled&&document.querySelector('[role="checkbox"]').getAttribute('aria-disabled')==='true'`)
 if(paused.length!==1)throw Error('Duplicate acceptance request')
 holdAccept=false;for(const request of paused.splice(0))await b.command('Fetch.continueRequest',{requestId:request.id},request.session)
 await until(`location.pathname==='/borrower/equipment'&&document.body.innerText.includes('This feature is not available yet')`)
 const receipt=(await ownStatus()).acceptance
 await evaluate(`document.querySelector('[role="checkbox"]').click();[...document.querySelectorAll('button')].find(e=>e.textContent==='Accept and continue').click()`,peer.sessionId)
 await until(`location.pathname==='/borrower/home'&&document.body.innerText.includes('This feature is not available yet')`,peer.sessionId)
 if(JSON.stringify((await ownStatus()).acceptance)!==JSON.stringify(receipt))throw Error('Peer duplicate replaced original evidence')
 checks.push({name:'Two tabs accepting the same version preserve one exact original receipt',pass:true})
 await navigate('/borrower/home');await until(`document.body.innerText.includes('This feature is not available yet')`)
 const third=await peerTab('/borrower/home');peers.push(third);await until(`document.body.innerText.includes('This feature is not available yet')`,third.sessionId)
 await check('Reload and third tab restore accepted current version',`location.pathname==='/borrower/home'&&!!document.querySelector('.borrower-shell')`)
 await navigate('/borrower/terms');await until(`document.body.innerText.includes('Acceptance recorded on')`);await capture('terms-accepted',390,'light');await capture('terms-accepted',390,'dark')
 await publish(2)
 await navigate('/borrower/home');await until(`document.querySelector('h1')?.textContent==='Updated terms'`)
 await check('A newer mandatory version requires another explicit acceptance',`document.querySelector('[role="checkbox"]').getAttribute('aria-checked')==='false'&&document.body.innerText.includes('The borrowing terms have been updated')`)
 await capture('terms-updated',390,'light');await capture('terms-updated',390,'dark')
 await click('[role="checkbox"]');const v3=await publish(3);await button('Accept and continue')
 await until(`document.body.innerText.includes(${JSON.stringify(v3.version)})&&document.querySelector('[role="checkbox"]')?.getAttribute('aria-checked')==='false'`)
 await check('Publication during review returns conflict and resets consent for new text',`document.querySelector('h1')?.textContent==='Updated terms'&&document.querySelector('[role="checkbox"]').getAttribute('aria-checked')==='false'`)
 if((await ownStatus()).state!=='updated')throw Error('Stale submission fabricated current consent')
 await capture('terms-stale',390,'light');await capture('terms-stale',390,'dark')
 await b.command('Network.setBlockedURLs',{urls:['*localhost:18084/api/v1/terms/*/accept*']});await click('[role="checkbox"]');await button('Accept and continue');await until(`document.body.innerText.includes('Your acceptance could not be confirmed')`)
 await check('Acceptance network failure never unlocks borrower request routes',`location.pathname==='/borrower/terms'&&document.querySelector('[role="checkbox"]')?.getAttribute('aria-checked')==='true'`)
 await capture('terms-accept-network-error',390,'dark');await b.command('Network.setBlockedURLs',{urls:[]});await button('Accept and continue');await until(`location.pathname==='/borrower/home'&&document.body.innerText.includes('This feature is not available yet')`)
 await evaluate(`window.__phase4bEvents=[];window.__phase4bChannel=new BroadcastChannel('elabtrack_v2.session.v1');window.__phase4bChannel.onmessage=e=>window.__phase4bEvents.push(e.data)`,peer.sessionId)
 await signOut();for(const p of peers)await until(`location.pathname==='/login'&&!!document.querySelector('#login-email')`,p.sessionId)
 await check('Cross-tab logout clears tokens, account and terms cache',`(async()=>{const s=${store};const {queryClient}=await import('/src/app/query-client.ts');return s.accessToken===null&&s.user===null&&!queryClient.getQueryCache().findAll({queryKey:['terms']}).length&&!document.querySelector('.borrower-shell')})()`,peer.sessionId)
 await check('Coordination messages contain only the existing non-secret lifecycle fields',`!!navigator.locks&&window.__phase4bEvents.some(e=>e.type==='logout')&&window.__phase4bEvents.every(e=>Object.keys(e).every(k=>['version','type','tabId','attemptId','epoch','timestamp','failure'].includes(k)))`,peer.sessionId)
 sessionEvents.push(...await evaluate('window.__phase4bEvents',peer.sessionId));for(const p of peers)await b.command('Target.closeTarget',{targetId:p.targetId});peers.length=0
 await login('BORROWER','/borrower/home');await until(`document.body.innerText.includes('This feature is not available yet')`)
 await check('Logout and subsequent real login preserve the accepted current version',`location.pathname==='/borrower/home'`)
 await b.command('Network.setBlockedURLs',{urls:['*localhost:18084/api/v1/terms/status*']});await navigate('/borrower/terms');await until(`document.body.innerText.includes('Unable to load borrowing terms')`)
 await check('Read failure remains a bounded manual retry with logout',`!![...document.querySelectorAll('button')].find(e=>e.textContent==='Try again')&&!![...document.querySelectorAll('button')].find(e=>e.textContent==='Sign Out')`)
 await capture('terms-read-network-error',390,'light');await b.command('Network.setBlockedURLs',{urls:[]});await button('Try again');await until(`document.body.innerText.includes('Acceptance recorded on')`)
 state('disable','BORROWER');await navigate('/borrower/account');await until(`location.pathname==='/login'&&!!document.querySelector('#login-email')`)
 await check('Disabling a previously authenticated borrower removes protected access',`(async()=>(${store}).accessToken===null)()`)
 state('enable','BORROWER')
 await fill('#login-email',fixtures.INACTIVE.email);await fill('#login-password',fixtures.INACTIVE.password);await button('Sign In');await until(`document.body.innerText.includes('Unable to sign in.')`)
 await check('Inactive login remains generic and clears the password',`document.querySelector('#login-password').value===''`)
 await wait(5200)
 for(const kind of ['STAFF','ADMIN']){await login(kind,'/staff/dashboard');await until(`document.body.innerText.includes('This feature is not available yet')`);await capture(kind.toLowerCase(),1280,'light');await capture(kind.toLowerCase(),1280,'dark');await check(kind+' operational navigation remains unaffected by borrower consent',`!!document.querySelector('.staff-shell')&&!document.querySelector('[role="checkbox"]')`);await navigate('/borrower/terms');await until(`document.body.innerText.includes('Access denied')`);await navigate('/staff/dashboard');await until(`document.body.innerText.includes('This feature is not available yet')`);await signOut()}
 await login('BORROWER','/borrower/home');await until(`document.body.innerText.includes('This feature is not available yet')`)
 const cookies=await b.command('Network.getCookies',{urls:[api+'/auth/refresh']});cookiePolicies=cookies.cookies.filter(c=>c.name==='elabtrack_v2_refresh').map(({name,httpOnly,secure,sameSite,path})=>({name,httpOnly,secure,sameSite,path}))
 await check('Access token remains memory only and refresh cookie is HttpOnly',`(async()=>{return !!(${store}).accessToken&&!document.cookie.includes('refresh')&&![...Object.keys(localStorage),...Object.keys(sessionStorage)].some(k=>/token|auth|session|jwt/i.test(k))})()`)
 await signOut()
 const report={generatedAt:new Date().toISOString(),browser:(await b.command('Browser.getVersion')).product,database:'elabtrack_v2_phase4b_test',syntheticTestTermsOnly:true,realAuthentication:true,cookiePolicies,sessionEvents,screens,checks,responses,applicationErrors:errors,applicationWarnings:warnings,networkFailureMethod:'Chromium URL blocking, no fabricated API success',pass:checks.every(c=>c.pass)&&screens.every(s=>s.pass)&&errors.length===0&&warnings.length===0&&cookiePolicies.length===1&&cookiePolicies[0].httpOnly}
 await fs.writeFile(path.join(output,'RESULTS.json'),JSON.stringify(report,null,2)+'\n');console.log(JSON.stringify({screens:screens.length,checks:checks.length,errors,warnings,pass:report.pass},null,2));if(!report.pass)process.exitCode=1
} finally {for(const p of peers)await b.command('Target.closeTarget',{targetId:p.targetId});await b.close()}
