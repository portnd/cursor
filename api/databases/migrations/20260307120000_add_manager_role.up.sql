-- Add MANAGER role to check_user_role constraint
-- MANAGER has global project visibility (like CEO) but no admin rights
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_user_role;
ALTER TABLE users ADD CONSTRAINT check_user_role CHECK (role IN ('CEO', 'MANAGER', 'PM', 'DEV'));
