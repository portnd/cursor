-- Allow CHIEF_ENGINEER role (same capabilities as ENGINEER; stored value CHIEF_ENGINEER)
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_user_role;
ALTER TABLE users ADD CONSTRAINT check_user_role CHECK (role IN ('CEO', 'MANAGER', 'PRODUCT_OWNER', 'ENGINEER', 'CHIEF_ENGINEER', 'SUPPORT'));
