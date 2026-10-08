-- Publication is immediate and mandatory; no draft CMS or invented content.
CREATE TABLE terms_versions (
 id UUID PRIMARY KEY,
 version TEXT NOT NULL UNIQUE CHECK (version ~ '^[A-Za-z0-9][A-Za-z0-9._-]{0,63}$'),
 title TEXT NOT NULL CHECK (octet_length(title) BETWEEN 1 AND 200 AND btrim(title)<>''),
 body TEXT NOT NULL CHECK (octet_length(body) BETWEEN 1 AND 65536 AND btrim(body)<>''),
 content_hash TEXT NOT NULL CHECK (content_hash ~ '^[0-9a-f]{64}$'),
 published_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 created_at TIMESTAMPTZ NOT NULL DEFAULT statement_timestamp(),
 published_at TIMESTAMPTZ NOT NULL DEFAULT statement_timestamp()
);
-- A single locked pointer prevents ambiguous selectors and publish/accept races.
CREATE TABLE terms_publication (
 id SMALLINT PRIMARY KEY CHECK(id=1),
 current_version_id UUID REFERENCES terms_versions(id) ON DELETE RESTRICT
);
INSERT INTO terms_publication(id,current_version_id) VALUES(1,NULL);
CREATE TABLE terms_acceptances (
 id UUID PRIMARY KEY,
 user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 terms_version_id UUID NOT NULL REFERENCES terms_versions(id) ON DELETE RESTRICT,
 accepted_at TIMESTAMPTZ NOT NULL DEFAULT statement_timestamp(),
 UNIQUE(user_id,terms_version_id),
 UNIQUE(id,user_id)
);
CREATE INDEX idx_terms_acceptances_version ON terms_acceptances(terms_version_id);
CREATE INDEX idx_terms_versions_publisher ON terms_versions(published_by);

-- Even owner DML cannot silently edit/delete consequential published history.
CREATE FUNCTION reject_terms_history_change() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN RAISE EXCEPTION 'Published terms and acceptance history are immutable'; END;
$$;
CREATE TRIGGER terms_versions_immutable BEFORE UPDATE OR DELETE OR TRUNCATE ON terms_versions
 FOR EACH STATEMENT EXECUTE FUNCTION reject_terms_history_change();
CREATE TRIGGER terms_acceptances_immutable BEFORE UPDATE OR DELETE OR TRUNCATE ON terms_acceptances
 FOR EACH STATEMENT EXECUTE FUNCTION reject_terms_history_change();
REVOKE ALL ON FUNCTION reject_terms_history_change() FROM PUBLIC;
-- Named grants only. Runtime cannot alter history, create schema or track migrations.
GRANT SELECT,INSERT ON terms_versions,terms_acceptances TO elabtrack_runtime;
GRANT SELECT ON terms_publication TO elabtrack_runtime;
GRANT UPDATE(current_version_id) ON terms_publication TO elabtrack_runtime;
