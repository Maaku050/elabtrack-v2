#!/usr/bin/env python3
"""Guarded offline presentation demo. Never writes normal or accepted Phase7 data."""
import argparse
import hashlib
import io
import json
import os
from pathlib import Path
import secrets
import shutil
import signal
import socket
import subprocess
import sys
import tarfile
import time
from urllib.parse import quote
from urllib.request import urlopen

ROOT=Path(__file__).resolve().parent.parent
PRIVATE=ROOT/'backend/tmp/presentation-demo'
BASELINE=PRIVATE/'snapshot'
DB='elabtrack_v2_presentation_demo'
PORT='54836'
CONTAINER='elabtrack_v2_presentation_demo_postgres'
GO=Path(os.environ.get('ELABTRACK_GO') or shutil.which('go') or 'go')
NODE=Path(os.environ.get('ELABTRACK_NODE') or shutil.which('node') or 'node')


def protected(path,data):
    path.parent.mkdir(parents=True,exist_ok=True,mode=0o700)
    fd=os.open(path,os.O_WRONLY|os.O_CREAT|os.O_TRUNC,0o600)
    os.fchmod(fd,0o600)
    with os.fdopen(fd,'w') as f:f.write(data)


def read_private(path):
    if path.stat().st_mode&0o077:raise RuntimeError('Private file permissions must be0600')
    return json.loads(path.read_text())


def run(args,cwd=ROOT,env=None,input=None):
    r=subprocess.run([str(x) for x in args],cwd=cwd,env=env,input=input,stdout=subprocess.PIPE,stderr=subprocess.PIPE)
    if r.returncode:
        protected(PRIVATE/'last-error.log',r.stderr.decode(errors='replace'))
        raise RuntimeError('Command failed; private last-error.log retained. Credentials are not printed.')
    return r.stdout


def target():
    v=read_private(PRIVATE/'config.json')
    if v['DB_NAME']!=DB or v['DB_PORT']!=PORT or v['DB_HOST']!='127.0.0.1' or v['APP_ENV']!='development':raise RuntimeError('Wrong demo target')
    info=read_private(PRIVATE/'target.json')
    if info['container']!=CONTAINER or not info['volume'].startswith('elabtrack_v2_presentation_demo_') or info['database']!=DB or info.get('disposable') is not True:raise RuntimeError('Wrong disposable identity')
    i=json.loads(run(['docker','inspect',CONTAINER]))[0]
    if i['Config']['Labels'].get('com.docker.compose.project')!='elabtrack-presentation-demo':raise RuntimeError('Wrong Compose ownership')
    if [m['Name'] for m in i['Mounts'] if m['Type']=='volume']!=[info['volume']]:raise RuntimeError('Wrong volume')
    ports=i['HostConfig']['PortBindings'].get('5432/tcp',[])
    if ports!=[{'HostIp':'127.0.0.1','HostPort':PORT}]:raise RuntimeError('Wrong loopback port')
    return v,info


def environment(migration=False):
    v,_=target()
    e={k:x for k,x in os.environ.items() if not k.startswith('PG') and k not in ('DATABASE_URL','MIGRATION_DATABASE_URL','PORT','BREVO_API_KEY','BREVO_SENDER_EMAIL')}
    e.update(v)
    e.update(APP_BIND_HOST='127.0.0.1',GOCACHE=os.environ.get('ELABTRACK_BUILD_CACHE','/tmp/elabtrack-phase1b-go-cache'),ELABTRACK_PRESENTATION_DEMO='1',BREVO_API_KEY='',BREVO_SENDER_EMAIL='',VITE_API_URL='http://127.0.0.1:18087/api/v1')
    if migration:e['MIGRATION_DATABASE_URL']='postgresql://elabtrack_migrator:'+quote(v['MIGRATION_DB_PASSWORD'],safe='')+'@127.0.0.1:'+PORT+'/'+DB+'?sslmode=disable'
    else:
        e.pop('MIGRATION_DB_PASSWORD',None);e.pop('BOOTSTRAP_DB_PASSWORD',None)
    return e


def sql(statement):
    target()
    command='PGPASSWORD="$MIGRATION_DB_PASSWORD" exec psql -X -h127.0.0.1 -Uelabtrack_migrator -d'+DB+' -vON_ERROR_STOP=1 -t -A'
    return run(['docker','exec','-i',CONTAINER,'sh','-c',command],input=statement.encode()).decode()


def snapshot_baseline():
    # Explicit allowlist; private configs, fixtures, ignored files and references never copied.
    BASELINE.mkdir(parents=True,exist_ok=True,mode=0o700)
    entries=['backend/cmd/api','backend/internal','backend/migrations','backend/database','backend/vendor','frontend/src','frontend/public']
    for name in entries:
        src=ROOT/name;dst=BASELINE/name
        if src.exists():shutil.copytree(src,dst,dirs_exist_ok=True,ignore=shutil.ignore_patterns('__pycache__','*.log'))
    for name in ['backend/go.mod','backend/go.sum','frontend/package.json','frontend/package-lock.json','frontend/index.html','frontend/components.json','frontend/tsconfig.json','frontend/tsconfig.app.json','frontend/tsconfig.node.json','frontend/vite.config.ts']:
        src=ROOT/name
        if src.exists():
            dst=BASELINE/name;dst.parent.mkdir(parents=True,exist_ok=True);shutil.copyfile(src,dst)
    for file in (ROOT/'backend/migrations').glob('00000[1-8]*.sql'):
        original=run(['git','show','c2741c5:backend/migrations/'+file.name])
        if hashlib.sha256(original).digest()!=hashlib.sha256(file.read_bytes()).digest():raise RuntimeError('Historical migration changed')
    hashes={str(f.relative_to(BASELINE)):hashlib.sha256(f.read_bytes()).hexdigest() for name in entries for f in (BASELINE/name).rglob('*') if f.is_file()}
    protected(PRIVATE/'source-hashes.json',json.dumps(hashes,indent=2)+'\n')
    seed=BASELINE/'backend/cmd/presentation-seed';seed.mkdir(parents=True,exist_ok=True)
    shutil.copyfile(ROOT/'integration/presentation/seed.go',seed/'main.go')


def prepare():
    PRIVATE.mkdir(parents=True,exist_ok=True,mode=0o700);PRIVATE.chmod(0o700)
    if (PRIVATE/'target.json').exists():
        target()
        n=int(sql('SELECT count(*) FROM users;').strip())
        if n:
            if n==28 and sql("SELECT seeded FROM presentation_demo_identity WHERE id=1;").strip()=='t':
                print('Presentation already initialized; no duplicate seed. Use status/start/build or deliberate reset.');return
            raise RuntimeError('Unexpected seeded target; no existing records changed')
    else:
        # A fresh working copy must never replace another presentation target.
        probe=subprocess.run(['docker','inspect',CONTAINER],stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL)
        if probe.returncode==0:raise RuntimeError('Presentation container already exists; use its owning working copy, never replace it from a fresh setup')
        snapshot_baseline()
        v={'APP_ENV':'development','APP_PORT':'18087','APP_BIND_HOST':'127.0.0.1','DB_HOST':'127.0.0.1','DB_PORT':PORT,'DB_NAME':DB,'DB_USER':'elabtrack_runtime','DB_PASSWORD':secrets.token_urlsafe(32),'DB_SSLMODE':'disable','BOOTSTRAP_DB_USER':'postgres','BOOTSTRAP_DB_PASSWORD':secrets.token_urlsafe(32),'MIGRATION_DB_PASSWORD':secrets.token_urlsafe(32),'JWT_SECRET':secrets.token_urlsafe(48),'FRONTEND_URL':'http://127.0.0.1:15177','ALLOWED_ORIGINS':'http://127.0.0.1:15177','STUDENT_EMAIL_DOMAINS':'students.example.invalid','LOG_FORMAT':'json','LOG_LEVEL':'info','RATE_LIMIT_WINDOW':'5s','NOTIFICATION_DUE_SOON_LEAD_SECONDS':'86400'}
        protected(PRIVATE/'config.json',json.dumps(v,indent=2)+'\n')
        protected(PRIVATE/'compose.env','\n'.join(k+'='+x for k,x in v.items())+'\n')
        stamp=time.strftime('%Y%m%dT%H%M%SZ',time.gmtime()).lower()
        volume='elabtrack_v2_presentation_demo_'+stamp+'_'+secrets.token_hex(3)+'_pgdata'
        protected(PRIVATE/'compose.yml','services:\n  postgres:\n    container_name: '+CONTAINER+'\n    ports: !override\n      - "127.0.0.1:'+PORT+':5432"\nvolumes:\n  elabtrack_v2_pgdata:\n    name: '+volume+'\n')
        protected(PRIVATE/'target.json',json.dumps({'container':CONTAINER,'volume':volume,'database':DB,'port':int(PORT),'baseline':'c2741c5 + verified working source','disposable':True},indent=2)+'\n')
        run(compose()+['up','-d','--wait','--pull','never','postgres'])
    target();env=environment(True)
    run([GO,'build','-o',PRIVATE/'api','./cmd/api'],cwd=BASELINE/'backend',env=env)
    migrations=sorted((BASELINE/'backend/migrations').glob('*.up.sql'))
    if len(migrations)!=11 or migrations[-1].name!='000011_profile_images.up.sql':raise RuntimeError('Presentation migration gate failed')
    run([PRIVATE/'api','--migrate-up'],cwd=BASELINE/'backend',env=env)
    status=run([PRIVATE/'api','--migrate-status'],cwd=BASELINE/'backend',env=env).decode();protected(PRIVATE/'migration-status.txt',status)
    sql((ROOT/'backend/database/runtime-grants.sql').read_text())
    sql("CREATE TABLE IF NOT EXISTS presentation_demo_identity(id integer PRIMARY KEY CHECK(id=1),identity text NOT NULL,seeded boolean NOT NULL DEFAULT false); INSERT INTO presentation_demo_identity(id,identity) VALUES(1,'elabtrack-presentation-demo') ON CONFLICT(id) DO NOTHING;")
    run([GO,'build','-o',PRIVATE/'seed','./cmd/presentation-seed'],cwd=BASELINE/'backend',env=env)
    accounts=[]
    for role,count,category,prefix in [('ADMIN',1,'','admin'),('STAFF',2,'','staff'),('BORROWER',20,'STUDENT','student'),('BORROWER',5,'FACULTY','faculty')]:
        for i in range(1,count+1):
            key='admin' if role=='ADMIN' else prefix+f'{i:02d}'
            if role=='STAFF':key='staff'+str(i)
            email=key+('@students.example.invalid' if category=='STUDENT' else '@example.invalid')
            names=['Avery','Rowan','Morgan','Sage','Quinn','River','Sky','Reese','Finley','Remy','Ellis','Robin','Casey','Cameron','Taylor','Harper','Riley','Jordan','Jamie','Alex']
            name='DEMO '+names[(i-1)%20]+' '+('Student' if category=='STUDENT' else 'Faculty' if category=='FACULTY' else prefix.title())+f' {i:02d}'
            accounts.append(dict(Key=key,Email=email,Name=name,Role=role,Category=category,StudentID='DEMO-S-'+f'{i:03d}' if category=='STUDENT' else '',Password=secrets.token_urlsafe(24)))
    credentials={x['Key']:{'email':x['Email'],'password':x['Password']} for x in accounts}
    protected(PRIVATE/'credentials.json',json.dumps(credentials,indent=2)+'\n')
    data=json.loads(run([PRIVATE/'seed'],cwd=BASELINE/'backend',env=env,input=json.dumps({'Accounts':accounts,'Assets':str(ROOT/'integration/presentation/assets')}).encode()))
    protected(PRIVATE/'fixtures.json',json.dumps(data,indent=2)+'\n')
    fixture={**data,**{k:{**data['accounts'][k],**credentials[k]} for k in credentials}}
    protected(PRIVATE/'browser-fixtures.json',json.dumps(fixture,indent=2)+'\n')
    build_frontend()
    print('Prepared separate presentation:40 equipment,20 Students,5 Faculty,2 Staff,1 Admin,68 bundled illustrations. Student20 retains actual first-use terms workflow.')


def build_frontend():
    snapshot_baseline();env=environment()
    run([GO,'build','-o',PRIVATE/'api','./cmd/api'],cwd=BASELINE/'backend',env=environment(True))
    node_modules=BASELINE/'frontend/node_modules'
    if not node_modules.exists():node_modules.symlink_to(ROOT/'frontend/node_modules',target_is_directory=True)
    config=(BASELINE/'frontend/vite.config.ts').read_text().replace('plugins: [react(), tailwindcss()],','''plugins: [react(), tailwindcss(), {name:'isolated-presentation',transformIndexHtml(html:string){return html.replace('</head>','<style>:root{--presentation-banner-height:28px}.borrower-bottom-nav{bottom:28px}body{padding-bottom:28px}</style></head>').replace('<body>', '<body><div role="status" style="position:fixed;bottom:0;left:0;right:0;z-index:9999;background:#20204b;color:#fff;text-align:center;font:12px system-ui;padding:6px;pointer-events:none">FSMO PRESENTATION · FICTIONAL DATA · PORT 15177</div>')}}],''')
    protected(BASELINE/'frontend/vite.demo.config.ts',config)
    env['PATH']=str(NODE.parent)+os.pathsep+env.get('PATH','')
    run([NODE,ROOT/'frontend/node_modules/vite/bin/vite.js','build','--config','vite.demo.config.ts'],cwd=BASELINE/'frontend',env=env)


def compose():return ['docker','compose','--env-file',PRIVATE/'compose.env','-p','elabtrack-presentation-demo','-f',ROOT/'docker-compose.yml','-f',PRIVATE/'compose.yml']


def stop():
    if not (PRIVATE/'processes.json').exists():return
    records=read_private(PRIVATE/'processes.json')
    for key,r in records.items():
        pid=r['pid'];proc=Path('/proc')/str(pid)/'cmdline'
        if not proc.exists():continue
        argv=proc.read_bytes().split(b'\0')
        if str(PRIVATE).encode() not in b' '.join(argv):raise RuntimeError('Refusing to stop unexpected process')
        os.killpg(pid,signal.SIGINT)
        for _ in range(150):
            if not proc.exists() or (proc.parent/'stat').read_text().split()[2]=='Z':break
            time.sleep(.1)
        else:raise RuntimeError('Owned process did not stop')
    (PRIVATE/'processes.json').unlink(missing_ok=True)
    print('Owned demo API/frontend stopped; databases and volumes retained.')


def start():
    target();run(compose()+['up','-d','--wait','--pull','never','postgres'])
    if (PRIVATE/'processes.json').exists():
        old=read_private(PRIVATE/'processes.json')
        if all((Path('/proc')/str(r['pid'])/'stat').exists() and (Path('/proc')/str(r['pid'])/'stat').read_text().split()[2]!='Z' for r in old.values()):print('Demo already running: http://127.0.0.1:15177/login');return
        stop()
    for port in (18087,15177):
        with socket.socket(socket.AF_INET,socket.SOCK_STREAM) as listener:
            listener.setsockopt(socket.SOL_SOCKET,socket.SO_REUSEADDR,1)
            try:listener.bind(('127.0.0.1',port))
            except OSError:raise RuntimeError('Presentation port '+str(port)+' is occupied; no other process is stopped')
    e=environment();e['PATH']=str(NODE.parent)+os.pathsep+e.get('PATH','');records={}
    try:
        for key,args,cwd in [('api',[PRIVATE/'api'],BASELINE/'backend'),('frontend',[NODE,ROOT/'frontend/node_modules/vite/bin/vite.js','preview','--config',BASELINE/'frontend/vite.demo.config.ts','--host','127.0.0.1','--port','15177','--strictPort'],BASELINE/'frontend')]:
            log=PRIVATE/(key+'.log');fd=os.open(log,os.O_WRONLY|os.O_CREAT|os.O_APPEND,0o600)
            with os.fdopen(fd,'w') as f:p=subprocess.Popen([str(x) for x in args],cwd=cwd,env=e,stdout=f,stderr=subprocess.STDOUT,start_new_session=True)
            records[key]={'pid':p.pid,'command':str(args[0])};protected(PRIVATE/'processes.json',json.dumps(records))
        for url in ['http://127.0.0.1:18087/api/v1/ready','http://127.0.0.1:15177/login']:
            for _ in range(100):
                try:
                    with urlopen(url,timeout=1) as r:assert r.status==200
                    break
                except Exception:time.sleep(.1)
            else:raise RuntimeError('Demo health gate failed')
    except Exception:
        stop();raise
    print('DEMO ready: http://127.0.0.1:15177/login · API18087 · PostgreSQL54836 · normal ports unchanged.')


def reset(confirm):
    target()
    if confirm!=DB:raise RuntimeError('Reset requires --confirm '+DB+'; only this disposable demo can be reset')
    if sql("SELECT identity||':'||seeded FROM presentation_demo_identity WHERE id=1;").strip()!='elabtrack-presentation-demo:true':raise RuntimeError('Disposable database identity not verified')
    stop();run(compose()+['down']) # Never -v. Retain prior scenario volume for recovery.
    archive=PRIVATE/('retained-config-'+time.strftime('%Y%m%dT%H%M%SZ',time.gmtime()));archive.mkdir(mode=0o700)
    for name in ['config.json','compose.env','compose.yml','target.json','credentials.json','fixtures.json','browser-fixtures.json','migration-status.txt']:
        p=PRIVATE/name
        if p.exists():p.rename(archive/name)
    prepare();print('Demo reset uses a new verified isolated volume; previous volume retained, normal data untouched. Run start.')


def credentials():
    target()
    if not sys.stdin.isatty() or not sys.stdout.isatty():raise RuntimeError('Retrieve credentials only in your own interactive terminal; never paste them into chat/logs')
    v=read_private(PRIVATE/'credentials.json')
    print('Private DEMO credentials only. Do not share or commit.')
    for key,x in v.items():print(key+': '+x['email']+'  password: '+x['password'])


def status():
    target()
    out=sql("BEGIN READ ONLY;SELECT json_build_object('database',current_database(),'migrations',(SELECT count(*) FROM schema_migrations),'accounts',(SELECT count(*) FROM users),'students',(SELECT count(*) FROM borrower_profiles WHERE borrower_type='STUDENT'),'faculty',(SELECT count(*) FROM borrower_profiles WHERE borrower_type='FACULTY'),'staff',(SELECT count(*) FROM users WHERE role='STAFF'),'admin',(SELECT count(*) FROM users WHERE role='ADMIN'),'equipment',(SELECT count(*) FROM equipment),'images',(SELECT count(*) FROM equipment_images),'avatars',(SELECT count(*) FROM profile_images WHERE image_id IS NOT NULL),'notifications',(SELECT count(*) FROM notifications),'terms_acceptances',(SELECT count(*) FROM terms_acceptances),'borrowings',(SELECT count(*) FROM borrowings),'stock',(SELECT json_build_object('available',sum(available),'reserved',sum(reserved),'checked_out',sum(checked_out),'damaged_held',sum(damaged_held),'total',sum(total_tracked),'violations',count(*) FILTER(WHERE total_tracked<>available+reserved+checked_out+damaged_held)) FROM equipment));COMMIT;")
    print(next(l for l in out.splitlines() if l.startswith('{')))



def main():
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('command',choices=['setup','start','stop','status','credentials','reset','build']);p.add_argument('--confirm');a=p.parse_args()
    if a.command=='reset':reset(a.confirm)
    else:globals()[{'setup':'prepare','build':'build_frontend'}.get(a.command,a.command)]()
if __name__=='__main__':
    try:main()
    except Exception as e:raise SystemExit(str(e))
