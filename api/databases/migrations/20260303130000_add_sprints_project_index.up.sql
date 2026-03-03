-- Speed up GetSprintsByProjectID(projectID)
CREATE INDEX IF NOT EXISTS idx_sprints_project_id ON sprints(project_id);
