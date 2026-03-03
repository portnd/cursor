-- Ensure tasks table has all columns required by the Task domain model (CreateTask / GORM).
-- Safe to run: every column uses IF NOT EXISTS.

-- Allow ai_estimated_minutes = 0 (CreateTask inserts 0 by default; old constraint was > 0 only)
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS check_ai_estimated_minutes;
ALTER TABLE tasks ADD CONSTRAINT check_ai_estimated_minutes CHECK (ai_estimated_minutes IS NULL OR ai_estimated_minutes >= 0);

-- Task code (e.g. mims-hdmap-001); unique per task
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS code VARCHAR(64) DEFAULT '';

-- Project / WBS / Gantt
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS project_id UUID REFERENCES projects(id) ON DELETE SET NULL;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS parent_id UUID REFERENCES tasks(id) ON DELETE SET NULL;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS start_date TIMESTAMP WITH TIME ZONE;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS end_date TIMESTAMP WITH TIME ZONE;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS progress INT NOT NULL DEFAULT 0;

-- Priority and estimation
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS priority VARCHAR(20) NOT NULL DEFAULT 'MEDIUM';
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS story_points INT NOT NULL DEFAULT 0;

-- Sprint and milestone (FKs added only if tables exist; references created in other migrations)
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS sprint_id UUID;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS milestone_id UUID;

-- Unique index on code (allow empty for legacy; partial unique where code != '')
CREATE UNIQUE INDEX IF NOT EXISTS idx_tasks_code_unique ON tasks(code) WHERE code != '';

-- Indexes for filters
CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks(project_id);
CREATE INDEX IF NOT EXISTS idx_tasks_parent_id ON tasks(parent_id);
CREATE INDEX IF NOT EXISTS idx_tasks_sprint_id ON tasks(sprint_id);
CREATE INDEX IF NOT EXISTS idx_tasks_milestone_id ON tasks(milestone_id);
