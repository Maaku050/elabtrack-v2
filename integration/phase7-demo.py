#!/usr/bin/env python3
"""Explicit, guarded Phase7 owner demo. No normal data, mail or auth bypass."""
import argparse
import hashlib
import io
import json
import os
from pathlib import Path
import secrets
import shutil
import signal
import subprocess
import sys
import tarfile
import time
from urllib.parse import quote
from urllib.request import urlopen

ROOT=Path(__file__).resolve().parent.parent
PRIVATE=ROOT/'backend/tmp/phase7-owner-demo'
BASELINE=PRIVATE/'baseline'
DB='elabtrack_v2_phase7_demo'
PORT='54835'
CONTAINER='elabtrack_v2_phase7_owner_demo_postgres'
GO=Path(shutil.which('go') or '/home/marvin/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.27.1.linux-amd64/bin/go')
NODE=Path(shutil.which('node') or '/home/marvin/.nvm/versions/node/v24.19.0/bin/node')


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
    if info['container']!=CONTAINER or not info['volume'].startswith('elabtrack_v2_phase7_owner_demo_') or info['database']!=DB or info.get('disposable') is not True:raise RuntimeError('Wrong disposable identity')
    i=json.loads(run(['docker','inspect',CONTAINER]))[0]
    if i['Config']['Labels'].get('com.docker.compose.project')!='elabtrack-phase7-owner-demo':raise RuntimeError('Wrong Compose ownership')
    if [m['Name'] for m in i['Mounts'] if m['Type']=='volume']!=[info['volume']]:raise RuntimeError('Wrong volume')
    ports=i['HostConfig']['PortBindings'].get('5432/tcp',[])
    if ports!=[{'HostIp':'127.0.0.1','HostPort':PORT}]:raise RuntimeError('Wrong loopback port')
    return v,info


def environment(migration=False):
    v,_=target()
    e={k:x for k,x in os.environ.items() if not k.startswith('PG') and k not in ('DATABASE_URL','MIGRATION_DATABASE_URL','PORT','BREVO_API_KEY','BREVO_SENDER_EMAIL')}
    e.update(v)
    e.update(GOCACHE='/tmp/elabtrack-phase1b-go-cache',ELABTRACK_PHASE7_DEMO='1',BREVO_API_KEY='',BREVO_SENDER_EMAIL='',VITE_API_URL='http://127.0.0.1:18086/api/v1')
    if migration:e['MIGRATION_DATABASE_URL']='postgresql://elabtrack_migrator:'+quote(v['MIGRATION_DB_PASSWORD'],safe='')+'@127.0.0.1:'+PORT+'/'+DB+'?sslmode=disable'
    else:
        e.pop('MIGRATION_DB_PASSWORD',None);e.pop('BOOTSTRAP_DB_PASSWORD',None)
    return e


def sql(statement):
    target()
    command='PGPASSWORD="$MIGRATION_DB_PASSWORD" exec psql -X -h127.0.0.1 -Uelabtrack_migrator -d'+DB+' -vON_ERROR_STOP=1 -t -A'
    return run(['docker','exec','-i',CONTAINER,'sh','-c',command],input=statement.encode()).decode()


def snapshot_baseline():
    if BASELINE.exists():
        if (PRIVATE/'baseline-commit.txt').exists() and (PRIVATE/'baseline-commit.txt').read_text().strip()!='73f3fe4':raise RuntimeError('Unexpected snapshot')
        return
    BASELINE.mkdir(parents=True,mode=0o700)
    data=run(['git','archive','73f3fe4','backend','frontend','integration'])
    with tarfile.open(fileobj=io.BytesIO(data)) as a:a.extractall(BASELINE,filter='data')
    protected(PRIVATE/'baseline-commit.txt','73f3fe4\n')


def prepare():
    PRIVATE.mkdir(parents=True,exist_ok=True,mode=0o700);PRIVATE.chmod(0o700)
    if (PRIVATE/'target.json').exists():
        target()
        n=int(sql('SELECT count(*) FROM users;').strip())
        if n:
            if n==5 and (BASELINE/'frontend/dist/index.html').exists() and (PRIVATE/'browser-fixtures.json').exists():
                print('Demo already initialized; nothing seeded or reset. Use status/start or explicit reset.');return
            raise RuntimeError('Existing partial seeded demo requires explicit guarded reset; no duplicate accounts created')
        print('Resuming verified empty demo initialization; no existing accounts changed.')
    else:
        snapshot_baseline();protected(PRIVATE/'baseline-commit.txt','73f3fe4\n')
        v={'APP_ENV':'development','APP_PORT':'18086','DB_HOST':'127.0.0.1','DB_PORT':PORT,'DB_NAME':DB,'DB_USER':'elabtrack_runtime','DB_PASSWORD':secrets.token_urlsafe(32),'DB_SSLMODE':'disable','BOOTSTRAP_DB_USER':'postgres','BOOTSTRAP_DB_PASSWORD':secrets.token_urlsafe(32),'MIGRATION_DB_PASSWORD':secrets.token_urlsafe(32),'JWT_SECRET':secrets.token_urlsafe(48),'FRONTEND_URL':'http://127.0.0.1:15176','ALLOWED_ORIGINS':'http://127.0.0.1:15176','STUDENT_EMAIL_DOMAINS':'students.example.invalid','LOG_FORMAT':'json','LOG_LEVEL':'info','RATE_LIMIT_WINDOW':'5s'}
        protected(PRIVATE/'config.json',json.dumps(v,indent=2)+'\n')
        protected(PRIVATE/'compose.env','\n'.join(k+'='+x for k,x in v.items())+'\n')
        stamp=time.strftime('%Y%m%dT%H%M%SZ',time.gmtime()).lower()
        volume='elabtrack_v2_phase7_owner_demo_'+stamp+'_'+secrets.token_hex(3)+'_pgdata'
        protected(PRIVATE/'compose.yml','services:\n  postgres:\n    container_name: '+CONTAINER+'\n    ports: !override\n      - "127.0.0.1:'+PORT+':5432"\nvolumes:\n  elabtrack_v2_pgdata:\n    name: '+volume+'\n')
        info={'container':CONTAINER,'volume':volume,'database':DB,'port':int(PORT),'baseline':'73f3fe4','disposable':True}
        protected(PRIVATE/'target.json',json.dumps(info,indent=2)+'\n')
        run(compose()+['up','-d','--wait','postgres'])
    target();env=environment(True)
    run([GO,'build','-o',PRIVATE/'api','./cmd/api'],cwd=BASELINE/'backend',env=env)
    migrations=sorted((BASELINE/'backend/migrations').glob('*.up.sql'))
    if len(migrations)!=8 or migrations[-1].name!='000008_borrowing.up.sql':raise RuntimeError('Snapshot migration gate failed')
    for p in (BASELINE/'backend/migrations').glob('*.sql'):
        if hashlib.sha256(p.read_bytes()).digest()!=hashlib.sha256((ROOT/'backend/migrations'/p.name).read_bytes()).digest():raise RuntimeError('Historical migration changed')
    run([PRIVATE/'api','--migrate-up'],cwd=BASELINE/'backend',env=env)
    status=run([PRIVATE/'api','--migrate-status'],cwd=BASELINE/'backend',env=env).decode();protected(PRIVATE/'migration-status.txt',status)
    if '000008_borrowing' not in status or '000009' in status:raise RuntimeError('Migration verification failed')
    sql((ROOT/'backend/database/runtime-grants.sql').read_text())
    sql("CREATE TABLE IF NOT EXISTS phase7_demo_identity(id integer PRIMARY KEY CHECK(id=1),identity text NOT NULL,seeded boolean NOT NULL DEFAULT false); INSERT INTO phase7_demo_identity(id,identity) VALUES(1,'elabtrack-phase7-owner-demo') ON CONFLICT(id) DO NOTHING;")
    seed=BASELINE/'backend/cmd/phase7-demo';seed.mkdir(exist_ok=True);shutil.copyfile(ROOT/'integration/phase7-demo/seed.go',seed/'main.go')
    run([GO,'build','-o',PRIVATE/'seed','./cmd/phase7-demo'],cwd=BASELINE/'backend',env=env)
    accounts=[]
    for key,email,name,role,category,sid in [('student','student.one@students.example.invalid','DEMO Student One','BORROWER','STUDENT','DEMO-001'),('student2','student.two@students.example.invalid','DEMO Student Two','BORROWER','STUDENT','DEMO-002'),('faculty','faculty.one@example.invalid','DEMO Faculty One','BORROWER','FACULTY',''),('staff','staff.one@example.invalid','DEMO Staff One','STAFF','',''),('admin','admin.one@example.invalid','DEMO Admin One','ADMIN','','')]:
        accounts.append(dict(Key=key,Email=email,Name=name,Role=role,Category=category,StudentID=sid,Password=secrets.token_urlsafe(24)))
    # Save private credentials before atomic seeding so a later disk failure cannot strand passwords.
    credentials={x['Key']:{'email':x['Email'],'password':x['Password']} for x in accounts}
    protected(PRIVATE/'credentials.json',json.dumps(credentials,indent=2)+'\n')
    data=json.loads(run([PRIVATE/'seed'],cwd=BASELINE/'backend',env=env,input=json.dumps({'Accounts':accounts}).encode()))
    protected(PRIVATE/'fixtures.json',json.dumps(data,indent=2)+'\n')
    fixture={**data,**{k:{**data[k],**credentials[k]} for k in credentials}};fixture['borrower']=fixture['student']
    protected(PRIVATE/'browser-fixtures.json',json.dumps(fixture,indent=2)+'\n')
    build_frontend()
    print('Prepared isolated Phase7 demo:5 accounts,8 equipment,3 categories,8 local images; zero seeded acceptances; migrations000001–000008 verified.')


def build_frontend():
    env=environment()
    node_modules=BASELINE/'frontend/node_modules'
    if not node_modules.exists():node_modules.symlink_to(ROOT/'frontend/node_modules',target_is_directory=True)
    config=(BASELINE/'frontend/vite.config.ts').read_text().replace('plugins: [react(), tailwindcss()],','''plugins: [react(), tailwindcss(), {name:'phase7-isolated-demo',transformIndexHtml(html:string){return html.replace('</head>','<style>.borrower-bottom-nav{bottom:26px}.cart-dock{bottom:calc(var(--borrower-nav-height) + env(safe-area-inset-bottom) + 34px)}body{padding-bottom:26px}</style></head>').replace('<body>', '<body><div role="status" style="position:fixed;bottom:0;left:0;right:0;z-index:9999;background:#20204b;color:#fff;text-align:center;font:12px system-ui;padding:5px;pointer-events:none">ISOLATED PHASE 7 DEMO · FICTIONAL DATA</div>')}}],''')
    protected(BASELINE/'frontend/vite.demo.config.ts',config)
    env['PATH']=str(NODE.parent)+os.pathsep+env.get('PATH','')
    run([NODE,ROOT/'frontend/node_modules/vite/bin/vite.js','build','--config','vite.demo.config.ts'],cwd=BASELINE/'frontend',env=env)

def compose():return ['docker','compose','--env-file',PRIVATE/'compose.env','-p','elabtrack-phase7-owner-demo','-f',ROOT/'docker-compose.yml','-f',PRIVATE/'compose.yml']


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
    target();run(compose()+['up','-d','--wait','postgres'])
    if (PRIVATE/'processes.json').exists():
        old=read_private(PRIVATE/'processes.json')
        if all((Path('/proc')/str(r['pid'])/'stat').exists() and (Path('/proc')/str(r['pid'])/'stat').read_text().split()[2]!='Z' for r in old.values()):print('Demo already running: http://127.0.0.1:15176/login');return
        stop()
    e=environment();e['PATH']=str(NODE.parent)+os.pathsep+e.get('PATH','');records={}
    try:
        for key,args,cwd in [('api',[PRIVATE/'api'],BASELINE/'backend'),('frontend',[NODE,ROOT/'frontend/node_modules/vite/bin/vite.js','preview','--config',BASELINE/'frontend/vite.demo.config.ts','--host','127.0.0.1','--port','15176','--strictPort'],BASELINE/'frontend')]:
            log=PRIVATE/(key+'.log');fd=os.open(log,os.O_WRONLY|os.O_CREAT|os.O_APPEND,0o600)
            with os.fdopen(fd,'w') as f:p=subprocess.Popen([str(x) for x in args],cwd=cwd,env=e,stdout=f,stderr=subprocess.STDOUT,start_new_session=True)
            records[key]={'pid':p.pid,'command':str(args[0])};protected(PRIVATE/'processes.json',json.dumps(records))
        for url in ['http://127.0.0.1:18086/api/v1/ready','http://127.0.0.1:15176/login']:
            for _ in range(100):
                try:
                    with urlopen(url,timeout=1) as r:assert r.status==200
                    break
                except Exception:time.sleep(.1)
            else:raise RuntimeError('Demo health gate failed')
    except Exception:
        stop();raise
    print('DEMO ready: http://127.0.0.1:15176/login · API18086 · PostgreSQL54835 · normal ports unchanged.')


def reset(confirm):
    target()
    if confirm!=DB:raise RuntimeError('Reset requires --confirm '+DB+'; only this disposable demo can be reset')
    if sql("SELECT identity||':'||seeded FROM phase7_demo_identity WHERE id=1;").strip()!='elabtrack-phase7-owner-demo:true':raise RuntimeError('Disposable database identity not verified')
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
    out=sql("BEGIN READ ONLY;SELECT json_build_object('database',current_database(),'migrations',(SELECT count(*) FROM schema_migrations),'accounts',(SELECT count(*) FROM users),'students',(SELECT count(*) FROM borrower_profiles WHERE borrower_type='STUDENT'),'faculty',(SELECT count(*) FROM borrower_profiles WHERE borrower_type='FACULTY'),'staff',(SELECT count(*) FROM users WHERE role='STAFF'),'admin',(SELECT count(*) FROM users WHERE role='ADMIN'),'equipment',(SELECT count(*) FROM equipment),'images',(SELECT count(*) FROM equipment_images),'terms_acceptances',(SELECT count(*) FROM terms_acceptances),'borrowings',(SELECT count(*) FROM borrowings),'stock',(SELECT json_build_object('available',sum(available),'reserved',sum(reserved),'checked_out',sum(checked_out),'damaged_held',sum(damaged_held),'total',sum(total_tracked),'violations',count(*) FILTER(WHERE total_tracked<>available+reserved+checked_out+damaged_held)) FROM equipment));COMMIT;")
    print(next(l for l in out.splitlines() if l.startswith('{')))


def expire():
    target();data=run([PRIVATE/'seed','expire'],cwd=BASELINE/'backend',env=environment(True));protected(PRIVATE/'expiry-result.json',data.decode());print('Deterministic demo request expired; actual reservation released once after two live sweeps. Open http://127.0.0.1:15176/borrower/borrowings/'+json.loads(data)['id']+' as Student One.')


def main():
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('command',choices=['setup','start','stop','status','credentials','expire','reset','build']);p.add_argument('--confirm');a=p.parse_args()
    if a.command=='reset':reset(a.confirm)
    else:globals()[{'setup':'prepare','build':'build_frontend'}.get(a.command,a.command)]()
if __name__=='__main__':
    try:main()
    except Exception as e:raise SystemExit(str(e))
