"""Real API eligibility checks against only the guarded owner demo."""
import importlib.util
import json
from pathlib import Path
import uuid
from urllib.request import Request,urlopen
from urllib.error import HTTPError
spec=importlib.util.spec_from_file_location('phase7_demo_tool',Path(__file__).resolve().parents[1]/'phase7-demo.py');demo=importlib.util.module_from_spec(spec);spec.loader.exec_module(demo)
demo.target();fixture=demo.read_private(demo.PRIVATE/'browser-fixtures.json');base='http://127.0.0.1:18086/api/v1';origin='http://127.0.0.1:15176';checks=[];tokens={};rejection_codes={}

def call(route,role=None,method='GET',body=None,access=None):
    headers={'Origin':origin,'Content-Type':'application/json','Idempotency-Key':str(uuid.uuid4())}
    if access:headers['Authorization']='Bearer '+access
    elif role:headers['Authorization']='Bearer '+token(role)
    req=Request(base+route,headers=headers,method=method,data=None if body is None else json.dumps(body).encode())
    try:
        with urlopen(req,timeout=15) as r:return r.status,json.loads(r.read())
    except HTTPError as e:return e.code,json.loads(e.read())

def token(role):
    if role in tokens:return tokens[role]
    s,r=call('/auth/login',method='POST',body={'email':fixture[role]['email'],'password':fixture[role]['password']})
    assert s==200,('legitimate login',role,s)
    tokens[role]=r['data']['access_token']
    return tokens[role]

def check(name,s,expected):
    assert s==expected,(name,s,expected);checks.append(name);print('PASS '+name)

def flag(column,value):
    assert column in ['activation_required','is_active']
    ident=str(uuid.UUID(fixture['student2']['id']))
    demo.sql('UPDATE users SET '+column+'='+('true' if value else 'false')+" WHERE id='"+ident+"';")

assert call('/terms/status','student')[1]['data']['state']=='accepted','Run borrower browser consent against the current demo version before API checks'
assert call('/terms/status','student2')[1]['data']['state']=='required','Student2 must remain unaccepted for negative eligibility checks; use a fresh reset'

before=call('/equipment/'+fixture['equipment']['id'],'staff')[1]['data']['stock']
item={'equipment_id':fixture['equipment']['id'],'quantity':1};request={'items':[item],'confirm':True}
check('Missing terms acceptance rejected by authoritative API',call('/borrowings','student2','POST',request)[0],409)
access=token('student2')
try:
    flag('activation_required',True)
    check('Pending activation rejects an existing access token',call('/borrowings',method='POST',body=request,access=access)[0],401)
    due='2099-12-01T04:00:00Z'
    direct={**request,'borrower_id':fixture['student2']['id'],'due_at':due,'physical_handover_confirmed':True}
    check('Direct issuance rejects pending activation',call('/borrowings/direct-checkout','staff','POST',direct)[0],403)
finally:flag('activation_required',False)
try:
    flag('is_active',False)
    check('Inactive Borrower rejects existing access token',call('/borrowings',method='POST',body=request,access=access)[0],403)
    check('Direct issuance rejects inactive Borrower',call('/borrowings/direct-checkout','staff','POST',direct)[0],403)
finally:flag('is_active',True)
inactive=next(x for x in fixture['equipment_records'] if x['status']=='INACTIVE')
zero=next(x for x in fixture['equipment_records'] if x['stock']['available']==0)
for label,x,q in [('Inactive equipment',inactive,1),('Zero-stock equipment',zero,1),('Insufficient available stock',fixture['equipment'],999)]:
    rejected=call('/borrowings','student','POST',{'items':[{'equipment_id':x['id'],'quantity':q}],'confirm':True})
    assert rejected[1]['error']['code']=='EQUIPMENT_NOT_AVAILABLE',(label,'wrong rejection cause')
    rejection_codes[label]=rejected[1]['error']['code']
    check(label+' rejected without stock mutation',rejected[0],409)
check('Borrower cannot perform Staff direct issuance',call('/borrowings/direct-checkout','student','POST',direct)[0],403)
check('Staff cannot create accounts',call('/borrowers','staff','POST',{})[0],403)
check('Unauthenticated borrowing read rejected',call('/borrowings')[0],401)
check('Strict request schema rejects photographic evidence fields',call('/borrowings','student','POST',{**request,'photo':'not-accepted'})[0],400)
stock=call('/equipment/'+fixture['equipment']['id'],'staff')[1]['data']['stock'];assert stock==before
checks.append('All rejected commands preserve original demo stock exactly')
current=call('/terms/current','admin')[1]['data']
assert current['version'].startswith('DEMO-'),'only synthetic demo versions may be revised'
next_version='DEMO-'+str(int(current['version'].split('-')[-1])+1)
new=call('/terms/versions','admin','POST',{'version':next_version,'title':'DEMONSTRATION TERMS — NOT OFFICIAL FSMO POLICY','body':'Updated synthetic demo terms for actual version-change and reacceptance testing only. Not approved institutional policy. No real email or borrowing is represented.','expected_current_version_id':current['id']})
check('Admin can publish a newer labeled DEMO version',new[0],201)
check('Existing acceptance cannot authorize a new request after version change',call('/borrowings','student','POST',request)[0],409)
check('Stale terms-version acceptance is rejected',call('/terms/'+current['id']+'/accept','student','POST',{})[0],409)
output={'result':'PASS','api':base,'origin':origin,'new_demo_version':next_version,'stock_rejection_codes':rejection_codes,'checks':checks,'stock':stock,'scope':'Real Phase7 APIs/PostgreSQL on isolated owner demo only; normal data untouched'}
p=demo.ROOT/'docs/project/verification/phase7-owner-demo/api-eligibility.json';p.parent.mkdir(parents=True,exist_ok=True);p.write_text(json.dumps(output,indent=2)+'\n')
print(str(len(checks))+' real API checks passed')
