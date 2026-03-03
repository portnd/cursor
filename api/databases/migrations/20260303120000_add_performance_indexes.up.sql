-- Add indexes to fix SLOW SQL: tasks by project_id, projects by code
-- Queries: GetTasksByProjectID(project_id), GetProjectByCode(code)

CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks(project_id);
CREATE INDEX IF NOT EXISTS idx_tasks_project_created_at ON tasks(project_id, created_at DESC);

-- projects.code lookup (GORM may already create unique index; this ensures fast lookup)
CREATE INDEX IF NOT EXISTS idx_projects_code_lookup ON projects(code);
