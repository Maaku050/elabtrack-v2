// Actual Chromium and the existing isolated HTTP/PostgreSQL browser harness.
// Output contains method/path/status only; never cookies, headers or credentials.
import fs from 'node:fs/promises'
import path from 'node:path'
import assert from 'node:assert/strict'
import { execFileSync } from 'node:child_process'
import { browser } from './browser-cdp.mjs'

const mode = process.argv.includes('--before') ? 'before' : 'after'
const base = 'http://localhost:15175'
const fixtures = JSON.parse(await fs.readFile('/tmp/elabtrack-batch1-fixtures.json', 'utf8'))
const out = path.resolve('../docs/project/verification/request-efficiency')
const b = await browser(), requests = [], errors = [], checks = []
const wait = ms => new Promise(resolve => setTimeout(resolve, ms))
b.on(m => {
  if (m.method === 'Runtime.exceptionThrown') errors.push(m.params.exceptionDetails.text)
  if (m.method === 'Network.requestWillBeSent') {
    const { requestId, request, type } = m.params
    const url = new URL(request.url)
    if (url.origin === 'http://localhost:18085' || type === 'Document') requests.push({ requestId, session: m.sessionId, method: request.method, path: url.pathname, type, status: null })
  }
  if (m.method === 'Network.responseReceived') {
    const request = requests.findLast(r => r.requestId === m.params.requestId && r.session === m.sessionId)
    if (request) request.status = m.params.response.status
  }
})
async function until(expr) {
  for (let i = 0; i < 150; i++) {
    try { if (await b.evaluate(expr)) return } catch (e) { if (!/context|Cannot find|Inspected target/.test(e.message)) throw e }
    await wait(100)
  }
  throw Error('Timed out: ' + expr)
}
async function click(selector) {
  const r = await b.evaluate(`(()=>{const e=document.querySelector(${JSON.stringify(selector)});if(!e)throw Error('Missing control');e.scrollIntoView({block:'nearest'});const r=e.getBoundingClientRect();return{x:r.x+r.width/2,y:r.y+r.height/2}})()`)
  for (const type of ['mousePressed', 'mouseReleased']) await b.command('Input.dispatchMouseEvent', { type, ...r, button: 'left', clickCount: 1 })
}
async function fill(selector, text) {
  await click(selector)
  for (const type of ['keyDown', 'keyUp']) await b.command('Input.dispatchKeyEvent', { type, key: 'a', code: 'KeyA', windowsVirtualKeyCode: 65, modifiers: 2 })
  await b.command('Input.insertText', { text })
}
async function settled() {
  await until("(async()=>{const {queryClient}=await import('/src/app/query-client.ts');return queryClient.isFetching()===0})()")
  await wait(100)
}
async function route(target) {
  const title = { '/staff/dashboard': 'Dashboard', '/staff/inventory': 'Inventory', '/staff/borrowers': 'Borrowers', '/admin/administration': 'Administration' }[target]
  await click(`a[href="${target}"]`)
  await until(`location.pathname===${JSON.stringify(target)} && document.querySelector('#staff-content h1')?.textContent===${JSON.stringify(title)}`)
  await settled()
}
async function login() {
  if (!await b.evaluate("location.pathname==='/login'&&!!document.querySelector('#login-email')")) await b.command('Page.navigate', { url: base + '/login' })
  await until("!!document.querySelector('#login-email')&&!document.querySelector('#login-email').disabled")
  await fill('#login-email', fixtures.admin.email); await fill('#login-password', fixtures.admin.password)
  await click('button[type=submit]')
  await until("!document.querySelector('#login-email')&&!!document.querySelector('#staff-content h1')")
  await settled()
  if (!await b.evaluate("location.pathname==='/staff/dashboard'")) await route('/staff/dashboard')
}
function measure(start) {
  const rows = requests.slice(start).map(({ method, path, type, status }) => ({ method, path, type, status }))
  const counts = {}
  for (const r of rows) { const key = `${r.method} ${r.path} ${r.status}`; counts[key] = (counts[key] || 0) + 1 }
  return { total: rows.length, application: rows.filter(r => r.method !== 'OPTIONS' && r.type !== 'Document').length, preflights: rows.filter(r => r.method === 'OPTIONS').length, documentLoads: rows.filter(r => r.type === 'Document').length, counts, requests: rows }
}
function fixtureAccountState(role, active = true) {
  assert.ok(['ADMIN', 'STAFF'].includes(role))
  assert.match(fixtures.admin.id, /^[0-9a-f-]{36}$/i)
  assert.equal(execFileSync('docker', ['inspect', '--format', '{{index .Config.Labels "com.docker.compose.project"}}', 'elabtrack_v2_batch1_postgres'], { encoding: 'utf8' }).trim(), 'elabtrack_v2_batch1')
  // Only a named synthetic fixture in the disposable, guarded test database.
  execFileSync('docker', ['exec', '-i', 'elabtrack_v2_batch1_postgres', 'psql', '-X', '-U', 'elabtrack_migrator', '-d', 'elabtrack_v2_batch1_test', '-v', 'ON_ERROR_STOP=1'], { input: `UPDATE users SET role='${role}', is_active=${active ? 'TRUE' : 'FALSE'} WHERE id='${fixtures.admin.id}'::uuid;`, stdio: ['pipe', 'ignore', 'pipe'] })
}
async function securityChecks() {
  // Focus/reconnect deliberately revalidate authority. No timer or effect loop.
  for (const event of ['focus', 'reconnect']) {
    const start = requests.length
    await b.evaluate(event === 'focus'
      ? "window.dispatchEvent(new Event('visibilitychange'))"
      : "(()=>{window.dispatchEvent(new Event('offline'));window.dispatchEvent(new Event('online'))})()")
    await settled(); await wait(300)
    assert.equal(measure(start).counts['GET /api/v1/auth/me 200'], 1)
  }
  checks.push('focus and reconnect each perform one fresh account check, without a loop')

  fixtureAccountState('STAFF')
  await click('a[href="/admin/administration"]')
  await until("document.querySelector('#staff-content h1')?.textContent==='Access denied'")
  assert.equal(await b.evaluate("!!document.querySelector('a[href=\"/admin/administration\"]')"), false)
  assert.equal(await b.evaluate("(async()=>{const {queryClient}=await import('/src/app/query-client.ts');return queryClient.getQueriesData({queryKey:['terms']}).length})()"), 0)
  checks.push('server role change denies cached Admin destination and clears private policy data')
  fixtureAccountState('ADMIN')
  await route('/staff/dashboard'); await route('/admin/administration'); await route('/staff/dashboard')

  fixtureAccountState('ADMIN', false)
  await click('a[href="/staff/inventory"]')
  await until("!!document.querySelector('#login-email')")
  assert.equal(await b.evaluate("!!document.querySelector('#staff-content')"), false)
  assert.equal(await b.evaluate("(async()=>{const {queryClient,isAuthenticatedQuery}=await import('/src/app/query-client.ts');return queryClient.getQueryCache().findAll().filter(q=>isAuthenticatedQuery(q)&&q.state.data!==undefined).length})()"), 0)
  checks.push('server deactivation rejects current-account access, ends session and clears privileged content/cache')
  fixtureAccountState('ADMIN')
  await login()

  // A rejected access credential exercises the existing 401 → rotating
  // cookie refresh → one retry path, without reading or exporting any token.
  let start = requests.length
  assert.equal(await b.evaluate("(async()=>{const {useAuthStore}=await import('/src/stores/auth-store.ts');const {authApi}=await import('/src/features/auth/api/auth.api.ts');useAuthStore.setState({accessToken:'synthetic-rejected-access'});const account=await authApi.me();return account.role==='ADMIN'&&account.is_active})()"), true)
  assert.equal(measure(start).counts['GET /api/v1/auth/me 401'], 1)
  assert.equal(measure(start).counts['POST /api/v1/auth/refresh 200'], 1)
  assert.equal(measure(start).counts['GET /api/v1/auth/me 200'], 1)
  checks.push('rejected access credential triggers one server-verified refresh and one retry, without a loop')
  assert.equal(await b.evaluate("(async()=>{const {apiClient}=await import('/src/lib/api-client.ts');try{await apiClient.refreshSession();return true}catch{return false}})()"), true)
  await route('/staff/inventory'); await route('/staff/dashboard')
  start = requests.length
  const documentEpoch = await b.evaluate('performance.timeOrigin')
  await b.command('Page.reload')
  await until(`performance.timeOrigin!==${documentEpoch} && document.querySelector('#staff-content h1')?.textContent==='Dashboard'`)
  await settled()
  assert.equal(measure(start).counts['POST /api/v1/auth/refresh 200'], 1)
  checks.push('explicit refresh and deliberate full-document reload restore valid session with one refresh')

  // Fail exactly one refresh before it reaches the server; recovery retains
  // the HttpOnly cookie and only restores authority after server confirmation.
  let interrupted = false
  await b.command('Fetch.enable', { patterns: [{ urlPattern: '*18085/api/v1/auth/refresh', requestStage: 'Request' }] })
  b.on(m => {
    if (m.method === 'Fetch.requestPaused') {
      if (!interrupted) { interrupted = true; void b.command('Fetch.failRequest', { requestId: m.params.requestId, errorReason: 'InternetDisconnected' }) }
      else void b.command('Fetch.continueRequest', { requestId: m.params.requestId })
    }
  })
  await b.command('Page.reload')
  await until("!!document.querySelector('.session-recovery')")
  assert.equal(await b.evaluate("!!document.querySelector('#staff-content')"), false)
  await b.command('Fetch.disable')
  await click('.session-recovery button')
  await until("document.querySelector('#staff-content h1')?.textContent==='Dashboard'")
  await settled()
  checks.push('transient refresh network failure withholds protected content; manual recovery succeeds using the existing cookie')

  // A real peer tab shares only cookies and coordination events, never JWTs.
  const { targetId } = await b.command('Target.createTarget', { url: 'about:blank' })
  const { sessionId } = await b.command('Target.attachToTarget', { targetId, flatten: true })
  const peer = async expression => {
    const result = await b.command('Runtime.evaluate', { expression, awaitPromise: true, returnByValue: true }, sessionId)
    if (result.exceptionDetails) throw Error('Peer evaluation failed')
    return result.result.value
  }
  await b.command('Page.enable', {}, sessionId)
  await b.command('Runtime.enable', {}, sessionId)
  await b.command('Page.navigate', { url: base + '/staff/dashboard' }, sessionId)
  for (let i = 0; i < 100; i++) { if (await peer("document.querySelector('#staff-content h1')?.textContent==='Dashboard'")) break; await wait(100) }
  assert.equal(await peer("document.querySelector('#staff-content h1')?.textContent==='Dashboard'"), true)
  await click('.workspace-signout')
  await until("!!document.querySelector('#login-email')&&!document.querySelector('#login-email').disabled")
  await wait(300)
  assert.equal(await peer("!!document.querySelector('#login-email')&&!document.querySelector('#staff-content')"), true)
  assert.equal(await peer("(async()=>{const {queryClient,isAuthenticatedQuery}=await import('/src/app/query-client.ts');return queryClient.getQueryCache().findAll().filter(q=>isAuthenticatedQuery(q)&&q.state.data!==undefined).length})()"), 0)
  await b.command('Target.closeTarget', { targetId })
  checks.push('actual logout propagates across two Chromium tabs and clears peer protected cache/content')
  await login()
  assert.equal(await b.evaluate("fetch('http://localhost:18085/api/v1/auth/logout',{method:'POST',credentials:'include'}).then(r=>r.status)"), 204)
  assert.equal(await b.evaluate("(async()=>{const {apiClient}=await import('/src/lib/api-client.ts');try{await apiClient.refreshSession();return 200}catch(e){return e.status}})()"), 401)
  await until("!!document.querySelector('#login-email')")
  await login(); await click('.workspace-signout')
  await until("!!document.querySelector('#login-email')&&!document.querySelector('#login-email').disabled")
  checks.push('server-revoked refresh cannot recover privileges; new login and intentional logout still succeed')
}
try {
  await b.command('Network.enable')
  await b.command('Emulation.setDeviceMetricsOverride', { width: 1440, height: 900, deviceScaleFactor: 1, mobile: false })
  if (mode === 'after') fixtureAccountState('ADMIN')
  await login()
  await b.evaluate("void(window.__auditShell=[document.querySelector('header'),document.querySelector('#staff-content')])")
  let start = requests.length
  const began = Date.now()
  for (let i = 0; i < 3; i++) for (const target of ['/staff/inventory', '/staff/borrowers', '/staff/dashboard']) await route(target)
  const navigation = { sequence: 'Inventory → Borrowers → Dashboard, three cycles after login', elapsedMs: Date.now() - began, ...measure(start) }
  assert.equal(navigation.counts['GET /api/v1/auth/me 200'], 9)
  assert.equal(navigation.documentLoads, 0)
  assert.equal(await b.evaluate("window.__auditShell.every((e,i)=>!!e&&e===[document.querySelector('header'),document.querySelector('#staff-content')][i])"), true)
  assert.ok(!Object.keys(navigation.counts).some(k => k.includes('/auth/refresh')))
  checks.push('one authoritative current-account request per route; persistent shell; no document load or session restoration during navigation')
  start = requests.length
  for (let i = 0; i < 3; i++) { await route('/admin/administration'); await until("document.body.textContent.includes('Official FSMO terms await institutional approval')"); await route('/staff/dashboard') }
  const administration = { sequence: 'Administration → Dashboard, three cycles', ...measure(start) }
  assert.equal(administration.counts['GET /api/v1/auth/me 200'], 6)
  assert.equal(administration.counts['GET /api/v1/terms/current 503'], mode === 'before' ? 3 : 1)
  assert.equal(administration.documentLoads, 0)
  checks.push('unpublished terms shows expected policy message; measured remount behavior')
  if (mode === 'after') await securityChecks()
  assert.deepEqual(errors, [])
  await fs.mkdir(out, { recursive: true })
  await fs.writeFile(path.join(out, mode + '.json'), JSON.stringify({ date: new Date().toISOString(), mode, environment: 'isolated synthetic accounts, localhost:15175 → localhost:18085; no normal data or live mail', navigation, administration, checks, errors }, null, 2) + '\n')
  console.log(JSON.stringify({ result: 'PASS', mode, navigation: navigation.counts, administration: administration.counts }))
} finally { if (mode === 'after') fixtureAccountState('ADMIN'); await b.close() }
