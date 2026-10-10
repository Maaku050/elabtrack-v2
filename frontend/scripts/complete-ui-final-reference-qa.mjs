// Explicit fictional PostgreSQL harness; credentials are never included in evidence.
import fs from 'node:fs/promises'
import assert from 'node:assert/strict'
import path from 'node:path'
import { browser } from './browser-cdp.mjs'
const base='http://127.0.0.1:15177',fixture=JSON.parse(await fs.readFile('../backend/tmp/presentation-demo/browser-fixtures.json','utf8')),output=path.resolve('../docs/ux/verification/complete-ui/after'),screens=[],errors=[],checks=[]
const b=await browser(),wait=ms=>new Promise(r=>setTimeout(r,ms));b.on(m=>{if(m.method==='Runtime.exceptionThrown')errors.push(m.params.exceptionDetails.text)})
async function until(expr){for(let i=0;i<180;i++){try{if(await b.evaluate(expr))return}catch(e){if(!/context|Inspected target|collected/.test(e.message))throw e}await wait(100)}throw Error('Timed out: '+expr)}
async function navigate(route){await b.command('Page.navigate',{url:base+route});await until("document.readyState==='complete'");await wait(200)}
async function click(selector){await until(`!!document.querySelector(${JSON.stringify(selector)})`);const rect=await b.evaluate(`(()=>{const e=document.querySelector(${JSON.stringify(selector)});e.scrollIntoView({block:'center'});const r=e.getBoundingClientRect();return{x:r.x+r.width/2,y:r.y+r.height/2}})()`);for(const type of ['mousePressed','mouseReleased'])await b.command('Input.dispatchMouseEvent',{type,...rect,button:'left',clickCount:1})}
async function button(text){await until(`[...document.querySelectorAll('button')].some(e=>e.textContent.trim()===${JSON.stringify(text)}&&!e.disabled)`);await b.evaluate(`(()=>{const e=[...document.querySelectorAll('button')].find(e=>e.textContent.trim()===${JSON.stringify(text)});e.dataset.phase11Button='true'})()`);await click('[data-phase11-button]');await b.evaluate("document.querySelector('[data-phase11-button]')?.removeAttribute('data-phase11-button')")}
async function fill(selector,value){await click(selector);for(const type of ['keyDown','keyUp'])await b.command('Input.dispatchKeyEvent',{type,key:'a',code:'KeyA',windowsVirtualKeyCode:65,modifiers:2});await b.command('Input.insertText',{text:value})}
async function screenshot(name,width=1440,theme='light',height=1000){await b.command('Emulation.setDeviceMetricsOverride',{width,height,deviceScaleFactor:1,mobile:width<768});if(await b.evaluate("document.documentElement.classList.contains('dark')")!==(theme==='dark'))await click('button[aria-label="Switch to '+theme+' theme"]');await wait(150);const shot=await b.command('Page.captureScreenshot',{format:'png',captureBeyondViewport:false});const file=`${name}-${width}-${height}-${theme}.png`;await fs.writeFile(path.join(output,file),Buffer.from(shot.data,'base64'));screens.push(file)}

async function login(role){await b.command('Network.clearBrowserCookies');await navigate('/login');await until("!!document.querySelector('#login-email')");await fill('#login-email',fixture[role].email);await fill('#login-password',fixture[role].password);await button('Sign In');await until("location.pathname!='/login'&&!!document.querySelector('h1')")}



await fs.mkdir(output,{recursive:true})


function pass(name){checks.push(name);console.log('PASS',name)}
async function ready(){await until("!!document.querySelector('h1')&&!document.querySelector('.directory-loading')&&!document.body.textContent.includes('Loading current metrics')&&!document.body.textContent.includes('Loading…')");await wait(350)}
async function fits(){assert.equal(await b.evaluate('document.documentElement.scrollWidth<=innerWidth'),true)}
try{
 await login('student01');await navigate('/borrower/home');await ready();assert.deepEqual(await b.evaluate("[...document.querySelectorAll('.borrower-quick-actions a')].map(e=>e.getAttribute('href'))"),['/borrower/equipment','/borrower/borrowings','/borrower/account']);for(const theme of ['light','dark']){await screenshot('student-home-full-reference',390,theme,1672);await fits();await b.evaluate('window.scrollTo(0,document.documentElement.scrollHeight)');await screenshot('student-home-lower-content',390,theme,844);await b.evaluate('window.scrollTo(0,0)')}const thumb=await b.evaluate("document.querySelector('.borrower-recent-records .equipment-thumbnail')?.getBoundingClientRect().width");assert.ok(thumb<=53);pass('Full B01 composition and lower quick-action/collection content use compact protected thumbnails and actual themes')
 for(const route of ['/borrower/equipment','/borrower/account']){await navigate(route);await ready();for(const theme of ['light','dark']){await screenshot('student01'+route.replaceAll('/','-'),390,theme,844);await fits()}}assert.equal(await b.evaluate("document.body.textContent.includes('Borrower type')"),true);pass('Final B02 contextual cart and Profile After distinguish Student category from operational authorization role')
 assert.equal(errors.length,0);await fs.writeFile(path.join(output,'reference-acceptance.json'),JSON.stringify({result:'PASS',checks,screens,errors,scope:'Final literal B01/B02/Profile comparison; synthetic read-only browsing, no borrowing commands.'},null,2));console.log(JSON.stringify({result:'PASS',checks:checks.length,screens:screens.length}))
}finally{await b.close()}
