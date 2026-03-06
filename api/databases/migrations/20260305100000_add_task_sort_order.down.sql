DROP INDEX IF EXISTS idx_tasks_epic_sort;
ALTER TABLE tasks DROP COLUMN IF EXISTS sort_order;
