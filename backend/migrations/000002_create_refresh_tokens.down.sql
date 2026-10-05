-- 000002_create_refresh_tokens down

DROP INDEX IF EXISTS idx_refresh_tokens_expired;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP INDEX IF EXISTS idx_refresh_tokens_token;
DROP TABLE IF EXISTS refresh_tokens;
