#!/usr/bin/env python3
"""Read-only presentation reconciliation; never prints credentials or private rows."""
import importlib.util,json,hashlib
from pathlib import Path
s=importlib.util.spec_from_file_location('demo',Path(__file__).resolve().parents[1]/'presentation-demo.py');m=importlib.util.module_from_spec(s);s.loader.exec_module(m)
queries={
 'counts':"SELECT json_build_array((SELECT count(*) FROM users),(SELECT count(*) FROM borrower_profiles WHERE borrower_type='STUDENT'),(SELECT count(*) FROM borrower_profiles WHERE borrower_type='FACULTY'),(SELECT count(*) FROM users WHERE role='STAFF'),(SELECT count(*) FROM users WHERE role='ADMIN'),(SELECT count(*) FROM equipment),(SELECT count(*) FROM equipment_images),(SELECT count(*) FROM profile_images WHERE image_id IS NOT NULL));",
 'stock':"SELECT count(*) FROM equipment WHERE total_tracked<>available+reserved+checked_out+damaged_held OR LEAST(available,reserved,checked_out,damaged_held)<0;",
 'ledger':"SELECT count(*) FROM equipment e LEFT JOIN LATERAL(SELECT * FROM inventory_movements m WHERE m.equipment_id=e.id ORDER BY sequence DESC LIMIT 1)m ON true WHERE e.stock_sequence>0 AND (m.sequence IS NULL OR ROW(e.available,e.reserved,e.checked_out,e.damaged_held,e.total_tracked,e.stock_sequence) IS DISTINCT FROM ROW(m.after_available,m.after_reserved,m.after_checked_out,m.after_damaged_held,m.after_total,m.sequence));",
 'movement_sums':"SELECT count(*) FROM equipment e LEFT JOIN LATERAL(SELECT COALESCE(sum(delta_available),0) a,COALESCE(sum(delta_reserved),0) r,COALESCE(sum(delta_checked_out),0) c,COALESCE(sum(delta_damaged_held),0) d,COALESCE(sum(delta_total),0) t FROM inventory_movements WHERE equipment_id=e.id)m ON true WHERE ROW(e.available,e.reserved,e.checked_out,e.damaged_held,e.total_tracked) IS DISTINCT FROM ROW(m.a,m.r,m.c,m.d,m.t);",
 'chronology':"SELECT count(*) FROM borrowings b JOIN terms_acceptances a ON a.id=b.acceptance_id JOIN terms_versions t ON t.id=a.terms_version_id WHERE a.accepted_at>b.created_at OR t.published_at>a.accepted_at OR b.checked_out_at<b.created_at OR b.due_at<=b.checked_out_at OR b.completed_at<b.checked_out_at;",
 'ledger_chronology':"SELECT count(*) FROM(SELECT created_at,lag(created_at) OVER(PARTITION BY equipment_id ORDER BY sequence) previous FROM inventory_movements)t WHERE created_at<previous;",
 'synthetic':"SELECT count(*) FROM users WHERE email NOT LIKE '%.invalid' OR name NOT LIKE 'DEMO %';",
 'replacement_overflow':"SELECT count(*) FROM replacement_obligations o WHERE required<(SELECT COALESCE(sum(quantity),0) FROM replacement_acceptances a WHERE a.obligation_id=o.id);",
 'notification_uniqueness':"SELECT count(*) FROM(SELECT event_key FROM notification_dispatch GROUP BY event_key HAVING count(*)>1)t;",
}
results={}
for name,q in queries.items():
 raw=m.sql('BEGIN READ ONLY;'+q+'COMMIT;');v=next(x for x in raw.splitlines() if x not in ('BEGIN','COMMIT',''))
 if name=='counts':assert json.loads(v)==[28,20,5,2,1,40,40,28],v
 else:assert int(v)==0,(name,v)
 results[name]='PASS'
assets=Path(__file__).parent/'assets';avatars=list((assets/'avatars').glob('*.png'));assert len({hashlib.sha256(p.read_bytes()).hexdigest() for p in avatars})==28
results['distinct_fictional_avatars']='PASS'
manifest=json.loads((assets/'MANIFEST.json').read_text());assert len(manifest['files'])==68
for entry in manifest['files']:
 file=assets/entry['file'];assert file.resolve().is_relative_to(assets.resolve())
 data=file.read_bytes();assert len(data)==entry['bytes'] and hashlib.sha256(data).hexdigest()==entry['sha256']
 assert hashlib.sha256(file.with_suffix('.svg').read_bytes()).hexdigest()==entry['source_sha256']
results['bundled_asset_provenance_hashes']='PASS';print(json.dumps(results,indent=2))
