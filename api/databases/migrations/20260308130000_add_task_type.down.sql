-- Revert task_type column
DROP INDEX IF EXISTS idx_tasks_task_type;
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS check_task_type;
ALTER TABLE tasks DROP COLUMN IF EXISTS task_type;
