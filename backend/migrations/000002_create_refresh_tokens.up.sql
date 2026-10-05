-- 000002_create_refresh_tokens up

CREATE TABLE refresh_tokens (
    id          UUID PRIMARY KEY,
    token       TEXT UNIQUE NOT NULL,
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked     BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Fast lookups by token (used on every refresh) and by user (logout-all).
CREATE INDEX idx_refresh_tokens_token    ON refresh_tokens (token);
CREATE INDEX idx_refresh_tokens_user_id  ON refresh_tokens (user_id) WHERE revoked = FALSE;

-- Periodic cleanup job can use this partial index to find expired tokens.
CREATE INDEX idx_refresh_tokens_expired ON refresh_tokens (expires_at) WHERE revoked = FALSE;
