-- Add task_type column to tasks table for Task Typology system (FEATURE, TASK, BUG)
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS task_type VARCHAR(20) NOT NULL DEFAULT 'TASK';

-- Constraint to enforce valid enum values
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS check_task_type;
ALTER TABLE tasks ADD CONSTRAINT check_task_type CHECK (task_type IN ('FEATURE', 'TASK', 'BUG'));

CREATE INDEX IF NOT EXISTS idx_tasks_task_type ON tasks(task_type);
