-- Remove previous_sprint_id column from tasks table
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS fk_tasks_previous_sprint;
DROP INDEX IF EXISTS idx_tasks_previous_sprint_id;
ALTER TABLE tasks DROP COLUMN previous_sprint_id;
