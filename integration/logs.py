"""Check captured real API output without displaying tested credentials."""
import json
from pathlib import Path

paths=[Path('/tmp/elabtrack-phase1h-startup.log'),Path('/tmp/elabtrack-phase1h-final.log')]
assert paths[-1].exists(), 'capture the current integration API log first'
paths=[p for p in paths if p.exists()]
raw='\n'.join(p.read_text() for p in paths)
secrets=json.loads(Path('/tmp/elabtrack-phase1g-secrets.json').read_text())
fixture=Path('backend/.env.phase1h')
if fixture.exists():
    values=dict(line.split('=',1) for line in fixture.read_text().splitlines() if line and not line.startswith('#'))
    secrets += [value for key,value in values.items() if key.endswith('_PASSWORD') or key=='JWT_SECRET']
secrets += ['wrong-synthetic-password','Authorization: Bearer','Cookie:','Set-Cookie:']
assert not any(s and s in raw for s in secrets), 'credential sentinel leaked (value withheld)'
rows=[json.loads(line) for p in paths for line in p.read_text().splitlines() if line.startswith('{')]
forbidden={'password','password_hash','token','access_token','refresh_token','token_hash','authorization','cookie','set-cookie','database_url','db_password','jwt_secret','body','headers','query','url'}
assert not any(forbidden.intersection(k.lower() for k in row) for row in rows)
current=[json.loads(line) for line in dict.fromkeys(paths[-1].read_text().splitlines()) if line.startswith('{')]
completed={}
for row in current:
    if row.get('event')=='http.request_completed':completed.setdefault(row['request_id'],[]).append(row)
browser=json.loads(Path('/tmp/elabtrack-phase1g-browser.json').read_text())
runtime=json.loads(Path('/tmp/elabtrack-phase1g-runtime.json').read_text())
for response in browser['correlation']+runtime['requests']:
    ident=response.get('id')
    if ident:
        match=completed.get(ident,[])
        assert len(match)==1 and match[0]['status']==response['status'], 'final response/log correlation mismatch'
proxy={}
for label in ['direct','direct_spoof','proxied_spoof']:
    row=completed[runtime['checks'][label]['requestId']][0]
    proxy[label]=row['client_ip']
assert proxy['direct']==proxy['direct_spoof']
assert '198.51.100.42' not in proxy.values()
events=sorted({r.get('security_event') for r in rows if r.get('security_event')})
assert {'auth.login_succeeded','auth.login_failed','auth.refresh_succeeded','auth.refresh_failed','auth.logout_completed'}.issubset(events)
statuses=sorted({r['status'] for r in current if 'status' in r})
assert {200,401,403,413,429,431,503}.issubset(statuses)
result={'secretSentinelsChecked':len(set(secrets)),'secretMatches':0,'structuredRowsChecked':len(rows),'correlatedResponses':len(browser['correlation'])+len(runtime['requests']),'securityEvents':events,'finalStatuses':statuses,'clientIPs':proxy,'parserCompletionsAccurate':True}
Path('/tmp/elabtrack-phase1i-logs.json').write_text(json.dumps(result,indent=2))
print(json.dumps(result,indent=2))
