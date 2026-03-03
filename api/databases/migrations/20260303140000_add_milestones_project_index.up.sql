-- Speed up GetMilestonesByProjectID(projectID)
CREATE INDEX IF NOT EXISTS idx_milestones_project_id ON milestones(project_id);
