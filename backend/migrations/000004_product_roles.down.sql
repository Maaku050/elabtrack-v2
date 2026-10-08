-- The old taxonomy cannot express STAFF. Refuse lossy rollback atomically;
-- an operator must explicitly reconcile those accounts before rolling back.
DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM users WHERE role = 'STAFF') THEN
        RAISE EXCEPTION 'Product STAFF accounts require explicit reconciliation before rollback';
    END IF;
END $$;
ALTER TABLE users DROP CONSTRAINT users_role_check;
UPDATE users SET role = CASE role WHEN 'BORROWER' THEN 'user' WHEN 'ADMIN' THEN 'admin' ELSE role END;
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'user';
ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (role IN ('user', 'admin'));
