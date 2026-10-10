import importlib.util
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch
spec=importlib.util.spec_from_file_location('presentation_recovery',Path(__file__).parent/'recovery.py')
r=importlib.util.module_from_spec(spec);spec.loader.exec_module(r)

class RecoveryGuards(unittest.TestCase):
 def test_normal_target_refused_before_any_docker_command(self):
  with patch.object(r.m,'run') as run:
   with self.assertRaisesRegex(RuntimeError,'foreign recovery'):r.query('elabtrack_v2_postgres','elabtrack_v2','SELECT 1;')
   run.assert_not_called()
 def test_unseeded_source_refused_before_backup(self):
  with patch.object(r.m,'target',return_value=({},{})),patch.object(r.m,'sql',return_value='elabtrack-presentation-demo:false'),patch.object(r.m,'run') as run:
   with self.assertRaisesRegex(RuntimeError,'unmarked/unseeded'):r.main()
   run.assert_not_called()
 def test_live_application_refused_before_snapshot_or_backup(self):
  with tempfile.TemporaryDirectory() as directory:
   private=Path(directory);(private/'processes.json').write_text('{}')
   with patch.object(r.m,'PRIVATE',private),patch.object(r.m,'target',return_value=({},{})),patch.object(r.m,'sql',return_value='elabtrack-presentation-demo:true'),patch.object(r.m,'read_private',return_value={'api':{'pid':123}}),patch.object(r.os,'kill'),patch.object(r,'snapshot') as snapshot,patch.object(r.m,'run') as run:
    with self.assertRaisesRegex(RuntimeError,'Stop the presentation'):r.main()
    snapshot.assert_not_called();run.assert_not_called()

if __name__=='__main__':unittest.main()
