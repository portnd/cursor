-- PM owners per project when teams/squads feature is disabled (CEO assigns one or more PMs)
CREATE TABLE IF NOT EXISTS project_pm_assignments (
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (project_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_project_pm_assignments_user_id ON project_pm_assignments(user_id);
