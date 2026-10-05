-- Development seed data. Idempotent: safe to run multiple times.
-- Never run in production.

-- Seed an admin user. Password: "password123" (bcrypt cost 10).
INSERT INTO users (id, email, name, password, role, is_active, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'admin@example.com',
    'Development Admin',
    '$2a$10$H01mdVWfETSpfAeUHxitQ.OCEGchTn1GIH9DSUoaL8VAQAB8pYnW2',
    'admin',
    TRUE,
    NOW(),
    NOW()
)
ON CONFLICT (email) DO NOTHING;

-- Seed a regular user. Password: "password123".
INSERT INTO users (id, email, name, password, role, is_active, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000002',
    'user@example.com',
    'Development User',
    '$2a$10$H01mdVWfETSpfAeUHxitQ.OCEGchTn1GIH9DSUoaL8VAQAB8pYnW2',
    'user',
    TRUE,
    NOW(),
    NOW()
)
ON CONFLICT (email) DO NOTHING;
