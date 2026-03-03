-- Revert task columns added for full Task model (optional; may break app if data exists)
DROP INDEX IF EXISTS idx_tasks_milestone_id;
DROP INDEX IF EXISTS idx_tasks_sprint_id;
DROP INDEX IF EXISTS idx_tasks_parent_id;
DROP INDEX IF EXISTS idx_tasks_project_id;
DROP INDEX IF EXISTS idx_tasks_code_unique;

ALTER TABLE tasks DROP COLUMN IF EXISTS milestone_id;
ALTER TABLE tasks DROP COLUMN IF EXISTS sprint_id;
ALTER TABLE tasks DROP COLUMN IF EXISTS story_points;
ALTER TABLE tasks DROP COLUMN IF EXISTS priority;
ALTER TABLE tasks DROP COLUMN IF EXISTS end_date;
ALTER TABLE tasks DROP COLUMN IF EXISTS start_date;
ALTER TABLE tasks DROP COLUMN IF EXISTS parent_id;
ALTER TABLE tasks DROP COLUMN IF EXISTS project_id;
ALTER TABLE tasks DROP COLUMN IF EXISTS progress;
ALTER TABLE tasks DROP COLUMN IF EXISTS code;
