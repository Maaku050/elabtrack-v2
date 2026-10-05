-- Rolling back also invalidates every session. Hashes cannot recover secrets.
-- Roll back application and schema together; never run the hash-only API here.
DELETE FROM refresh_tokens;
DROP INDEX idx_refresh_tokens_replaced_by;
DROP INDEX idx_refresh_tokens_terminal;
DROP INDEX idx_refresh_tokens_user_id;
ALTER TABLE refresh_tokens
    DROP CONSTRAINT refresh_tokens_token_hash_key,
    DROP CONSTRAINT refresh_tokens_hash_format,
    DROP CONSTRAINT refresh_tokens_lifetime,
    DROP CONSTRAINT refresh_tokens_replacement_revoked,
    DROP CONSTRAINT refresh_tokens_not_self_replaced,
    DROP COLUMN replaced_by,
    DROP COLUMN revoked_at;
ALTER TABLE refresh_tokens RENAME COLUMN token_hash TO token;
ALTER TABLE refresh_tokens
    ADD COLUMN revoked BOOLEAN NOT NULL DEFAULT FALSE,
    ADD CONSTRAINT refresh_tokens_token_key UNIQUE (token);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens (token);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id) WHERE revoked = FALSE;
CREATE INDEX idx_refresh_tokens_expired ON refresh_tokens (expires_at) WHERE revoked = FALSE;
