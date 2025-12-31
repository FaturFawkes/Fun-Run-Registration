-- Seed admin user for development
-- Email: admin@tautaurun.com
-- Password: Admin123!
-- 
-- IMPORTANT: This is for development/testing only
-- Change credentials in production!

INSERT INTO admins (email, password_hash, created_at)
VALUES (
    'admin@tautaurun.com',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYlK4Qr1WZK',
    CURRENT_TIMESTAMP
)
ON CONFLICT (email) DO NOTHING;

-- Verify admin was created
SELECT 'Admin user created successfully:' AS message, email, created_at 
FROM admins 
WHERE email = 'admin@tautaurun.com';
