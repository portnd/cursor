DROP INDEX IF EXISTS idx_tasks_is_komgrip;
ALTER TABLE tasks DROP COLUMN IF EXISTS is_komgrip;
