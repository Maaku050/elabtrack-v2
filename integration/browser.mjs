import assert from 'node:assert/strict'
import { createServer } from 'node:http'
import { readFile, writeFile } from 'node:fs/promises'
import { randomUUID, createHash } from 'node:crypto'
import { execFileSync } from 'node:child_process'
import { pathToFileURL } from 'node:url'

// Use a small external Playwright installation, not a production dependency.
const { chromium } = await import(pathToFileURL(process.env.PLAYWRIGHT_MODULE).href)
const base = 'http://localhost:15173'
const driver = base + '/phase1g.local/index.html'
const email = `phase1g-browser-${randomUUID()}@example.invalid`
const password = randomUUID() + '-Synthetic-Only'
const secrets = JSON.parse(await readFile('/tmp/elabtrack-phase1g-secrets.json', 'utf8'))
secrets.push(password)
const evidence = { browser: 'Chromium', checks: [], correlation: [], crossTab: {} }
const check = (name, detail = true) => { evidence.checks.push({ name, detail }); console.log('PASS ' + name) }
const sql = (query) => execFileSync('psql', ['-h','127.0.0.1','-p','15432','-U','elabtrack_runtime','-d','elabtrack_v2_integration','-At','-c',query], { env: { ...process.env, PGPASSWORD: process.env.DB_PASSWORD }, encoding: 'utf8' }).trim()
const cookieName = 'elabtrack_v2_refresh'
const browser = await chromium.launch({ headless: true })
const servers = []
let accountID
try {
  const context = await browser.newContext()
  const page = await context.newPage()
  const cspErrors = []
  const counts = new Map()
  const count = (p) => counts.get(p) ?? 0
  function observe(p) {
    p.on('console', (msg) => { if (/content.security.policy|violates.*directive|refused to (load|execute|apply|connect)/i.test(msg.text())) cspErrors.push(msg.text()) })
    p.on('request', req => {
      const pathname = new URL(req.url()).pathname
      if (pathname.startsWith('/api/v1')) counts.set(p, count(p) + (pathname === '/api/v1/auth/refresh' ? 1 : 0))
      const bearer = req.headers().authorization?.replace(/^Bearer /, '')
      if (bearer) secrets.push(bearer)
    })
    p.on('response', async r => {
      if (new URL(r.url()).pathname.startsWith('/api/v1')) {
        const id = (await r.allHeaders())['x-request-id']
        if (id) evidence.correlation.push({ id, status: r.status() })
        const setCookie=await r.headerValue('set-cookie')
        for(const match of (setCookie??'').matchAll(/elabtrack_v2_refresh=([0-9a-f]{64})/g)) {
          secrets.push(match[1],createHash('sha256').update(match[1]).digest('hex'))
        }
      }
    })
  }
  observe(page)
  const response = await page.goto(base)
  await page.getByRole('heading', { name:'eLabTrack V2',exact:true }).waitFor()
  const headers = await response.allHeaders()
  assert.match(headers['content-security-policy'], /script-src 'self'/)
  assert.equal(headers['x-frame-options'],'DENY')
  assert.equal(headers['x-content-type-options'],'nosniff')
  assert.equal(headers['referrer-policy'],'no-referrer')
  assert.match(headers['permissions-policy'],/camera=\(\)/)
  assert.equal(headers['strict-transport-security'],undefined)
  for (const viewport of [{width:390,height:844},{width:768,height:1024},{width:1440,height:900}]) {
    await page.setViewportSize(viewport)
    assert(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth + 1))
  }
  await page.getByRole('button',{name:'Switch to dark theme'}).click()
  assert(await page.evaluate(()=>document.documentElement.classList.contains('dark')))
  await page.getByRole('button',{name:'Switch to light theme'}).click()
  await page.getByRole('link',{name:'Check service connection'}).click()
  await page.getByRole('button',{name:'Check connection'}).click()
  await page.getByText('API connected.',{exact:true}).waitFor()
  await page.reload()
  await page.getByRole('heading',{name:'eLabTrack V2 service connection'}).waitFor()
  assert.equal(cspErrors.length,0)
  check('production SPA headers, routes, API proxy, theme and three responsive viewports')

  const fixtureResponse = await fetch(base + '/api/v1/auth/register', { method:'POST',headers:{Origin:base,'Content-Type':'application/json'},body:JSON.stringify({email,password,name:'Synthetic browser integration account'}) })
  assert.equal(fixtureResponse.status,201)
  const fixture = (await fixtureResponse.json()).data
  accountID = fixture.user.id
  assert.match(accountID,/^[0-9a-f-]{36}$/)
  secrets.push(sql(`SELECT password FROM users WHERE id='${accountID}'`))
  secrets.push(fixture.access_token,fixtureResponse.headers.get('set-cookie').split(';')[0].split('=')[1])
  // The browser logs in through the retained actual central client.
  await page.goto(driver)
  await page.waitForFunction(()=>window.phase1g && window.phase1g.state().status==='unauthenticated')
  assert(await page.evaluate(()=>window.phase1g.initialMemoryEmpty))
  await page.evaluate(({email,password})=>window.phase1g.login(email,password),{email,password})
  assert.equal((await page.evaluate(()=>window.phase1g.state())).status,'authenticated')
  assert.equal((await page.evaluate(()=>window.phase1g.me())).email,email)
  let cookies = (await context.cookies()).filter(c=>c.name===cookieName)
  assert.equal(cookies.length,1)
  assert(cookies[0].httpOnly && !cookies[0].secure && cookies[0].sameSite==='Lax' && cookies[0].path==='/api/v1/auth')
  secrets.push(cookies[0].value)
  const storage = await page.evaluate(()=>({ local:Object.entries(localStorage),session:Object.entries(sessionStorage),cookie:document.cookie }))
  assert(!storage.cookie.includes(cookieName))
  for (const [key,value] of [...storage.local,...storage.session]) { assert(!/access.?token|refresh.?token|auth.?token/i.test(key));assert(!secrets.some(s=>s && value.includes(s))) }
  check('browser login, safe current account, HttpOnly/Lax/path cookie, memory-only tokens')

  let before = count(page)
  const oldCookie = cookies[0].value
  await page.reload()
  await page.waitForFunction(()=>window.phase1g?.state().status==='authenticated')
  assert.equal(count(page)-before,1)
  assert(await page.evaluate(()=>window.phase1g.initialMemoryEmpty))
  cookies=(await context.cookies()).filter(c=>c.name===cookieName)
  assert(cookies[0].value!==oldCookie,'cookie rotates');secrets.push(cookies[0].value)
  check('reload starts with empty memory and exactly one real bootstrap refresh')

  before=count(page)
  await page.evaluate(()=>window.phase1g.invalidateAccess())
  const results=await page.evaluate(()=>Promise.allSettled(Array.from({length:5},()=>window.phase1g.me())).then(r=>r.map(x=>x.status)))
  assert.deepEqual(results,Array(5).fill('fulfilled'))
  assert.equal(count(page)-before,1)
  check('five concurrent real protected 401s share one refresh and retry successfully')

  let protectedCalls=0
  await page.route('**/api/v1/auth/me',async route=>{
    protectedCalls++
    if(protectedCalls===1) await route.fulfill({status:401,contentType:'application/json',body:JSON.stringify({success:false,error:{code:'UNAUTHORIZED',message:'Authentication required.',requestId:randomUUID()},data:null,meta:null})})
    else await route.continue()
  })
  before=count(page)
  await page.evaluate(()=>window.phase1g.me())
  assert.equal(protectedCalls,2);assert.equal(count(page)-before,1)
  await page.unroute('**/api/v1/auth/me')
  check('controlled protected 401 performs one actual refresh and one retry')
  protectedCalls=0
  await page.route('**/api/v1/auth/me',async route=>{protectedCalls++;await route.fulfill({status:401,contentType:'application/json',body:JSON.stringify({success:false,error:{code:'UNAUTHORIZED',message:'Authentication required.',requestId:randomUUID()},data:null,meta:null})})})
  before=count(page)
  assert.equal(await page.evaluate(()=>window.phase1g.me().then(()=>false,()=>true)),true)
  assert.equal(protectedCalls,2);assert.equal(count(page)-before,1)
  assert.equal((await page.evaluate(()=>window.phase1g.state())).status,'unauthenticated')
  await page.unroute('**/api/v1/auth/me')
  check('persistent protected 401 stops after one retry and clears memory')

  await page.evaluate(({email,password})=>window.phase1g.login(email,password),{email,password})
  sql(`UPDATE refresh_tokens SET revoked_at=now() WHERE user_id='${accountID}' AND revoked_at IS NULL`)
  before=count(page);await page.reload()
  await page.waitForFunction(()=>window.phase1g?.state().status==='unauthenticated')
  assert.equal(count(page)-before,1)
  assert.equal((await context.cookies()).filter(c=>c.name===cookieName).length,0)
  check('revoked-session reload fails once and clears the dead cookie')

  await page.evaluate(({email,password})=>window.phase1g.login(email,password),{email,password})
  const tabB=await context.newPage();observe(tabB)
  await tabB.goto(driver)
  await tabB.waitForFunction(()=>window.phase1g?.state().status==='authenticated')
  assert.equal((await page.evaluate(()=>window.phase1g.state())).status,'authenticated')
  // Hold requests at the browser network boundary until both tabs have sent
  // their shared cookie; no backend relaxation or fabricated refresh response.
  const pending=[]
  let release
  const barrier=new Promise(r=>{release=r})
  const cookieHeaders=[]
  for(const p of [page,tabB]) await p.route('**/api/v1/auth/refresh',async route=>{
    cookieHeaders.push((await route.request().allHeaders()).cookie)
    pending.push(route);if(pending.length===2) release()
    await barrier;await route.continue()
  })
  const concurrent=await Promise.all([page,tabB].map(p=>p.evaluate(()=>window.phase1g.refresh().then(()=> 'success',e=>`denied:${e.status}`))))
  assert(cookieHeaders[0]===cookieHeaders[1],'both tabs sent the same cookie')
  assert.equal(concurrent.filter(x=>x==='success').length,1)
  assert.equal(concurrent.filter(x=>x==='denied:401').length,1)
  for(const p of [page,tabB]) await p.unroute('**/api/v1/auth/refresh')
  evidence.crossTab={ refreshResults:concurrent, cookieSurvived:(await context.cookies()).some(c=>c.name===cookieName), states:await Promise.all([page,tabB].map(p=>p.evaluate(()=>window.phase1g.state().status))) }
  check('two tabs share the cookie but racing rotation forces one tab unauthenticated',evidence.crossTab)

  await page.evaluate(({email,password})=>window.phase1g.login(email,password),{email,password})
  await tabB.reload();await tabB.waitForFunction(()=>window.phase1g?.state().status==='authenticated')
  await page.evaluate(()=>window.phase1g.logout())
  assert.equal((await page.evaluate(()=>window.phase1g.state())).status,'unauthenticated')
  assert.equal((await tabB.evaluate(()=>window.phase1g.state())).status,'authenticated')
  assert.equal((await tabB.evaluate(()=>window.phase1g.me())).email,email)
  await tabB.evaluate(()=>window.phase1g.invalidateAccess())
  assert(await tabB.evaluate(()=>window.phase1g.me().then(()=>false,()=>true)))
  assert.equal((await tabB.evaluate(()=>window.phase1g.state())).status,'unauthenticated')
  evidence.crossTab.logoutOtherTabAccessRemainedValid=true
  evidence.crossTab.recommendation='ADD CROSS-TAB SESSION COORDINATION BEFORE PRODUCT UI'
  check('logout clears one tab; other tab remains authenticated until refresh fails')

  // Local same-site origins with different ports exercise real browser CORS.
  for(const port of [14173,14174]) {
    const server=createServer((_,res)=>{res.setHeader('Content-Type','text/html');res.end('<!doctype html><title>Local CORS fixture</title>')})
    await new Promise(resolve=>server.listen(port,'127.0.0.1',resolve));servers.push(server)
  }
  const allowed=await context.newPage();observe(allowed)
  await allowed.goto('http://localhost:14173')
  const allowedHealth=await allowed.evaluate(async base=>{
    const r=await fetch(base+'/api/v1/health',{credentials:'include',headers:{'X-Request-ID':'valid-but-ignored'}})
    return {status:r.status,id:r.headers.get('x-request-id'),data:await r.json()}
  },base)
  assert.equal(allowedHealth.status,200);assert.match(allowedHealth.id,/^[0-9a-f-]{36}$/)
  await allowed.evaluate(async ({base,email,password})=>{const r=await fetch(base+'/api/v1/auth/login',{method:'POST',credentials:'include',headers:{'Content-Type':'application/json'},body:JSON.stringify({email,password})});if(r.status!==200)throw Error('allowed login failed')},{base,email,password})
  const guardedCookie=(await context.cookies()).find(c=>c.name===cookieName);secrets.push(guardedCookie.value)
  const disallowed=await context.newPage();observe(disallowed)
  await disallowed.goto('http://localhost:14174')
  // Chromium suppresses Playwright's normal response event for CORS-blocked
  // responses. Its network extra-info event still provides the real wire status.
  const cdp=await context.newCDPSession(disallowed)
  await cdp.send('Network.enable')
  const wireURLs=new Map(), wireResponses=new Map()
  cdp.on('Network.requestWillBeSent',e=>wireURLs.set(e.requestId,e.request.url))
  cdp.on('Network.responseReceivedExtraInfo',e=>wireResponses.set(e.requestId,e))
  for(const endpoint of ['refresh','logout']) {
    const denied=await disallowed.evaluate(async ({base,endpoint})=>{try{await fetch(base+'/api/v1/auth/'+endpoint,{method:'POST',credentials:'include'});return false}catch{return true}},{base,endpoint})
    const wire=[...wireResponses].find(([id])=>wireURLs.get(id)===base+'/api/v1/auth/'+endpoint)?.[1]
    assert(denied);assert.equal(wire?.statusCode,403)
    const h=Object.fromEntries(Object.entries(wire.headers).map(([k,v])=>[k.toLowerCase(),v]))
    assert.equal(h['access-control-allow-origin'],undefined);assert.equal(h['access-control-allow-credentials'],undefined)
  }
  assert((await context.cookies()).find(c=>c.name===cookieName).value===guardedCookie.value,'untrusted origin cannot mutate cookie')
  const allowedRefresh=await allowed.evaluate(async base=>{const r=await fetch(base+'/api/v1/auth/refresh',{method:'POST',credentials:'include'});return r.status},base)
  assert.equal(allowedRefresh,200)
  const preflight=await fetch(base+'/api/v1/auth/refresh',{method:'OPTIONS',headers:{Origin:'http://localhost:14173','Access-Control-Request-Method':'POST','Access-Control-Request-Headers':'content-type'}})
  assert.equal(preflight.headers.get('access-control-allow-origin'),'http://localhost:14173')
  assert.equal(preflight.headers.get('access-control-allow-credentials'),'true')
  check('real browser allowed-origin credentials/preflight and disallowed refresh/logout CSRF denial')
  assert.equal(cspErrors.length,0)
  evidence.chromiumVersion=browser.version()
  await writeFile('/tmp/elabtrack-phase1g-browser.json',JSON.stringify(evidence,null,2))
} finally {
  for(const server of servers) await new Promise(resolve=>server.close(resolve))
  if(accountID) sql(`DELETE FROM users WHERE id='${accountID}'`)
  await writeFile('/tmp/elabtrack-phase1g-secrets.json',JSON.stringify(secrets),{mode:0o600})
  await browser.close()
}
