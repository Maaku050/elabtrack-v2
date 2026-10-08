-- Existing identities and credential/session state are retained without inferred categories.
ALTER TABLE users ADD COLUMN activation_required BOOLEAN NOT NULL DEFAULT FALSE,
 ADD COLUMN activation_delivery_status TEXT NOT NULL DEFAULT 'NOT_REQUIRED'
 CHECK(activation_delivery_status IN ('NOT_REQUIRED','PENDING','UNCONFIGURED','ACCEPTED','FAILED','UNKNOWN'));
CREATE TABLE borrower_profiles (
 user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE RESTRICT,
 borrower_type TEXT NOT NULL CHECK(borrower_type IN ('STUDENT','FACULTY')),
 student_id TEXT UNIQUE,
 course TEXT NOT NULL DEFAULT '' CHECK(octet_length(course)<=128),
 contact_number TEXT NOT NULL DEFAULT '' CHECK(octet_length(contact_number)<=32),
 CHECK((borrower_type='STUDENT' AND student_id IS NOT NULL AND octet_length(student_id) BETWEEN 1 AND 64 AND btrim(student_id)<>'' AND student_id=btrim(student_id)) OR (borrower_type='FACULTY' AND student_id IS NULL))
);
CREATE FUNCTION guard_borrower_profile_role() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
 IF TG_TABLE_NAME='borrower_profiles' THEN
  IF NOT EXISTS(SELECT 1 FROM users WHERE id=NEW.user_id AND role='BORROWER') THEN RAISE EXCEPTION 'Borrower profile requires BORROWER role'; END IF;
 ELSIF NEW.role<>'BORROWER' AND EXISTS(SELECT 1 FROM borrower_profiles WHERE user_id=NEW.id) THEN RAISE EXCEPTION 'Classified borrower role cannot change';
 END IF;
 RETURN NEW;
END;
$$;
CREATE TRIGGER borrower_profile_role BEFORE INSERT OR UPDATE ON borrower_profiles FOR EACH ROW EXECUTE FUNCTION guard_borrower_profile_role();
CREATE TRIGGER classified_borrower_role BEFORE UPDATE OF role ON users FOR EACH ROW EXECUTE FUNCTION guard_borrower_profile_role();
CREATE TABLE account_activation_tokens (
 token_hash TEXT PRIMARY KEY CHECK(token_hash ~ '^[0-9a-f]{64}$'),
 user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 created_at TIMESTAMPTZ NOT NULL,
 expires_at TIMESTAMPTZ NOT NULL CHECK(expires_at>created_at),
 invalidated BOOLEAN NOT NULL DEFAULT FALSE,
 delivery_status TEXT NOT NULL DEFAULT 'PENDING' CHECK(delivery_status IN ('PENDING','UNCONFIGURED','ACCEPTED','FAILED','UNKNOWN')),
 provider_message_id TEXT NOT NULL DEFAULT '' CHECK(octet_length(provider_message_id)<=256)
);
CREATE UNIQUE INDEX account_activation_one_live ON account_activation_tokens(user_id) WHERE invalidated=FALSE;
CREATE TABLE account_operation_receipts (
 actor_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 operation TEXT NOT NULL CHECK(octet_length(operation) BETWEEN 1 AND 100),
 command_key UUID NOT NULL,
 payload_hash TEXT NOT NULL CHECK(payload_hash ~ '^[0-9a-f]{64}$'),
 result JSONB NOT NULL,
 created_at TIMESTAMPTZ NOT NULL DEFAULT statement_timestamp(),
 PRIMARY KEY(actor_id,operation,command_key)
);
CREATE TABLE account_audit_events (
 id UUID PRIMARY KEY,
 actor_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 account_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 action TEXT NOT NULL CHECK(octet_length(action) BETWEEN 1 AND 100),
 occurred_at TIMESTAMPTZ NOT NULL DEFAULT statement_timestamp()
);
CREATE INDEX account_audit_target_time ON account_audit_events(account_id,occurred_at DESC,id);
CREATE TABLE account_roster_batches (
 id UUID PRIMARY KEY,
 actor_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 operation TEXT NOT NULL CHECK(operation IN ('CREATE','DEACTIVATE')),
 preview JSONB NOT NULL,
 expires_at TIMESTAMPTZ NOT NULL,
 result JSONB,
 created_at TIMESTAMPTZ NOT NULL DEFAULT statement_timestamp()
);
CREATE INDEX account_batches_actor ON account_roster_batches(actor_id,created_at DESC,id);
CREATE FUNCTION reject_account_history_change() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN RAISE EXCEPTION 'Account audit and committed receipts are immutable'; END;
$$;
CREATE TRIGGER account_audit_immutable BEFORE UPDATE OR DELETE OR TRUNCATE ON account_audit_events FOR EACH STATEMENT EXECUTE FUNCTION reject_account_history_change();
CREATE TRIGGER account_receipts_immutable BEFORE UPDATE OR DELETE OR TRUNCATE ON account_operation_receipts FOR EACH STATEMENT EXECUTE FUNCTION reject_account_history_change();
CREATE FUNCTION guard_account_batch_confirmation() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
 IF OLD.result IS NOT NULL OR ROW(OLD.id,OLD.actor_id,OLD.operation,OLD.preview,OLD.expires_at,OLD.created_at) IS DISTINCT FROM ROW(NEW.id,NEW.actor_id,NEW.operation,NEW.preview,NEW.expires_at,NEW.created_at) OR NEW.result IS NULL THEN RAISE EXCEPTION 'Reviewed roster is immutable and confirms once'; END IF;
 RETURN NEW;
END;
$$;
CREATE TRIGGER account_batch_once BEFORE UPDATE ON account_roster_batches FOR EACH ROW EXECUTE FUNCTION guard_account_batch_confirmation();
REVOKE ALL ON FUNCTION guard_borrower_profile_role(),reject_account_history_change(),guard_account_batch_confirmation() FROM PUBLIC;
GRANT SELECT,INSERT ON borrower_profiles,account_activation_tokens,account_roster_batches,account_operation_receipts,account_audit_events TO elabtrack_runtime;
GRANT UPDATE(course,contact_number) ON borrower_profiles TO elabtrack_runtime;
GRANT UPDATE(invalidated,delivery_status,provider_message_id) ON account_activation_tokens TO elabtrack_runtime;
GRANT UPDATE(result) ON account_roster_batches TO elabtrack_runtime;
