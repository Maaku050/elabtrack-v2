// Dependency-free CDP browser QA. Uses real Vite pages and an anonymous HTTP
// refresh fixture only; no authenticated session or business API is mocked.
import fs from 'node:fs/promises'
import path from 'node:path'
import { browser } from './browser-cdp.mjs'
const args = process.argv.slice(2)
const argument = (key, fallback) => args.includes(key) ? args[args.indexOf(key) + 1] : fallback
const base = argument('--url', 'http://127.0.0.1:5173')
const production = argument('--production-url', '')
let currentOrigin = new URL(base).origin
const output = path.resolve(argument('--output', '../docs/ux/verification/phase3b'))
await fs.mkdir(output, { recursive: true })
const b = await browser(), requests = [], errors = [], warnings = [], expectedAuthResponses = [], cases = [], checks = []
b.on(async message => {
  if (message.method === 'Network.requestWillBeSent') requests.push(message.params.request.url)
  if (message.method === 'Runtime.exceptionThrown') errors.push(message.params.exceptionDetails)
  if (message.method === 'Runtime.consoleAPICalled' && message.params.type === 'error') errors.push(message.params.args)
  if (message.method === 'Runtime.consoleAPICalled' && message.params.type === 'warning') warnings.push(message.params.args)
  if (message.method === 'Log.entryAdded' && message.params.entry.level === 'warning') warnings.push(message.params.entry)
  if (message.method === 'Log.entryAdded' && message.params.entry.level === 'error') {
    const entry = message.params.entry
    if (entry.source === 'network' && entry.url?.endsWith('/api/v1/auth/refresh') && entry.text.includes('401')) expectedAuthResponses.push(entry)
    else errors.push(entry)
  }
  if (message.method === 'Fetch.requestPaused') {
    const { requestId, request } = message.params
    if (!request.url.endsWith('/api/v1/auth/refresh')) {
      errors.push({ unexpectedApi: request.url })
      await b.command('Fetch.failRequest', { requestId, errorReason: 'BlockedByClient' })
      return
    }
    await b.command('Fetch.fulfillRequest', {
      requestId, responseCode: request.method === 'OPTIONS' ? 204 : 401,
      responseHeaders: [
        { name: 'Access-Control-Allow-Origin', value: currentOrigin },
        { name: 'Access-Control-Allow-Credentials', value: 'true' },
        { name: 'Access-Control-Allow-Headers', value: 'Content-Type,X-CSRF-Token' },
        { name: 'Access-Control-Allow-Methods', value: 'POST,GET,OPTIONS' },
        { name: 'Content-Type', value: 'application/json' },
      ],
      body: request.method === 'OPTIONS' ? '' : Buffer.from(JSON.stringify({ success: false, message: 'Anonymous QA session', data: null, meta: null, error: { code: 'UNAUTHORIZED' } })).toString('base64'),
    })
  }
})
async function until(expression) {
  for (let attempt = 0; attempt < 80; attempt++) {
    if (await b.evaluate(expression)) return
    await new Promise(resolve => setTimeout(resolve, 100))
  }
  throw new Error(`Timed out: ${expression}`)
}
async function settle() {
  await b.evaluate('document.fonts.ready.then(()=>new Promise(resolve=>requestAnimationFrame(()=>requestAnimationFrame(resolve))))')
  await new Promise(resolve => setTimeout(resolve, 150))
}
async function click(selector) {
  const rect = await b.evaluate(`(()=>{const e=document.querySelector(${JSON.stringify(selector)});if(!e)throw Error('Missing control');e.scrollIntoView({block:'center'});const r=e.getBoundingClientRect();return {x:r.x+r.width/2,y:r.y+r.height/2};})()`)
  await b.command('Input.dispatchMouseEvent', { type: 'mousePressed', ...rect, button: 'left', clickCount: 1 })
  await b.command('Input.dispatchMouseEvent', { type: 'mouseReleased', ...rect, button: 'left', clickCount: 1 })
  await settle()
}
async function shot(filename, full = false) {
  const metrics = await b.command('Page.getLayoutMetrics')
  const params = { format: 'png', captureBeyondViewport: full }
  if (full) params.clip = { x: 0, y: 0, width: metrics.cssContentSize.width, height: metrics.cssContentSize.height, scale: 1 }
  const result = await b.command('Page.captureScreenshot', params)
  await fs.writeFile(path.join(output, filename), Buffer.from(result.data, 'base64'))
}
const screens = [
  { id: 'B01', route: '/__preview/borrower/home', title: 'Good morning, Alex!', widths: [390, 320, 768, 1280] },
  { id: 'B02', route: '/__preview/borrower/equipment', title: 'Equipment Catalog', widths: [390, 320, 768, 1280] },
  { id: 'S01-dashboard', route: '/__preview/staff/dashboard', title: 'Dashboard', widths: [1024, 1440] },
  { id: 'S01-pending', route: '/__preview/staff/requests', title: 'Pending Requests', widths: [1024, 1440] },
]
try {
  await b.command('Network.enable'); await b.command('Log.enable')
  await b.command('Fetch.enable', { patterns: [{ urlPattern: '*/api/v1/*', requestStage: 'Request' }] })
  for (const screen of screens) for (const width of (args.includes('--quick') ? [screen.widths.at(-1) > 1280 ? 1440 : 390] : screen.widths)) for (const theme of (args.includes('--quick') ? ['light'] : ['light', 'dark'])) {
    await b.command('Emulation.setDeviceMetricsOverride', { width, height: 900, deviceScaleFactor: 1, mobile: width < 768 })
    await b.command('Emulation.setTouchEmulationEnabled', { enabled: width < 768 })
    await b.command('Page.navigate', { url: base + screen.route })
    await until(`document.querySelector('h1')?.textContent === ${JSON.stringify(screen.title)}`)
    await settle()
    if (theme === 'dark') await click('button[aria-label="Switch to dark theme"]')
    await b.evaluate('window.scrollTo(0,0)'); await settle()
    const result = await b.evaluate(`(()=>{
      const visible=e=>{const r=e.getBoundingClientRect();const s=getComputedStyle(e);return r.width>0&&r.height>0&&s.visibility!=='hidden'&&s.display!=='none'};
      const clipCandidates=[...document.querySelectorAll('.app-brand-copy,.page-heading h1,.metric-copy p,.equipment-copy h2,.panel-heading h2,.activity-copy strong')].filter(visible);
      const clipped=clipCandidates.filter(e=>e.scrollWidth>e.clientWidth+1).map(e=>e.textContent);
      const smallTargets=[...document.querySelectorAll('button,a,input')].filter(visible).filter(e=>!e.closest('.selection-cell')&&!e.classList.contains('skip-link')).map(e=>({name:e.getAttribute('aria-label')||e.textContent,width:e.getBoundingClientRect().width,height:e.getBoundingClientRect().height})).filter(r=>r.width<43.9||r.height<43.9);
      const scroll=document.querySelector('.table-scroll');
      const boxes=Object.fromEntries(['.borrower-top-bar','.borrower-bottom-nav','.staff-sidebar','.staff-top-bar','.page-heading','.metric-card','.equipment-card','.pending-table'].map(selector=>{const e=document.querySelector(selector),r=e?.getBoundingClientRect();return [selector,r?{x:r.x,y:r.y,width:r.width,height:r.height}:null]}));
      return {boxes,title:document.querySelector('h1').textContent,viewport:innerWidth,clientWidth:document.documentElement.clientWidth,pageWidth:document.documentElement.scrollWidth,clipped,smallTargets,theme:document.documentElement.classList.contains('dark')?'dark':'light',table:scroll?{width:scroll.clientWidth,contentWidth:scroll.scrollWidth,tabIndex:scroll.tabIndex,label:scroll.getAttribute('aria-label')}:null};
    })()`)
    await shot(`${screen.id}-${width}-${theme}.png`)
    if (screen.id.startsWith('B') && width === 390 && theme === 'light') await shot(`${screen.id}-${width}-${theme}-full.png`, true)
    await b.evaluate('window.scrollTo(0,document.documentElement.scrollHeight)'); await settle()
    result.footer = await b.evaluate(`(()=>{const n=document.querySelector('.borrower-bottom-nav'),p=document.querySelector('.preview-notice');return n?{navigationTop:n.getBoundingClientRect().top,contentBottom:p.getBoundingClientRect().bottom}:null;})()`)
    result.dock = await b.evaluate(`(()=>{const d=document.querySelector('.cart-dock'),n=document.querySelector('.borrower-bottom-nav');return d?{bottom:d.getBoundingClientRect().bottom,navigationTop:n.getBoundingClientRect().top}:null;})()`)
    result.id = screen.id; result.width = width; result.file = `${screen.id}-${width}-${theme}.png`
    result.pass = result.pageWidth <= result.clientWidth && result.clipped.length === 0 && result.smallTargets.length === 0 && result.theme === theme && (!result.footer || result.footer.contentBottom <= result.footer.navigationTop) && (!result.dock || result.dock.bottom <= result.dock.navigationTop - 6)
    cases.push(result); console.log(`${result.pass ? 'PASS' : 'FAIL'} ${screen.id} ${width} ${theme}`)
  }
  // Real pointer + keyboard checks against retained Base UI overlays.
  await b.command('Emulation.setDeviceMetricsOverride', { width: 390, height: 900, deviceScaleFactor: 1, mobile: false })
  await b.command('Page.navigate', { url: base + '/__preview/borrower/equipment' })
  await until(`document.querySelector('h1')?.textContent==='Equipment Catalog'`); await settle()
  await click('[aria-label="Equipment filters"]')
  await until(`!!document.querySelector('[role="dialog"]')`)
  const focusInFilter = await b.evaluate(`document.querySelector('[role="dialog"]').contains(document.activeElement)`)
  await b.command('Input.dispatchKeyEvent', { type: 'keyDown', key: 'Escape', code: 'Escape', windowsVirtualKeyCode: 27 })
  await b.command('Input.dispatchKeyEvent', { type: 'keyUp', key: 'Escape', code: 'Escape', windowsVirtualKeyCode: 27 })
  await settle()
  checks.push({ name: 'Filter Sheet Escape closes and restores focus', pass: focusInFilter && await b.evaluate(`!document.querySelector('[role="dialog"]') && document.activeElement.getAttribute('aria-label')==='Equipment filters'`) })
  await click('.cart-dock button')
  await until(`!!document.querySelector('[role="dialog"]')`)
  let trapped = true
  for (let i = 0; i < 8; i++) {
    await b.command('Input.dispatchKeyEvent', { type: 'keyDown', key: 'Tab', code: 'Tab', windowsVirtualKeyCode: 9 })
    await b.command('Input.dispatchKeyEvent', { type: 'keyUp', key: 'Tab', code: 'Tab', windowsVirtualKeyCode: 9 })
    await settle()
    trapped = trapped && await b.evaluate(`document.querySelector('[role="dialog"]').contains(document.activeElement)`)
  }
  await shot('B02-selection-dialog.png')
  checks.push({ name: 'Selection dialog traps keyboard focus; contains no business action', pass: trapped && await b.evaluate(`![...document.querySelectorAll('[role="dialog"] button')].some(e=>/submit|approve|return|pay|clear fine/i.test(e.textContent))`) })
  await b.command('Input.dispatchKeyEvent', { type: 'keyDown', key: 'Escape', code: 'Escape', windowsVirtualKeyCode: 27 }); await settle()
  await b.command('Page.navigate', { url: base + '/__preview/staff/dashboard' }); await until(`document.querySelector('h1')?.textContent==='Dashboard'`); await settle()
  await click('[aria-label="Toggle navigation"]'); await until(`!!document.querySelector('[role="dialog"]')`)
  await shot('S01-navigation-sheet.png')
  checks.push({ name: 'Staff smaller-width navigation Sheet is accessible', pass: await b.evaluate(`document.querySelector('[role="dialog"]').contains(document.activeElement) && !!document.querySelector('[role="dialog"] nav[aria-label="Staff navigation"]')`) })
  await b.command('Page.navigate', { url: base + '/__preview/admin/dashboard' }); await until(`document.querySelector('h1')?.textContent==='Dashboard'`); await settle()
  checks.push({ name: 'Admin shell includes labeled Admin additions without authenticating', pass: await b.evaluate(`(async()=>{const {useAuthStore}=await import('/src/stores/auth-store.ts');return !!document.querySelector('nav[aria-label="Admin navigation"]') && useAuthStore.getState().accessToken===null && useAuthStore.getState().user===null;})()`) })
  await b.command('Emulation.setDeviceMetricsOverride', { width: 1024, height: 900, deviceScaleFactor: 1, mobile: false })
  await b.command('Page.navigate', { url: base + '/__preview/staff/requests' }); await until(`document.querySelector('h1')?.textContent==='Pending Requests'`); await settle()
  await b.evaluate(`document.querySelector('.table-scroll').focus()`)
  await b.command('Input.dispatchKeyEvent', { type: 'keyDown', key: 'ArrowRight', code: 'ArrowRight', windowsVirtualKeyCode: 39 })
  await b.command('Input.dispatchKeyEvent', { type: 'keyUp', key: 'ArrowRight', code: 'ArrowRight', windowsVirtualKeyCode: 39 }); await settle()
  checks.push({ name: 'Tablet table scroll is keyboard reachable and bounded', pass: await b.evaluate(`document.querySelector('.table-scroll').scrollLeft>0 && window.scrollX===0`) })
  await b.command('Emulation.setDeviceMetricsOverride', { width: 1440, height: 900, deviceScaleFactor: 1, mobile: false })
  await b.command('Page.navigate', { url: base + '/__preview/staff/dashboard' }); await until(`document.querySelector('h1')?.textContent==='Dashboard'`); await settle()
  await click('[aria-label="Toggle navigation"]')
  checks.push({ name: 'Desktop sidebar collapse preserves workspace and named navigation', pass: await b.evaluate(`document.querySelector('.staff-sidebar').getBoundingClientRect().width===76 && document.querySelector('nav[aria-label="Staff navigation"] a').getAttribute('aria-label')==='Dashboard'`) })
  await click('[aria-label="Toggle navigation"]')
  // A narrow viewport exercises reflow corresponding to 200% of 640px;
  // full browser zoom and virtual keyboard behavior remain separate pilot checks.
  await b.command('Emulation.setEmulatedMedia', { features: [{ name: 'prefers-reduced-motion', value: 'reduce' }] })
  checks.push({ name: 'Reduced motion removes spinning animations', pass: await b.evaluate(`(()=>{const e=document.createElement('div');e.className='animate-spin';document.body.append(e);const result=getComputedStyle(e).animationName==='none';e.remove();return result;})()`) })
  const bundles = (await fs.readdir('dist/assets')).filter(file => file.endsWith('.js'))
  const leaked = []
  for (const file of bundles) {
    const body = await fs.readFile(path.join('dist/assets', file), 'utf8')
    if (/__preview|PREVIEW-0215|BorrowerHomePreview|EquipmentCatalogPreview/.test(body)) leaked.push(file)
  }
  checks.push({ name: 'Production JS excludes development routes and synthetic fixtures', pass: leaked.length === 0 })
  if (production) {
    currentOrigin = new URL(production).origin
    await b.command('Page.navigate', { url: production + '/__preview/borrower/home' })
    await until(`document.querySelector('h1')?.textContent.includes('Page not found')`)
    checks.push({ name: 'Production preview URL resolves to existing404', pass: await b.evaluate(`!document.body.innerText.includes('Synthetic data')`) })
    await b.command('Page.navigate', { url: production + '/' })
    await until(`document.querySelector('h1')?.textContent==='eLabTrack V2'`)
    checks.push({ name: 'Production foundation root remains available', pass: true })
  }
  const report = { productionUrl: production || null, productionBundles: bundles, generatedAt: new Date().toISOString(), browser: (await b.command('Browser.getVersion')).product, url: base, anonymousRefreshFixture: true, cases, checks, applicationErrors: errors, applicationWarnings: warnings, expectedAnonymous401Count: expectedAuthResponses.length, businessRequests: requests.filter(url => url.includes('/api/v1/') && !url.endsWith('/auth/refresh')), requests: [...new Set(requests)].filter(url => url.includes('/api/v1/')), pass: cases.every(c => c.pass) && checks.every(c => c.pass) && errors.length === 0 && warnings.length === 0 }
  await fs.writeFile(path.join(output, 'RESULTS.json'), JSON.stringify(report, null, 2) + '\n')
  console.log(JSON.stringify({ cases: cases.length, checks, applicationErrors: errors, applicationWarnings: warnings, pass: report.pass }, null, 2))
  if (!report.pass) process.exitCode = 1
} finally { await b.close() }
