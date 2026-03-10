-- Revert MANAGER role: restore original check constraint (CEO, PM, DEV only)
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_user_role;
ALTER TABLE users ADD CONSTRAINT check_user_role CHECK (role IN ('CEO', 'PM', 'DEV'));
-- Update any MANAGER users back to DEV
UPDATE users SET role = 'DEV' WHERE role = 'MANAGER';
