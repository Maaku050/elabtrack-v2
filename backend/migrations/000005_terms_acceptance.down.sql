-- Refuse destruction of ANY published/accepted history, even if now superseded.
DO $$ BEGIN
 IF EXISTS(SELECT 1 FROM terms_versions) OR EXISTS(SELECT 1 FROM terms_acceptances) THEN
  RAISE EXCEPTION 'Terms history requires an explicitly reviewed preservation plan before rollback';
 END IF;
END $$;
DROP TABLE terms_acceptances;
DROP TABLE terms_publication;
DROP TABLE terms_versions;
DROP FUNCTION reject_terms_history_change();
