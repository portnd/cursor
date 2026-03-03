DROP INDEX IF EXISTS idx_tasks_epic_id;
ALTER TABLE tasks DROP COLUMN IF EXISTS epic_id;
DROP INDEX IF EXISTS idx_epics_project_id;
DROP TABLE IF EXISTS epics;
