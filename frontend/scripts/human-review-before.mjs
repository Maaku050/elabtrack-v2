// Explicit fictional PostgreSQL harness; credentials are never included in evidence.
import fs from 'node:fs/promises'
import path from 'node:path'
import { browser } from './browser-cdp.mjs'
const base='http://127.0.0.1:15177',fixture=JSON.parse(await fs.readFile('../backend/tmp/presentation-demo/browser-fixtures.json','utf8')),output=path.resolve('../docs/ux/verification/human-review/before'),screens=[],errors=[]
const b=await browser(),wait=ms=>new Promise(r=>setTimeout(r,ms));b.on(m=>{if(m.method==='Runtime.exceptionThrown')errors.push(m.params.exceptionDetails.text)})
async function until(expr){for(let i=0;i<180;i++){try{if(await b.evaluate(expr))return}catch(e){if(!/context|Inspected target|collected/.test(e.message))throw e}await wait(100)}throw Error('Timed out: '+expr)}
async function navigate(route){await b.command('Page.navigate',{url:base+route});await until("document.readyState==='complete'");await wait(200)}
async function click(selector){await until(`!!document.querySelector(${JSON.stringify(selector)})`);const rect=await b.evaluate(`(()=>{const e=document.querySelector(${JSON.stringify(selector)});e.scrollIntoView({block:'center'});const r=e.getBoundingClientRect();return{x:r.x+r.width/2,y:r.y+r.height/2}})()`);for(const type of ['mousePressed','mouseReleased'])await b.command('Input.dispatchMouseEvent',{type,...rect,button:'left',clickCount:1})}
async function button(text){await until(`[...document.querySelectorAll('button')].some(e=>e.textContent.trim()===${JSON.stringify(text)}&&!e.disabled)`);await b.evaluate(`(()=>{const e=[...document.querySelectorAll('button')].find(e=>e.textContent.trim()===${JSON.stringify(text)});e.dataset.phase11Button='true'})()`);await click('[data-phase11-button]');await b.evaluate("document.querySelector('[data-phase11-button]')?.removeAttribute('data-phase11-button')")}
async function fill(selector,value){await click(selector);for(const type of ['keyDown','keyUp'])await b.command('Input.dispatchKeyEvent',{type,key:'a',code:'KeyA',windowsVirtualKeyCode:65,modifiers:2});await b.command('Input.insertText',{text:value})}
async function screenshot(name,width=1440,theme='light',height=1000){await b.command('Emulation.setDeviceMetricsOverride',{width,height,deviceScaleFactor:1,mobile:width<768});if(await b.evaluate("document.documentElement.classList.contains('dark')")!==(theme==='dark'))await click('button[aria-label="Switch to '+theme+' theme"]');await wait(150);const shot=await b.command('Page.captureScreenshot',{format:'png',captureBeyondViewport:false});const file=`${name}-${width}-${height}-${theme}.png`;await fs.writeFile(path.join(output,file),Buffer.from(shot.data,'base64'));screens.push(file)}

async function login(role){await b.command('Network.clearBrowserCookies');await navigate('/login');await until("!!document.querySelector('#login-email')");await fill('#login-email',fixture[role].email);await fill('#login-password',fixture[role].password);await button('Sign In');await until("location.pathname!='/login'&&!!document.querySelector('h1')")}



await fs.mkdir(output,{recursive:true})


async function ready(){await until("!!document.querySelector('h1')&&!document.querySelector('.directory-loading')&&!document.body.textContent.includes('Loading current metrics')&&!document.body.textContent.includes('Loading…')");await wait(350)}
try {
 await login('admin');
 for (const route of ['/staff/dashboard','/staff/requests']) {await navigate(route);await ready();for(const theme of ['light','dark'])await screenshot(route.replaceAll('/','-'),1366,theme,768)}
 await navigate('/staff/requests/direct');await ready();await fill('input[aria-label="Search issuance borrower"]',fixture.faculty01.name);await click('button[aria-label="Select borrower '+fixture.faculty01.name+'"]');await until("document.querySelectorAll('.selection-equipment').length>0");for(const theme of ['light','dark'])await screenshot('direct-browse',1366,theme,768);
 await login('student01');
 for(const route of ['/borrower/home','/borrower/account','/borrower/equipment']){await navigate(route);await ready();for(const theme of ['light','dark'])await screenshot(route.replaceAll('/','-'),390,theme,844)}
 for(const width of [320,390,440]){await navigate('/borrower/equipment');await ready();for(const theme of ['light','dark']){await screenshot('catalog-closed',width,theme,844);await click('.catalog-filter-panel summary');await screenshot('catalog-open',width,theme,844);await click('.catalog-filter-panel summary')}}
 await fs.writeFile(path.join(output,'acceptance.json'),JSON.stringify({scope:'Read-only baseline before human-review changes; fictional isolated accounts only.',screens,errors},null,2));console.log(JSON.stringify({screens:screens.length,errors}));
}finally{await b.close()}
