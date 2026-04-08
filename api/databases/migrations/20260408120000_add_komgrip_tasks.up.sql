ALTER TABLE tasks ADD COLUMN IF NOT EXISTS is_komgrip BOOLEAN NOT NULL DEFAULT FALSE;
CREATE INDEX IF NOT EXISTS idx_tasks_is_komgrip ON tasks (is_komgrip) WHERE is_komgrip = TRUE;
