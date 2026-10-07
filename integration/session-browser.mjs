import assert from 'node:assert/strict'
import { createHash } from 'node:crypto'

const digest = value => createHash('sha256').update(value).digest('hex')
const channelName = 'elabtrack_v2.session.v1'

export async function instrumentSessionMessages(context) {
  await context.addInitScript(name => {
    const Native = window.BroadcastChannel
    window.phase1iEvents = []
    window.phase1iReceived = 0
    window.BroadcastChannel = class extends Native {
      constructor(channel) { super(channel); if (channel === name) this.addEventListener('message', () => { window.phase1iReceived++ }) }
      postMessage(message) { if (this.name === name) window.phase1iEvents.push(structuredClone(message)); return super.postMessage(message) }
    }
  }, channelName)
}

/** Real pages import the actual singleton API/store/coordinator. Only network
 * timing, synthetic runtime-role DML and ignored driver cache probes are test
 * seams; no authentication response/credential/bootstrap endpoint is invented. */
export async function verifySessionTabs({ browser, context, page, observe, secrets, check, evidence, sql, accountID, email, password, base, driver, cookieName }) {
  const state = p => p.evaluate(() => window.phase1g.state())
  const authenticated = p => p.waitForFunction(() => window.phase1g?.state().status === 'authenticated')
  const unauthenticated = p => p.waitForFunction(() => window.phase1g?.state().status === 'unauthenticated')
  const login = () => page.evaluate(({email,password}) => window.phase1g.login(email,password), {email,password})
  const newTab = async (target = context) => { const p = await target.newPage(); observe(p); await p.goto(driver); await authenticated(p); return p }
  const cookie = async target => (await target.cookies()).find(c => c.name === cookieName)
  const privateCaches = tabs => Promise.all(tabs.map(p => p.evaluate(() => window.phase1g.cachePrivate())))
  const assertCachesCleared = async tabs => {
    for (const p of tabs) assert.deepEqual(await p.evaluate(() => window.phase1g.cacheState()), { privatePresent: false, publicPresent: true })
  }
  await login()
  let tabs = [page, await newTab()]
  assert.deepEqual(await page.evaluate(() => window.phase1g.coordination()), { exclusive: true, broadcast: true })

  async function pressure(tabs) {
    const presented = [], results = []
    let active = 0, max = 0, release, started
    const gate = new Promise(resolve => { release = resolve })
    const first = new Promise(resolve => { started = resolve })
    const responded = response => {
      if (new URL(response.url()).pathname === '/api/v1/auth/refresh') { active--; results.push(response.status()) }
    }
    const intercept = async route => {
      const headers = await route.request().allHeaders()
      const value = headers.cookie?.match(/elabtrack_v2_refresh=([0-9a-f]{64})/)?.[1]
      assert(value, 'refresh cookie present (value withheld)')
      secrets.push(value, digest(value)); presented.push(digest(value))
      max = Math.max(max, ++active)
      if (presented.length === 1) { started(); await gate }
      await route.continue()
    }
    for (const p of tabs) { p.on('response', responded); await p.route('**/api/v1/auth/refresh', intercept) }
    const tasks = tabs.map(p => p.evaluate(() => window.phase1g.refresh().then(() => 'success', e => `denied:${e.status}`)))
    try {
      await first
      // Observe actual queued browser locks while the first HTTP request is
      // held, rather than assume other evaluate calls started simultaneously.
      await tabs[0].waitForFunction(async count => (await navigator.locks.query()).pending.filter(l => l.name === 'elabtrack_v2.session-cookie.v1').length === count, tabs.length - 1)
      assert.equal(presented.length, 1)
      release(); assert.deepEqual(await Promise.all(tasks), Array(tabs.length).fill('success'))
      assert.equal(max, 1); assert.equal(new Set(presented).size, tabs.length)
      assert.deepEqual(results, Array(tabs.length).fill(200))
      for (const p of tabs) assert.equal((await state(p)).status, 'authenticated')
      const rows = sql(`SELECT count(*) FROM refresh_tokens WHERE token_hash IN (${presented.map(h => `'${h}'`).join(',')}) AND revoked_at IS NOT NULL AND replaced_by IS NOT NULL`)
      assert.equal(Number(rows), tabs.length)
      const final = await cookie(context); assert(final)
      assert.equal(sql(`SELECT count(*) FROM refresh_tokens WHERE token_hash='${digest(final.value)}' AND revoked_at IS NULL`), '1')
      return { tabs: tabs.length, refreshStatuses: results, maxNetworkRefreshesInFlight: max, distinctPresentedCredentials: new Set(presented).size, consumedAndLinked: Number(rows), cookieSurvived: true, states: await Promise.all(tabs.map(async p => (await state(p)).status)) }
    } finally {
      release()
      for (const p of tabs) { await p.unroute('**/api/v1/auth/refresh', intercept); p.off('response', responded) }
    }
  }
  evidence.crossTab = await pressure(tabs)
  check('two real tabs serialize cookie rotations and each retains its own memory access', evidence.crossTab)
  tabs.push(await newTab())
  evidence.threeTabs = await pressure(tabs)
  check('three real tabs serialize multiple waiting refreshes', evidence.threeTabs)

  let requests = 0
  const count = req => { if (new URL(req.url()).pathname === '/api/v1/auth/refresh') requests++ }
  for (const p of tabs) p.on('request', count)
  await Promise.all(tabs.map(p => p.reload())); await Promise.all(tabs.map(authenticated))
  assert.equal(requests, 3)
  for (const p of tabs) assert(await p.evaluate(() => window.phase1g.initialMemoryEmpty))
  evidence.simultaneousReload = { tabs: 3, refreshRequests: requests, states: await Promise.all(tabs.map(async p => (await state(p)).status)) }
  check('simultaneous three-tab reload restores once per document without a refresh storm', evidence.simultaneousReload)

  // Deliver a real channel event from an older attempt after all three newer
  // bootstrap successes. Only allowlisted metadata is synthesized here.
  const old = (await tabs[1].evaluate(() => window.phase1iEvents)).find(m => m.type === 'refresh-started')
  // The previous page was reloaded, so use its attempt before another new turn.
  await page.evaluate(() => window.phase1g.refresh())
  const received = await Promise.all(tabs.map(p => p.evaluate(() => window.phase1iReceived)))
  await page.evaluate(({name,old}) => { const c = new BroadcastChannel(name); c.postMessage({ ...old, type: 'session-invalidated' }); c.close() }, {name: channelName,old})
  await Promise.all(tabs.map((p,i) => p.waitForFunction(previous => window.phase1iReceived > previous, received[i])))
  for (const p of tabs) assert.equal((await state(p)).status, 'authenticated')
  check('real late stale invalidation cannot erase newer authenticated state')

  await privateCaches(tabs)
  const beforeLogout = requests
  await page.evaluate(() => window.phase1g.logout()); await Promise.all(tabs.map(unauthenticated))
  await assertCachesCleared(tabs); assert.equal(requests, beforeLogout); assert.equal(await cookie(context), undefined)
  for (const p of tabs) assert.equal(await p.evaluate(() => window.phase1g.me().then(() => false, () => true)), true)
  assert.equal(requests, beforeLogout)
  evidence.logoutPropagation = { tabs: 3, promptlyUnauthenticated: true, privateCachesRemoved: true, publicCachesRetained: true, peerRefreshRequests: 0, cookieCleared: true, subsequentProtectedRequestsDenied: true }
  check('logout broadcasts promptly, removes private caches and leaves public cache', evidence.logoutPropagation)

  await login(); await Promise.all(tabs.slice(1).map(p => p.reload())); await Promise.all(tabs.map(authenticated))
  await privateCaches(tabs)
  let replyStarted, releaseReply
  const reply = new Promise(resolve => { replyStarted = resolve })
  const release = new Promise(resolve => { releaseReply = resolve })
  let logoutRequests = 0
  const observeLogout = req => { if (new URL(req.url()).pathname === '/api/v1/auth/logout') logoutRequests++ }
  tabs[1].on('request', observeLogout)
  await page.route('**/api/v1/auth/refresh', async route => {
    const actual = await route.fetch(); assert.equal(actual.status(), 200)
    replyStarted(); await release; await route.fulfill({ response: actual })
  })
  const delayed = page.evaluate(() => window.phase1g.refresh().then(() => 'resurrected', e => e.code))
  await reply
  const peerLogout = tabs[1].evaluate(() => window.phase1g.logout())
  await Promise.all(tabs.map(unauthenticated)); await assertCachesCleared(tabs)
  assert.equal(logoutRequests, 0, 'logout must wait for the pending replacement response')
  releaseReply(); assert.equal(await delayed, 'SESSION_CHANGED'); await peerLogout
  assert.equal(logoutRequests, 1); assert.equal(await cookie(context), undefined)
  await page.unroute('**/api/v1/auth/refresh'); tabs[1].off('request', observeLogout)
  evidence.logoutDuringRefresh = { committedRefreshResponseHeld: true, peerMemoryClearedBeforeReply: true, logoutWaitedForReplacement: true, oldReplyCannotRestoreMemory: true, replacementCookieRevokedAndCleared: true }
  check('peer logout fences a delayed successful refresh and revokes its replacement afterward', evidence.logoutDuringRefresh)

  await login(); await Promise.all(tabs.slice(1).map(p => p.reload())); await Promise.all(tabs.map(authenticated))
  const logoutStatuses = []
  const logoutResponse = r => { if (new URL(r.url()).pathname === '/api/v1/auth/logout') logoutStatuses.push(r.status()) }
  for (const p of tabs) p.on('response', logoutResponse)
  await Promise.all(tabs.map(p => p.evaluate(() => window.phase1g.logout())))
  await Promise.all(tabs.map(unauthenticated)); assert.equal(await cookie(context), undefined)
  assert.deepEqual(logoutStatuses, [204, 204, 204])
  for (const p of tabs) p.off('response', logoutResponse)
  await page.reload(); await unauthenticated(page)
  evidence.concurrentLogout = { tabs: 3, statuses: logoutStatuses, peerMessagesCannotCancelRevocation: true, cookieCleared: true, reloadRemainsUnauthenticated: true }
  check('three simultaneous logout intents all reach the server and reload stays logged out', evidence.concurrentLogout)

  await login()
  for (const p of tabs.slice(1)) assert.equal((await state(p)).status, 'unauthenticated')
  await Promise.all(tabs.slice(1).map(p => p.reload())); await Promise.all(tabs.map(authenticated))
  await privateCaches(tabs)
  sql(`UPDATE users SET is_active=false, updated_at=now() WHERE id='${accountID}'`)
  assert.equal(await tabs[1].evaluate(() => window.phase1g.me().then(() => 200, e => e.status)), 403)
  await Promise.all(tabs.map(unauthenticated)); await assertCachesCleared(tabs)
  evidence.disablePropagation = { authoritativeCurrentAccountStatus: 403, allTabsUnauthenticated: true, privateCachesRemoved: true, cookieMutation: false }
  check('runtime-role account disable is authoritative and invalidates all three tabs', evidence.disablePropagation)
  sql(`UPDATE users SET is_active=true, updated_at=now() WHERE id='${accountID}'`)
  await login(); await Promise.all(tabs.slice(1).map(p => p.reload())); await Promise.all(tabs.map(authenticated))
  await privateCaches(tabs)
  sql(`UPDATE refresh_tokens SET revoked_at=now(), updated_at=now() WHERE user_id='${accountID}' AND revoked_at IS NULL`)
  assert.equal(await page.evaluate(() => window.phase1g.refresh().then(() => 200, e => e.status)), 401)
  await Promise.all(tabs.map(unauthenticated)); await assertCachesCleared(tabs)
  assert(await cookie(context), 'invalid cookie is left untouched until explicit logout')
  evidence.revocationPropagation = { refreshStatus: 401, allTabsUnauthenticated: true, privateCachesRemoved: true, rejectedRefreshLeavesCookieUntouched: true }
  check('exclusive refresh denial propagates generic session invalidation', evidence.revocationPropagation)

  await login(); await Promise.all(tabs.slice(1).map(p => p.reload())); await Promise.all(tabs.map(authenticated))
  // A real owning page closes before its intercepted request reaches the API.
  // The browser releases ownership without relying on any completion event.
  let ownerStarted, releaseClosed
  const owned = new Promise(resolve => { ownerStarted = resolve })
  const closed = new Promise(resolve => { releaseClosed = resolve })
  const owner = tabs[2]
  await owner.route('**/api/v1/auth/refresh', async route => { ownerStarted(); await closed; try { await route.continue() } catch { /* closed target */ } })
  const disappearing = owner.evaluate(() => window.phase1g.refresh()).catch(() => undefined)
  await owned
  const recovering = page.evaluate(() => window.phase1g.refresh())
  await page.waitForFunction(async () => (await navigator.locks.query()).pending.some(l => l.name === 'elabtrack_v2.session-cookie.v1'))
  await owner.close(); releaseClosed(); await disappearing; await recovering
  assert.equal((await state(page)).status, 'authenticated'); assert(await cookie(context))
  tabs = tabs.slice(0, 2)
  evidence.ownerClose = { realOwnerPageClosed: true, completionBroadcastRequired: false, waitingTabRecovered: true, cookieSurvived: true }
  check('closing the real owner releases its Web Lock and the waiting tab recovers', evidence.ownerClose)

  // Bypass the coordinator deliberately to preserve the server race regression.
  // Both actual requests carry an identical credential. Hold the actual losing
  // response until after the successful response is installed in the browser.
  let releaseRequests, releaseLoser, winnerFinished
  const ready = new Promise(resolve => { releaseRequests = resolve })
  const late = new Promise(resolve => { releaseLoser = resolve })
  const winner = new Promise(resolve => { winnerFinished = resolve })
  const rawCookies = [], rawStatuses = [], responseCookies = []
  for (const p of tabs) await p.route('**/api/v1/auth/refresh', async route => {
    const raw = (await route.request().allHeaders()).cookie
    rawCookies.push(raw); if (rawCookies.length === 2) releaseRequests()
    await ready
    const actual = await route.fetch()
    rawStatuses.push(actual.status()); responseCookies.push({ status: actual.status(), mutatesCookie: !!actual.headers()['set-cookie'] })
    if (actual.status() === 401) await late
    await route.fulfill({ response: actual })
  })
  const rawTasks = tabs.map(p => p.evaluate(() => fetch('/api/v1/auth/refresh', { method: 'POST', credentials: 'include' }).then(r => r.status)).then(status => { if (status === 200) winnerFinished(); return status }))
  await winner
  const winnerCookie = await cookie(context); assert(winnerCookie)
  releaseLoser(); const outcomes = await Promise.all(rawTasks)
  assert.deepEqual([...outcomes].sort(), [200, 401]); assert.equal(rawCookies[0], rawCookies[1])
  assert.deepEqual(responseCookies.find(r => r.status === 401), { status: 401, mutatesCookie: false })
  assert.equal((await cookie(context)).value, winnerCookie.value)
  for (const p of tabs) await p.unroute('**/api/v1/auth/refresh')
  await page.evaluate(() => window.phase1g.refresh())
  evidence.uncoordinatedRace = { sameCredentialConfirmed: true, statuses: rawStatuses.sort(), losingResponseDeliveredLast: true, losingResponseSetCookie: false, winnerCookieSurvived: true, successorRefreshSucceeded: true }
  check('real late replay 401 cannot delete the winner cookie even when coordination is bypassed', evidence.uncoordinatedRace)

  for (const p of tabs) {
    const storage = await p.evaluate(async () => ({ local: Object.entries(localStorage), session: Object.entries(sessionStorage), databases: (await indexedDB.databases()).map(d => d.name), cookie: document.cookie, messages: window.phase1iEvents }))
    assert(!storage.cookie.includes(cookieName)); assert.deepEqual(storage.databases, [])
    for (const [key,value] of [...storage.local,...storage.session]) { assert(!/access.?token|refresh.?token|auth.?token/i.test(key)); assert(!secrets.some(s => s && value.includes(s)), 'storage credential sentinel (withheld)') }
    for (const message of storage.messages) {
      const keys = ['version','type','tabId','attemptId','epoch','timestamp', ...(message.type === 'refresh-failed' ? ['failure'] : [])]
      assert(Object.keys(message).every(k => keys.includes(k)), 'unexpected coordination field (withheld)')
      assert(!secrets.some(s => s && JSON.stringify(message).includes(s)), 'coordination credential sentinel (withheld)')
    }
  }
  evidence.storageAndMessages = { localAndSessionAuthCredentials: 0, indexedDBDatabases: 0, readableRefreshCookies: 0, credentialBroadcasts: 0 }
  check('real Chromium storage/IndexedDB and observed channel messages contain no usable credentials', evidence.storageAndMessages)

  const fallback = await browser.newContext()
  try {
    await fallback.addInitScript(() => { Object.defineProperty(window, 'BroadcastChannel', { value: undefined }); Object.defineProperty(navigator, 'locks', { value: undefined }) })
    await fallback.addCookies(await context.cookies())
    const peers = await Promise.all([fallback.newPage(), fallback.newPage()])
    let release
    const barrier = new Promise(resolve => { release = resolve }), headers = []
    for (const p of peers) { observe(p); await p.route('**/api/v1/auth/refresh', async route => { headers.push((await route.request().allHeaders()).cookie); if (headers.length === 2) release(); await barrier; await route.continue() }) }
    await Promise.all(peers.map(p => p.goto(driver)))
    await Promise.all(peers.map(p => p.waitForFunction(() => ['authenticated','unauthenticated'].includes(window.phase1g?.state().status))))
    assert.equal(headers[0], headers[1]); assert(await cookie(fallback))
    const states = await Promise.all(peers.map(async p => (await state(p)).status))
    assert.deepEqual([...states].sort(), ['authenticated','unauthenticated'])
    for (const p of peers) await p.unroute('**/api/v1/auth/refresh')
    const loser = peers[states.indexOf('unauthenticated')]
    await loser.reload(); await authenticated(loser)
    for (const p of peers) { assert.equal((await state(p)).status, 'authenticated'); assert.deepEqual(await p.evaluate(() => window.phase1g.coordination()), { exclusive: false, broadcast: false }) }
    evidence.fallback = { primitivesUnavailable: ['BroadcastChannel','Web Locks'], initialStates: states, cookieSurvivedRace: true, deliberateReloadRecovered: true, crossTabLogoutOptimization: false }
    check('actual Chromium fallback safely preserves the winning cookie and supports deliberate recovery', evidence.fallback)
    const locksOnly = await browser.newContext()
    try {
      await locksOnly.addInitScript(() => { Object.defineProperty(window, 'BroadcastChannel', { value: undefined }) })
      await locksOnly.addCookies(await fallback.cookies())
      const pages = await Promise.all([locksOnly.newPage(), locksOnly.newPage()]), presented = [], statuses = []
      for (const p of pages) {
        observe(p)
        p.on('request', req => { if (new URL(req.url()).pathname === '/api/v1/auth/refresh') presented.push(req) })
        p.on('response', r => { if (new URL(r.url()).pathname === '/api/v1/auth/refresh') statuses.push(r.status()) })
      }
      await Promise.all(pages.map(p => p.goto(driver))); await Promise.all(pages.map(authenticated))
      assert.deepEqual(statuses, [200, 200]); assert.equal(presented.length, 2)
      const values = await Promise.all(presented.map(async req => (await req.allHeaders()).cookie))
      assert.notEqual(values[0], values[1])
      for (const p of pages) assert.deepEqual(await p.evaluate(() => window.phase1g.coordination()), { exclusive: true, broadcast: false })
      evidence.locksOnlyFallback = { broadcastUnavailable: true, webLocksAvailable: true, bootstrapStatuses: statuses, distinctPresentedCredentials: true, peerLifecycleNotificationUnavailable: true }
      check('actual Chromium without BroadcastChannel still serializes bootstrap using Web Locks', evidence.locksOnlyFallback)
    } finally { await locksOnly.close() }
  } finally { await fallback.close() }
  for (const p of tabs) p.off('request', count)
  await page.evaluate(() => window.phase1g.logout())
}
