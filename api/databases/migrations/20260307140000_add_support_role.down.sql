-- Revert SUPPORT role: restore constraint without SUPPORT
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_user_role;
ALTER TABLE users ADD CONSTRAINT check_user_role CHECK (role IN ('CEO', 'MANAGER', 'PM', 'DEV'));
-- Reassign any SUPPORT users to DEV
UPDATE users SET role = 'DEV' WHERE role = 'SUPPORT';
