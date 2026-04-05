-- Rename user role DEV → ENGINEER (time_log work_type 'DEV' is unchanged)
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_user_role;
UPDATE users SET role = 'ENGINEER' WHERE role = 'DEV';
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'ENGINEER';
ALTER TABLE users ADD CONSTRAINT check_user_role CHECK (role IN ('CEO', 'MANAGER', 'PM', 'ENGINEER', 'SUPPORT'));
