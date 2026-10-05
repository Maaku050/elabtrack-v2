-- 000001_create_users up

-- citext must be enabled before the table that uses it.
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users (
    id          UUID PRIMARY KEY,
    email       CITEXT UNIQUE NOT NULL,
    name        TEXT NOT NULL,
    password    TEXT NOT NULL,
    role        TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('user','admin')),
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes: email is already unique (above). Add a created_at index for
-- default ordering on list endpoints.
CREATE INDEX idx_users_created_at ON users (created_at DESC);
