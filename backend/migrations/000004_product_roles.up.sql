-- Atomic role reconciliation; identity, password, status and sessions retained.
ALTER TABLE users DROP CONSTRAINT users_role_check;
UPDATE users SET role = CASE role WHEN 'user' THEN 'BORROWER' WHEN 'admin' THEN 'ADMIN' ELSE role END;
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'BORROWER';
ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (role IN ('BORROWER', 'STAFF', 'ADMIN'));
