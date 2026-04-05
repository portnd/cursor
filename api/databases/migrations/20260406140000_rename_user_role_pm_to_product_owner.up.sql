-- User role PM → PRODUCT_OWNER (display: Product Owner)
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_user_role;
UPDATE users SET role = 'PRODUCT_OWNER' WHERE role = 'PM';
ALTER TABLE users ADD CONSTRAINT check_user_role CHECK (role IN ('CEO', 'MANAGER', 'PRODUCT_OWNER', 'ENGINEER', 'SUPPORT'));
