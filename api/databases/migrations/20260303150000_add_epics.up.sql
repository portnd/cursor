-- Create epics table (Dimension 1: Hierarchy - Epic level)
CREATE TABLE IF NOT EXISTS epics (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    status      VARCHAR(32) NOT NULL DEFAULT 'PLANNING', -- PLANNING, IN_PROGRESS, DONE
    color       VARCHAR(16) DEFAULT '#6366f1',           -- UI accent color (hex)
    sort_order  INT DEFAULT 0,
    start_date  TIMESTAMPTZ,
    end_date    TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_epics_project_id ON epics(project_id);

-- Add epic_id FK to tasks (Task links to its parent Epic)
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS epic_id UUID REFERENCES epics(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_tasks_epic_id ON tasks(epic_id);
