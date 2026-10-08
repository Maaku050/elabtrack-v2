// Real local API/browser smoke. No API interception, fake role or persisted token.
// Creates one temporary unprivileged fixture; cleans only that fixture in finally.
import fs from 'node:fs/promises'
import path from 'node:path'
import { randomUUID } from 'node:crypto'
import { execFileSync } from 'node:child_process'
import { browser } from './browser-cdp.mjs'

const args = process.argv.slice(2)
const output = path.resolve(args.includes('--output') ? args[args.indexOf('--output') + 1] : '/tmp/elabtrack-local-environment-smoke')
const root = path.resolve(import.meta.dirname, '../..')
const config = Object.fromEntries((await fs.readFile(path.join(root, 'backend/.env'), 'utf8')).split('\n').filter(l => l && !l.startsWith('#') && l.includes('=')).map(l => [l.slice(0, l.indexOf('=')), l.slice(l.indexOf('=') + 1)]))
if (config.APP_ENV !== 'development' || config.DB_NAME !== 'elabtrack_v2' || config.DB_USER !== 'elabtrack_runtime') throw Error('Refusing non-development runtime smoke')
await fs.mkdir(output, { recursive: true })
const b = await browser(), checks = [], cases = [], responses = [], exceptions = [], warnings = []
let fixtureId, peerSession
b.on(m => {
  if (m.method === 'Runtime.exceptionThrown') exceptions.push('Unexpected JavaScript exception')
  if (m.method === 'Runtime.consoleAPICalled' && m.params.type === 'error') exceptions.push('Unexpected application console error')
  if (m.method === 'Runtime.consoleAPICalled' && m.params.type === 'warning') warnings.push('Unexpected browser warning')
  if (m.method === 'Network.responseReceived' && m.params.response.url.includes('/api/v1/')) responses.push({ url: m.params.response.url, status: m.params.response.status })
})
await b.command('Network.enable')
const pause = ms => new Promise(resolve => setTimeout(resolve, ms))
async function evaluate(expression, session) {
  const r = await b.command('Runtime.evaluate', { expression, awaitPromise: true, returnByValue: true }, session)
  if (r.exceptionDetails) throw Error('Browser evaluation failed (request data omitted)')
  return r.result.value
}
const state = session => evaluate(`import('/src/stores/auth-store.ts').then(({useAuthStore})=>{const s=useAuthStore.getState();return {status:s.status,authenticated:s.isAuthenticated,id:s.user?.id,role:s.user?.role,accessPresent:!!s.accessToken}})`, session)
async function until(expression, session) {
  for (let i = 0; i < 100; i++) { if (await evaluate(expression, session)) return; await pause(100) }
  throw Error('Browser condition timed out')
}
function check(name, pass) { if (!pass) throw Error(name); checks.push({ name, pass }); console.log('PASS', name) }
async function newDocument(operation, session) {
  const marker = randomUUID()
  await evaluate(`window.__localReadinessNavigation=${JSON.stringify(marker)}`, session)
  await operation()
  await until(`window.__localReadinessNavigation!==${JSON.stringify(marker)}&&document.readyState!=='loading'`, session)
}
const reload = session => newDocument(() => b.command('Page.reload', {}, session), session)
async function navigate(route) {
  await newDocument(() => b.command('Page.navigate', { url: 'http://localhost:5173' + route }))
  await until(`!!document.querySelector('h1')`)
  await until(`import('/src/stores/auth-store.ts').then(m=>!['idle','bootstrapping'].includes(m.useAuthStore.getState().status))`)
}
function fixtureSQL(action) {
  const r = JSON.parse(execFileSync('python3', [path.join(root, 'integration/local-readiness-sql.py')], { input: JSON.stringify({ action, id: fixtureId }), encoding: 'utf8' }))
  if (!r.ok) throw Error('Fixture SQL action failed')
}
const input = { email: `local-readiness-${randomUUID()}@example.invalid`, name: 'Synthetic local readiness fixture', password: randomUUID() + randomUUID().slice(0, 20) }
const authenticate = () => evaluate(`import('/src/features/auth/api/auth.api.ts').then(m=>m.authApi.login(${JSON.stringify({ email: input.email, password: input.password })})).then(r=>({id:r.user.id,role:r.user.role}))`)
try {
  await navigate('/')
  check('Anonymous visitor is unauthenticated; restored connectivity removes error banner', (await state()).status === 'unauthenticated' && !await evaluate(`document.body.innerText.includes('Unable to restore your session.')`))
  const health = await evaluate(`fetch('http://localhost:8080/api/v1/ready',{credentials:'include'}).then(async r=>({status:r.status,body:await r.json(),id:!!r.headers.get('X-Request-ID')}))`)
  check('Real credentialed API readiness and exposed request ID', health.status === 200 && health.body.success === true && health.id)
  const created = await evaluate(`import('/src/features/auth/api/auth.api.ts').then(m=>m.authApi.register(${JSON.stringify(input)})).then(r=>({id:r.user.id,role:r.user.role}))`)
  fixtureId = created.id
  check('Synthetic fixture receives only the existing user role', created.role === 'user' && (await state()).authenticated)
  const cookie = (await b.command('Network.getCookies', { urls: ['http://localhost:8080/api/v1/auth/refresh'] })).cookies.find(c => c.name === 'elabtrack_v2_refresh')
  check('Host-only HttpOnly Lax refresh cookie; local HTTP Secure=false', cookie?.httpOnly && cookie.path === '/api/v1/auth' && cookie.sameSite === 'Lax' && !cookie.secure && cookie.domain === 'localhost')
  check('Access remains memory-only and refresh is hidden from document.cookie', await evaluate(`import('/src/stores/auth-store.ts').then(m=>{const s=m.useAuthStore.getState();return !!s.accessToken&&!document.cookie.includes('elabtrack_v2_refresh')&&[localStorage,sessionStorage].every(storage=>Object.keys(storage).every(k=>!['access_token','refresh_token','auth','auth-storage','token'].includes(k)&&!storage.getItem(k)?.includes(s.accessToken)))})`))
  await reload()
  await until(`import('/src/stores/auth-store.ts').then(m=>m.useAuthStore.getState().status==='authenticated')`)
  const me = await evaluate(`import('/src/features/auth/api/auth.api.ts').then(m=>m.authApi.me()).then(u=>({id:u.id,role:u.role}))`)
  check('Valid cookie restores a fresh document and authorizes current-account read', me.id === fixtureId && me.role === 'user')

  const peer = await b.command('Target.createTarget', { url: 'http://localhost:5173/' })
  peerSession = (await b.command('Target.attachToTarget', { targetId: peer.targetId, flatten: true })).sessionId
  await b.command('Runtime.enable', {}, peerSession)
  await b.command('Page.enable', {}, peerSession)
  await b.command('Network.enable', {}, peerSession)
  await until(`import('/src/stores/auth-store.ts').then(m=>m.useAuthStore.getState().status==='authenticated')`, peerSession)
  await Promise.all([reload(), reload(peerSession)])
  for (const session of [undefined, peerSession]) await until(`import('/src/stores/auth-store.ts').then(m=>m.useAuthStore.getState().status==='authenticated')`, session)
  check('Two simultaneous real tab reloads restore their own memory sessions', (await state()).id === fixtureId && (await state(peerSession)).id === fixtureId)
  await evaluate(`import('/src/features/auth/api/auth.api.ts').then(m=>m.authApi.logout())`)
  await until(`import('/src/stores/auth-store.ts').then(m=>m.useAuthStore.getState().status==='unauthenticated')`, peerSession)
  check('Real logout clears both tabs without shared credentials', !(await state()).accessPresent && !(await state(peerSession)).accessPresent)
  await b.command('Target.closeTarget', { targetId: peer.targetId })
  peerSession = undefined

  await authenticate()
  const revoked = (await b.command('Network.getCookies', { urls: ['http://localhost:8080/api/v1/auth/refresh'] })).cookies.find(c => c.name === 'elabtrack_v2_refresh')
  await evaluate(`import('/src/features/auth/api/auth.api.ts').then(m=>m.authApi.logout())`)
  await b.command('Network.setCookie', { name: revoked.name, value: revoked.value, domain: revoked.domain, path: revoked.path, httpOnly: true, secure: false, sameSite: 'Lax', expires: revoked.expires })
  await navigate('/')
  const retained = (await b.command('Network.getCookies', { urls: ['http://localhost:8080/api/v1/auth/refresh'] })).cookies.find(c => c.name === revoked.name)
  check('Revoked refresh is safely anonymous and does not clear possibly newer cookies', (await state()).status === 'unauthenticated' && retained?.value === revoked.value && !await evaluate(`document.body.innerText.includes('Unable to restore your session.')`))
  await authenticate()
  fixtureSQL('expire')
  await navigate('/')
  check('Expired server session is safely anonymous with no connection-error banner', (await state()).status === 'unauthenticated' && !await evaluate(`document.body.innerText.includes('Unable to restore your session.')`))
  await b.command('Network.clearBrowserCookies')

  const screens = [
    ['B01', '/__preview/borrower/home', 'Good morning, Alex!', [320, 390, 768, 1280]],
    ['B02', '/__preview/borrower/equipment', 'Equipment Catalog', [320, 390, 768, 1280]],
    ['S01-dashboard', '/__preview/staff/dashboard', 'Dashboard', [1024, 1440]],
    ['S01-pending', '/__preview/staff/requests', 'Pending Requests', [1024, 1440]],
  ]
  for (const [id, route, heading, widths] of screens) for (const width of widths) {
    await b.command('Emulation.setDeviceMetricsOverride', { width, height: 900, deviceScaleFactor: 1, mobile: width < 768 })
    await navigate(route)
    await until(`document.querySelector('h1')?.textContent===${JSON.stringify(heading)}`)
    for (const theme of ['light', 'dark']) {
      const dark = await evaluate(`document.documentElement.classList.contains('dark')`)
      if (dark !== (theme === 'dark')) await evaluate(`document.querySelector('[aria-label="Switch to ${theme} theme"]').click()`)
      await pause(100)
      const data = await evaluate(`({width:innerWidth,pageWidth:document.documentElement.scrollWidth,dark:document.documentElement.classList.contains('dark'),errorBanner:document.body.innerText.includes('Unable to restore your session.'),synthetic:document.body.innerText.includes('Development visual preview'),nav:!!document.querySelector('nav')})`)
      check(`${id} ${width} ${theme} live render`, data.pageWidth <= data.width && data.dark === (theme === 'dark') && !data.errorBanner && data.synthetic && data.nav && !(await state()).authenticated)
      const capture = await b.command('Page.captureScreenshot', { format: 'png', captureBeyondViewport: false })
      const file = `${id}-${width}-${theme}.png`
      await fs.writeFile(path.join(output, file), Buffer.from(capture.data, 'base64'))
      cases.push({ id, route, width, theme, file, pass: true })
    }
  }
  await b.command('Emulation.setDeviceMetricsOverride', { width: 390, height: 900, deviceScaleFactor: 1, mobile: true })
  await navigate('/__preview/borrower/home')
  await evaluate(`document.querySelector('.borrower-bottom-nav a[href="/__preview/borrower/equipment"]').click()`)
  await until(`document.querySelector('h1')?.textContent==='Equipment Catalog'`)
  check('Borrower mobile navigation opens the actual catalog route', await evaluate(`location.pathname==='/__preview/borrower/equipment'&&document.querySelectorAll('.borrower-bottom-nav a, .borrower-bottom-nav button').length===4`))
  await b.command('Emulation.setDeviceMetricsOverride', { width: 1024, height: 900, deviceScaleFactor: 1, mobile: false })
  await navigate('/__preview/staff/requests')
  await evaluate(`document.querySelector('[aria-label="Pending requests table; scroll horizontally for additional columns"]').focus()`)
  await b.command('Input.dispatchKeyEvent', { type: 'keyDown', key: 'ArrowRight', code: 'ArrowRight', windowsVirtualKeyCode: 39 })
  await b.command('Input.dispatchKeyEvent', { type: 'keyUp', key: 'ArrowRight', code: 'ArrowRight', windowsVirtualKeyCode: 39 })
  await pause(200)
  check('Staff tablet table is keyboard scrollable inside its named region', await evaluate(`document.querySelector('[aria-label="Pending requests table; scroll horizontally for additional columns"]').scrollLeft>0&&document.documentElement.scrollWidth<=innerWidth`))
  const denied = await fetch('http://localhost:8080/api/v1/auth/refresh', { method: 'POST', headers: { Origin: 'http://localhost:5174', 'Content-Type': 'application/json' } })
  check('Unconfigured Origin is denied without credentialed CORS', denied.status === 403 && !denied.headers.has('Access-Control-Allow-Origin'))
  const missing = await fetch('http://localhost:8080/api/v1/auth/refresh', { method: 'POST' })
  check('Cookie-changing request without Origin remains denied', missing.status === 403)
  check('No unexpected JavaScript exceptions or browser warnings', exceptions.length === 0 && warnings.length === 0)
  check('Expected anonymous refresh 401 is observed, not connection refusal', responses.some(r => r.url.endsWith('/auth/refresh') && r.status === 401))
  const report = { generatedAt: new Date().toISOString(), browser: (await b.command('Browser.getVersion')).product, liveApi: true, apiInterception: false, checks, cases, responses, exceptions, warnings, fixtureRole: 'user', fixtureCleanup: false, pass: true }
  fixtureSQL('cleanup')
  fixtureId = undefined
  report.fixtureCleanup = true
  await fs.writeFile(path.join(output, 'RESULTS.json'), JSON.stringify(report, null, 2) + '\n')
  console.log(JSON.stringify({ checks: checks.length, cases: cases.length, fixtureCleanup: true, pass: true }))
} finally {
  try { if (fixtureId) fixtureSQL('cleanup') } finally { await b.close() }
}
