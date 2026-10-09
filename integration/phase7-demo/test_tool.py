"""Guard tests without database mutation; real DB/browser evidence is separate."""
import importlib.util
from pathlib import Path
import tempfile
import json
import os
import unittest
from unittest.mock import patch
spec=importlib.util.spec_from_file_location('phase7_demo_tool',Path(__file__).resolve().parents[1]/'phase7-demo.py')
demo=importlib.util.module_from_spec(spec);spec.loader.exec_module(demo)

class DemoGuards(unittest.TestCase):
    def test_private_credentials_permissions(self):
        with tempfile.TemporaryDirectory() as tmp:
            p=Path(tmp)/'private.json';demo.protected(p,'{}');self.assertEqual(p.stat().st_mode&0o777,0o600)
            p.chmod(0o644)
            with self.assertRaisesRegex(RuntimeError,'permissions'):demo.read_private(p)

    def test_normal_database_is_rejected_before_docker(self):
        with tempfile.TemporaryDirectory() as tmp,patch.object(demo,'PRIVATE',Path(tmp)),patch.object(demo,'run') as run:
            demo.protected(Path(tmp)/'config.json',json.dumps({'DB_NAME':'elabtrack_v2','DB_PORT':'5434','DB_HOST':'127.0.0.1','APP_ENV':'development'}))
            with self.assertRaisesRegex(RuntimeError,'Wrong demo target'):demo.target()
            run.assert_not_called()

    def test_normal_volume_is_rejected_before_docker(self):
        with tempfile.TemporaryDirectory() as tmp,patch.object(demo,'PRIVATE',Path(tmp)),patch.object(demo,'run') as run:
            demo.protected(Path(tmp)/'config.json',json.dumps({'DB_NAME':demo.DB,'DB_PORT':demo.PORT,'DB_HOST':'127.0.0.1','APP_ENV':'development'}))
            demo.protected(Path(tmp)/'target.json',json.dumps({'container':demo.CONTAINER,'database':demo.DB,'volume':'elabtrack_v2_pgdata','disposable':True}))
            with self.assertRaisesRegex(RuntimeError,'disposable identity'):demo.target()
            run.assert_not_called()

    def test_reset_requires_exact_confirmation(self):
        with patch.object(demo,'target'),patch.object(demo,'sql') as sql,patch.object(demo,'stop') as stop:
            for value in [None,'yes','elabtrack_v2','production']:
                with self.assertRaisesRegex(RuntimeError,'Reset requires'):demo.reset(value)
            sql.assert_not_called();stop.assert_not_called()

    def test_mismatched_marker_prevents_reset(self):
        with patch.object(demo,'target'),patch.object(demo,'sql',return_value='production:true'),patch.object(demo,'stop') as stop:
            with self.assertRaisesRegex(RuntimeError,'identity not verified'):demo.reset(demo.DB)
            stop.assert_not_called()

    def test_credentials_reject_noninteractive_output(self):
        with patch.object(demo,'target'),patch.object(demo.sys.stdin,'isatty',return_value=False):
            with self.assertRaisesRegex(RuntimeError,'interactive terminal'):demo.credentials()

if __name__=='__main__':unittest.main()
