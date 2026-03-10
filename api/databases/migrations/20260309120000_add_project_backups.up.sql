CREATE TABLE IF NOT EXISTS project_backups (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    label       VARCHAR(255) NOT NULL DEFAULT '',
    payload     JSONB NOT NULL DEFAULT '{}',
    created_by  INTEGER,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_project_backups_project_id ON project_backups(project_id);
CREATE INDEX IF NOT EXISTS idx_project_backups_created_at ON project_backups(created_at DESC);
