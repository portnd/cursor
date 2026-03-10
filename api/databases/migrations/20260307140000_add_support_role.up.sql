-- Add SUPPORT role to check_user_role constraint
-- SUPPORT: back-office / non-technical staff (HR, accounting, etc.)
-- Their salaries are treated as company overhead, not billable dev cost.
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_user_role;
ALTER TABLE users ADD CONSTRAINT check_user_role CHECK (role IN ('CEO', 'MANAGER', 'PM', 'DEV', 'SUPPORT'));
