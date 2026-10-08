// Real forms, PostgreSQL-backed API and Chromium; no fabricated auth responses.
import fs from 'node:fs/promises'
import path from 'node:path'
import { execFileSync } from 'node:child_process'
import { browser } from './browser-cdp.mjs'
const base = 'http://localhost:5173'
const output = path.resolve('../docs/ux/verification/phase4a')
const fixtures = JSON.parse(await fs.readFile('/tmp/elabtrack-phase4a-fixtures.json', 'utf8'))
await fs.mkdir(output, { recursive: true })
const b = await browser(), checks = [], screens = [], errors = [], warnings = [], responses = []
const sessionEvents = []
let holdRefresh = false, paused = [], cookiePolicies = []
b.on(async message => {
  if (message.method === 'Runtime.exceptionThrown') errors.push({ text: message.params.exceptionDetails.text })
  if (message.method === 'Runtime.consoleAPICalled' && ['error', 'warning'].includes(message.params.type)) {
    const destination = message.params.type === 'error' ? errors : warnings
    destination.push({ text: message.params.args.map(arg => arg.value ?? arg.type).join(' ') })
  }
  if (message.method === 'Network.responseReceived' && message.params.response.url.includes('/api/v1/')) responses.push({ path: new URL(message.params.response.url).pathname, status: message.params.response.status })
  if (message.method === 'Fetch.requestPaused') {
    if (holdRefresh && message.params.request.url.endsWith('/auth/refresh') && message.params.request.method === 'POST') paused.push({ id: message.params.requestId, session: message.sessionId })
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
  const marker = `navigation-${Date.now()}-${Math.random()}`
  await evaluate(`document.documentElement.dataset.phase4a=${JSON.stringify(marker)}`, session)
  await b.command('Page.navigate', { url: base + route }, session)
  await until(`document.documentElement.dataset.phase4a !== ${JSON.stringify(marker)} && document.readyState==='complete'`, session)
}
async function click(selector) {
  const rect = await b.evaluate(`(()=>{const e=document.querySelector(${JSON.stringify(selector)});if(!e)throw Error('Missing control');e.scrollIntoView({block:'center'});const r=e.getBoundingClientRect();return {x:r.x+r.width/2,y:r.y+r.height/2};})()`)
  await b.command('Input.dispatchMouseEvent', { type: 'mousePressed', ...rect, button: 'left', clickCount: 1 })
  await b.command('Input.dispatchMouseEvent', { type: 'mouseReleased', ...rect, button: 'left', clickCount: 1 })
}
async function button(text) { await b.evaluate(`(()=>{const e=[...document.querySelectorAll('button')].find(e=>e.textContent.trim()===${JSON.stringify(text)});if(!e)throw Error('Missing named button');e.dataset.phase4aButton='true';})()`); await click('[data-phase4a-button="true"]'); await b.evaluate(`document.querySelector('[data-phase4a-button]')?.removeAttribute('data-phase4a-button')`) }
async function key(key, code, number) { for (const type of ['keyDown','keyUp']) await b.command('Input.dispatchKeyEvent', { type, key, code, windowsVirtualKeyCode: number }) }
async function fill(selector, value) {
  await click(selector); await b.command('Input.dispatchKeyEvent', { type: 'keyDown', key: 'a', code: 'KeyA', windowsVirtualKeyCode: 65, modifiers: 2 }); await b.command('Input.dispatchKeyEvent', { type: 'keyUp', key: 'a', code: 'KeyA', windowsVirtualKeyCode: 65, modifiers: 2 }); await b.command('Input.insertText', { text: value })
}
async function login(kind, destination) {
  const fixture = fixtures[kind]
  await until(`!!document.querySelector('#login-email') && !document.querySelector('#login-email').disabled && !![...document.querySelectorAll('button')].find(e=>e.textContent.trim()==='Sign In' && !e.disabled)`)
  await fill('#login-email', fixture.email); await fill('#login-password', fixture.password); await button('Sign In')
  await until(`location.pathname===${JSON.stringify(destination)} && document.body.innerText.includes('This feature is not available yet')`)
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
  const geometry = await b.evaluate(`(()=>{const visible=e=>e.getBoundingClientRect().width>0&&e.getBoundingClientRect().height>0&&getComputedStyle(e).visibility!=='hidden';return {width:innerWidth,pageWidth:document.documentElement.scrollWidth,clientWidth:document.documentElement.clientWidth,theme:document.documentElement.classList.contains('dark')?'dark':'light',smallTargets:[...document.querySelectorAll('button,a,input')].filter(visible).filter(e=>!e.classList.contains('skip-link')).map(e=>({name:e.getAttribute('aria-label')||e.textContent.trim(),width:e.getBoundingClientRect().width,height:e.getBoundingClientRect().height})).filter(e=>e.width<43.9||e.height<43.9),h1:document.querySelector('h1')?.textContent};})()`)
  const result = await b.command('Page.captureScreenshot', { format: 'png', captureBeyondViewport: false })
  const file = `${name}-${width}-${mode}.png`; await fs.writeFile(path.join(output, file), Buffer.from(result.data, 'base64'))
  console.log('CAPTURE', file)
  screens.push({ file, ...geometry, pass: geometry.pageWidth <= geometry.clientWidth && geometry.smallTargets.length === 0 && geometry.theme === mode })
  if (!screens.at(-1).pass) throw new Error(`Responsive geometry failed: ${file}: ${JSON.stringify(geometry)}`)
}
async function check(name, expression, session) { const pass = !!await evaluate(expression, session); console.log(pass ? 'PASS' : 'FAIL', name); checks.push({ name, pass }); if (!pass) throw new Error(`Check failed: ${name}`) }
function state(action, kind, role) { execFileSync('python3', ['../integration/phase4a-fixture-state.py', action, kind, ...(role ? [role] : [])], { stdio: 'pipe' }) }
const store = `(await import('/src/stores/auth-store.ts')).useAuthStore.getState()`
try {
  await b.command('Network.enable')
  await b.command('Fetch.enable', { patterns: [{ urlPattern: '*/api/v1/auth/refresh', requestStage: 'Request' }] })
  await navigate('/borrower/equipment'); await until(`location.pathname==='/login' && !!document.querySelector('#login-email')`)
  await check('Anonymous deep link denies protected content', `!document.body.innerText.includes('This feature is not available yet') && !document.querySelector('.borrower-bottom-nav')`)
  for (const width of [320,390,1280]) for (const mode of ['light','dark']) await capture('login', width, mode)
  await viewport(390); await theme('light'); await button('Sign In'); await until(`!!document.querySelector('#email-error')`)
  await check('Empty form validates without API submission', `document.querySelector('#login-email').getAttribute('aria-invalid')==='true'`)
  await fill('#login-email', fixtures.BORROWER.email); await fill('#login-password', 'deliberately-wrong-input'); await button('Sign In'); await until(`document.body.innerText.includes('Unable to sign in.')`)
  await check('Invalid login is generic, clears password and restores keyboard focus', `document.querySelector('#login-password').value==='' && document.activeElement.id==='login-password'`)
  await fill('#login-password', fixtures.BORROWER.password); await click('[aria-label="Show password"]'); await check('Password visibility is accessible', `document.querySelector('#login-password').type==='text' && document.querySelector('[aria-label="Hide password"]').getAttribute('aria-pressed')==='true'`); await click('[aria-label="Hide password"]')
  await button('Sign In'); await until(`location.pathname==='/borrower/equipment' && document.body.innerText.includes('This feature is not available yet')`)
  await check('Allowed deep link restored after real login', `location.pathname==='/borrower/equipment'`)
  for (const width of [320,390,1280]) for (const mode of ['light','dark']) await capture('borrower', width, mode)
  await check('No synthetic business records in authenticated Borrower shell', `!document.body.innerText.includes('Alex') && !document.body.innerText.includes('PREVIEW-') && !!document.querySelector('nav[aria-label="Borrower navigation"]')`)
  await navigate('/borrower/account'); await until(`document.body.innerText.includes('Account details')`); await capture('account',390,'light')
  await check('Account is real safe read-only metadata', `document.body.innerText.includes('Synthetic BORROWER') && !document.querySelector('input')`)
  await navigate('/staff/dashboard'); await until(`document.body.innerText.includes('Access denied')`); await capture('forbidden',390,'light'); await capture('forbidden',390,'dark')
  await check('Borrower cannot enter Staff shell', `!document.querySelector('.staff-shell')`)
  await navigate('/borrower/home'); await until(`document.body.innerText.includes('This feature is not available yet')`)
  const cookies = await b.command('Network.getCookies', { urls: ['http://localhost:8080/api/v1/auth/refresh'] })
  cookiePolicies = cookies.cookies.filter(c=>c.name==='elabtrack_v2_refresh').map(({name,httpOnly,secure,sameSite,path})=>({name,httpOnly,secure,sameSite,path}))
  await check('Access only in memory; refresh cookie hidden from JavaScript', `(async()=>{const s=${store};return !!s.accessToken && s.user.role==='BORROWER' && !document.cookie.includes('refresh') && ![...Object.keys(localStorage),...Object.keys(sessionStorage)].some(k=>/token|auth|session|jwt/i.test(k));})()`)
  await key('Tab','Tab',9); await check('Keyboard focus reaches a named interactive element', `['A','BUTTON','INPUT'].includes(document.activeElement.tagName)`)
  // Hold the real refresh request, then forward it unchanged: no fake response.
  holdRefresh=true; await navigate('/borrower/home'); await until(`document.body.innerText.includes('Restoring session')`)
  await check('Reload never flashes protected content while refresh is pending', `!document.querySelector('.borrower-shell') && !document.querySelector('.staff-shell')`)
  holdRefresh=false; for(const request of paused.splice(0)) await b.command('Fetch.continueRequest',{requestId:request.id},request.session)
  await until(`document.body.innerText.includes('This feature is not available yet')`)
  await check('Real refresh restores Borrower session after reload', `(async()=>{return (${store}).user?.role==='BORROWER'})()`)
  // Separate page context: shared cookie, Web Lock and non-secret lifecycle.
  const {targetId}=await b.command('Target.createTarget',{url:'about:blank'}); const {sessionId:peer}=await b.command('Target.attachToTarget',{targetId,flatten:true})
  await b.command('Page.enable',{},peer); await b.command('Runtime.enable',{},peer)
  await b.command('Page.navigate',{url:base+'/borrower/home'},peer); await until(`document.body.innerText.includes('This feature is not available yet')`,peer)
  await evaluate(`window.__phase4aEvents=[];window.__phase4aChannel=new BroadcastChannel('elabtrack_v2.session.v1');window.__phase4aChannel.onmessage=e=>window.__phase4aEvents.push(e.data)`,peer)
  await signOut(); await until(`location.pathname==='/login' && !!document.querySelector('#login-email')`,peer)
  await check('Cross-tab logout clears peer memory and private UI', `(async()=>{const s=${store};return s.accessToken===null&&s.user===null&&!document.querySelector('.borrower-shell')})()`,peer)
  await check('Native coordination uses non-secret lifecycle messages', `!!navigator.locks && window.__phase4aEvents.some(e=>e.type==='logout') && window.__phase4aEvents.every(e=>Object.keys(e).every(k=>['version','type','tabId','attemptId','epoch','timestamp','failure'].includes(k)))`,peer)
  sessionEvents.push(...await evaluate('window.__phase4aEvents',peer)); await b.command('Target.closeTarget',{targetId})
  await b.evaluate('history.back()'); await until(`!!document.querySelector('#login-email') && !document.querySelector('.borrower-shell')`); await check('Back navigation after logout cannot expose protected content', `!document.querySelector('.borrower-shell') && !document.querySelector('.staff-shell')`)
  await b.command('Page.reload'); await until(`!!document.querySelector('#login-email')`); await check('Logout stays anonymous after reload', `(async()=>{return (${store}).accessToken===null})()`)
  await login('STAFF','/staff/dashboard')
  for (const width of [1024,1440]) for(const mode of ['light','dark']) await capture('staff',width,mode)
  await check('Sidebar sign-out remains readable in both themes', `(()=>{const button=document.querySelector('.staff-sidebar .workspace-signout');const style=getComputedStyle(button);return style.backgroundColor==='rgba(0, 0, 0, 0)' && style.color===getComputedStyle(document.querySelector('.staff-sidebar')).color})()`)
  await click('[aria-label="Toggle navigation"]'); await settle()
  await check('Collapsed production sidebar retains accessible sign out without overflow', `document.querySelector('.staff-sidebar').getBoundingClientRect().width===76 && document.querySelector('.staff-sidebar .workspace-signout').getAttribute('aria-label')==='Sign Out' && document.documentElement.scrollWidth<=document.documentElement.clientWidth`)
  await click('[aria-label="Toggle navigation"]'); await settle()
  await check('Staff has operational navigation without Admin additions', `!!document.querySelector('nav[aria-label="Staff navigation"]') && !document.body.innerText.includes('Administration')`)
  await navigate('/admin/reports'); await until(`document.body.innerText.includes('Access denied')`); await capture('staff-forbidden',1024,'dark')
  await navigate('/staff/dashboard'); await until(`document.body.innerText.includes('This feature is not available yet')`); await signOut()
  await login('ADMIN','/staff/dashboard'); for(const mode of ['light','dark']) await capture('admin',1440,mode)
  await navigate('/admin/reports'); await until(`document.querySelector('h1')?.textContent==='Reports' && document.body.innerText.includes('This feature is not available yet')`)
  await check('Admin inherits operations and enters administrative destination', `!!document.querySelector('nav[aria-label="Admin navigation"]') && document.body.innerText.includes('Administration')`)
  const oldAccess = await b.evaluate(`(async()=>(${store}).accessToken)()`)
  state('role','ADMIN','BORROWER')
  await b.evaluate(`(async()=>{const {queryClient}=await import('/src/app/query-client.ts');queryClient.setQueryData(['users','directory'],{private:true});})()`)
  await b.evaluate(`document.querySelector('a[href="/staff/dashboard"]').click()`)
  await until(`document.body.innerText.includes('Access denied')`)
  await check('Database demotion clears old private cache and blocks old workspace', `(async()=>{const s=${store};const {queryClient}=await import('/src/app/query-client.ts');return s.user.role==='BORROWER'&&!queryClient.getQueryData(['users','directory'])&&!document.querySelector('.staff-shell')})()`)
  await check('Previously issued Admin JWT cannot authorize current demoted account', `(async()=>{return (await fetch('http://localhost:8080/api/v1/users/',{headers:{Authorization:'Bearer '+${JSON.stringify(oldAccess)}}})).status===403})()`)
  state('role','ADMIN','ADMIN')
  await navigate('/staff/dashboard'); await until(`document.body.innerText.includes('This feature is not available yet')`)
  state('disable','ADMIN'); await navigate('/staff/account'); await until(`location.pathname==='/login' && !!document.querySelector('#login-email')`)
  await check('Disabled account after issuance loses protected browser access', `(async()=>{return (${store}).accessToken===null&&!document.querySelector('.staff-shell')})()`)
  state('enable','ADMIN')
  await fill('#login-email',fixtures.INACTIVE.email); await fill('#login-password',fixtures.INACTIVE.password); await button('Sign In'); await until(`document.body.innerText.includes('Unable to sign in.')`)
  await check('Inactive login uses generic non-enumerating feedback', `!document.body.innerText.includes('disabled by') && document.querySelector('#login-password').value===''`)
  await wait(5200)
  await login('BORROWER','/borrower/home')
  await b.command('Network.emulateNetworkConditions',{offline:true,latency:0,downloadThroughput:-1,uploadThroughput:-1})
  await signOut(); await until(`document.body.innerText.includes('Server sign-out could not be confirmed')`)
  await capture('logout-network-error',390,'dark')
  await check('Failed sign out still removes local private authority', `(async()=>{return (${store}).accessToken===null && !document.querySelector('.borrower-shell')})()`)
  await b.command('Network.emulateNetworkConditions',{offline:false,latency:0,downloadThroughput:-1,uploadThroughput:-1}); await button('Retry sign out'); await until(`!document.body.innerText.includes('Server sign-out could not be confirmed')`)
  await login('BORROWER','/borrower/home'); state('revoke','BORROWER'); await b.command('Page.reload'); await until(`!!document.querySelector('#login-email')`)
  await check('Server-revoked session cannot restore after reload', `(async()=>{return (${store}).accessToken===null&&!document.querySelector('.borrower-shell')})()`)
  // Fail only API transport so Vite modules remain available for honest recovery.
  await b.command('Network.setBlockedURLs',{urls:['*localhost:8080/api/v1/auth/refresh*']}); await b.command('Page.reload'); await until(`document.body.innerText.includes('Unable to restore your session')`)
  const refreshBefore=responses.filter(r=>r.path.endsWith('/auth/refresh')).length; await wait(1200)
  await check('Transient bootstrap failure displays a bounded manual retry', `!![...document.querySelectorAll('button')].find(e=>e.textContent==='Retry session') && !document.querySelector('.borrower-shell')`)
  await capture('session-network-error',390,'light')
  await b.command('Network.setBlockedURLs',{urls:[]}); await button('Retry session'); await until(`!!document.querySelector('#login-email') && !document.body.innerText.includes('Unable to restore your session')`)
  checks.push({name:'No automatic refresh loop during offline observation',pass:responses.filter(r=>r.path.endsWith('/auth/refresh')).length<=refreshBefore+1})
  await check('Public self-registration remains unavailable', `(async()=>{return (await fetch('http://localhost:8080/api/v1/auth/register',{method:'POST',headers:{'Content-Type':'application/json'},body:'{}'})).status===404})()`)
  const report={generatedAt:new Date().toISOString(),browser:(await b.command('Browser.getVersion')).product,database:'elabtrack_v2_phase4a_test',realAuthentication:true,cookiePolicies,sessionEvents,screens,checks,responses,applicationErrors:errors,applicationWarnings:warnings,networkFailureMethod:'Chromium offline and URL blocking; original API responses never fabricated',pass:checks.every(c=>c.pass)&&screens.every(s=>s.pass)&&errors.length===0&&warnings.length===0&&cookiePolicies.length===1&&cookiePolicies[0].httpOnly}
  await fs.writeFile(path.join(output,'RESULTS.json'),JSON.stringify(report,null,2)+'\n')
  console.log(JSON.stringify({screens:screens.length,checks,applicationErrors:errors,applicationWarnings:warnings,pass:report.pass},null,2))
  if(!report.pass)process.exitCode=1
} finally { await b.close() }
