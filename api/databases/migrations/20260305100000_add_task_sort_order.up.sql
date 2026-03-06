-- Backlog drag-and-drop: order of tasks within an epic (or unassigned)
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS sort_order INT NOT NULL DEFAULT 0;
CREATE INDEX IF NOT EXISTS idx_tasks_epic_sort ON tasks(epic_id, sort_order);
