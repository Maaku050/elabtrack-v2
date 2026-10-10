import { spawn } from 'node:child_process';
export async function browser({ showScrollbars = false } = {}) {
 const child=spawn(process.env.CHROMIUM_EXECUTABLE || 'chromium',['--headless',...(showScrollbars ? [] : ['--hide-scrollbars']),'--no-sandbox','--disable-gpu','--disable-dev-shm-usage','--disable-crash-reporter','--remote-debugging-pipe'],{env:process.env,stdio:['ignore','ignore','pipe','pipe','pipe']});
 let buffer='',id=0;const pending=new Map(),listeners=[];
 child.stderr.on('data',()=>{});
 child.on('error',error=>{for(const p of pending.values())p.reject(error);});
 for(const stream of [child.stdio[3],child.stdio[4]])stream.on('error',error=>{for(const p of pending.values())p.reject(error);});
 child.on('exit',code=>{for(const p of pending.values())p.reject(new Error(`Browser exited ${code}`));});
 child.stdio[4].on('data',chunk=>{buffer+=chunk;let i;while((i=buffer.indexOf('\0'))>=0){const msg=JSON.parse(buffer.slice(0,i));buffer=buffer.slice(i+1);if(msg.id){const p=pending.get(msg.id);pending.delete(msg.id);if(msg.error)p?.reject(new Error(JSON.stringify(msg.error)));else p?.resolve(msg.result);}else listeners.forEach(fn=>fn(msg));}});
 function send(method,params={},sessionId){return new Promise((resolve,reject)=>{const key=++id;pending.set(key,{resolve,reject});child.stdio[3].write(JSON.stringify({id:key,method,params,sessionId})+'\0');});}
 const {targetId}=await send('Target.createTarget',{url:'about:blank'});const {sessionId}=await send('Target.attachToTarget',{targetId,flatten:true});
 const command=(method,params={},targetSession=sessionId)=>send(method,params,targetSession);
 await command('Page.enable');await command('Runtime.enable');
 return {command,on:fn=>listeners.push(fn),close:()=>send('Browser.close'),evaluate:async expression=>{const r=await command('Runtime.evaluate',{expression,awaitPromise:true,returnByValue:true});if(r.exceptionDetails)throw new Error(JSON.stringify(r.exceptionDetails));return r.result.value;}};
}
