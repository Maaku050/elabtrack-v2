"""Target guard regression tests; no real database writes."""
import importlib.util
import unittest
import tempfile
import subprocess
from pathlib import Path
from unittest.mock import patch
s=importlib.util.spec_from_file_location('presentation_tool',Path(__file__).resolve().parents[1]/'presentation-demo.py')
m=importlib.util.module_from_spec(s);s.loader.exec_module(m)
class Guards(unittest.TestCase):
 def test_wrong_database_is_rejected_before_docker(self):
  with patch.object(m,'read_private',return_value={'DB_NAME':'elabtrack_v2','DB_PORT':'5434','DB_HOST':'127.0.0.1','APP_ENV':'development'}),patch.object(m,'run') as run:
   with self.assertRaisesRegex(RuntimeError,'Wrong demo target'):m.target()
   run.assert_not_called()
 def test_reset_requires_exact_confirmation_and_never_stops_other_target(self):
  with patch.object(m,'target'),patch.object(m,'stop') as stop,patch.object(m,'sql') as sql:
   with self.assertRaisesRegex(RuntimeError,'Reset requires'):m.reset('elabtrack_v2')
   stop.assert_not_called();sql.assert_not_called()
 def test_reset_refuses_unverified_database_marker(self):
  with patch.object(m,'target'),patch.object(m,'stop') as stop,patch.object(m,'sql',return_value='wrong:true'):
   with self.assertRaisesRegex(RuntimeError,'identity not verified'):m.reset(m.DB)
   stop.assert_not_called()
 def test_fresh_working_copy_refuses_existing_container_before_snapshot(self):
  with tempfile.TemporaryDirectory() as directory,patch.object(m,'PRIVATE',Path(directory)),patch.object(m.subprocess,'run',return_value=subprocess.CompletedProcess([],0)),patch.object(m,'snapshot_baseline') as snapshot:
   with self.assertRaisesRegex(RuntimeError,'already exists'):m.prepare()
   snapshot.assert_not_called()
 def test_occupied_application_port_refused_before_process_start(self):
  with tempfile.TemporaryDirectory() as directory,patch.object(m,'PRIVATE',Path(directory)),patch.object(m,'target'),patch.object(m,'run'),patch.object(m.socket,'socket') as socket,patch.object(m,'environment') as environment:
   socket.return_value.__enter__.return_value.bind.side_effect=OSError('occupied')
   with self.assertRaisesRegex(RuntimeError,'occupied'):m.start()
   environment.assert_not_called()
 def test_credentials_cannot_enter_noninteractive_logs(self):
  with patch.object(m,'target'),patch.object(m.sys.stdin,'isatty',return_value=False),patch.object(m,'read_private') as read:
   with self.assertRaisesRegex(RuntimeError,'interactive terminal'):m.credentials()
   read.assert_not_called()
 def test_runtime_environment_excludes_owner_and_mail_secrets(self):
  v={'DB_NAME':m.DB,'MIGRATION_DB_PASSWORD':'temporary-test-owner','BOOTSTRAP_DB_PASSWORD':'temporary-test-bootstrap'}
  with patch.object(m,'target',return_value=(v,{})),patch.dict(m.os.environ,{'DATABASE_URL':'wrong','BREVO_API_KEY':'temporary-test-key'}):
   e=m.environment(False);self.assertNotIn('MIGRATION_DB_PASSWORD',e);self.assertNotIn('BOOTSTRAP_DB_PASSWORD',e);self.assertNotIn('DATABASE_URL',e);self.assertEqual(e['BREVO_API_KEY'],'')
if __name__=='__main__':unittest.main()
