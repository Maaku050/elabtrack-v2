-- V2 has no evidenced deployed production sessions. Explicitly invalidate all
-- existing raw-token sessions; users must log in again. No bearer backfill.
DELETE FROM refresh_tokens;

DROP INDEX idx_refresh_tokens_token;
DROP INDEX idx_refresh_tokens_user_id;
DROP INDEX idx_refresh_tokens_expired;
ALTER TABLE refresh_tokens DROP CONSTRAINT refresh_tokens_token_key;
ALTER TABLE refresh_tokens RENAME COLUMN token TO token_hash;
ALTER TABLE refresh_tokens DROP COLUMN revoked;
ALTER TABLE refresh_tokens
    ADD COLUMN revoked_at TIMESTAMPTZ,
    ADD COLUMN replaced_by UUID REFERENCES refresh_tokens(id) ON DELETE SET NULL,
    ADD CONSTRAINT refresh_tokens_token_hash_key UNIQUE (token_hash),
    ADD CONSTRAINT refresh_tokens_hash_format CHECK (token_hash ~ '^[0-9a-f]{64}$'),
    ADD CONSTRAINT refresh_tokens_lifetime CHECK (expires_at > created_at),
    ADD CONSTRAINT refresh_tokens_replacement_revoked CHECK (replaced_by IS NULL OR revoked_at IS NOT NULL),
    ADD CONSTRAINT refresh_tokens_not_self_replaced CHECK (replaced_by IS NULL OR replaced_by <> id);

-- UNIQUE already provides the indexed credential lookup; no duplicate index.
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id) WHERE revoked_at IS NULL;
CREATE INDEX idx_refresh_tokens_terminal ON refresh_tokens
    (LEAST(expires_at, COALESCE(revoked_at, expires_at)), id);
CREATE INDEX idx_refresh_tokens_replaced_by ON refresh_tokens (replaced_by) WHERE replaced_by IS NOT NULL;
