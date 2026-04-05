ALTER TABLE users DROP CONSTRAINT IF EXISTS check_user_role;
UPDATE users SET role = 'PM' WHERE role = 'PRODUCT_OWNER';
ALTER TABLE users ADD CONSTRAINT check_user_role CHECK (role IN ('CEO', 'MANAGER', 'PM', 'ENGINEER', 'SUPPORT'));
