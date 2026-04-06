-- Indexes for project details/task paging queries

-- Tasks listing for project details and paging (created_at desc + id desc ordering)
CREATE INDEX IF NOT EXISTS idx_tasks_project_created_at_id_desc
    ON tasks(project_id, created_at DESC, id DESC);

-- Sprints list by project (sort_order + created_at)
CREATE INDEX IF NOT EXISTS idx_sprints_project_sort_created
    ON sprints(project_id, sort_order ASC, created_at ASC);

-- Milestones list by project (due_date)
CREATE INDEX IF NOT EXISTS idx_milestones_project_due_date
    ON milestones(project_id, due_date ASC);

-- Epics list by project (sort_order + start_date)
CREATE INDEX IF NOT EXISTS idx_epics_project_sort_start
    ON epics(project_id, sort_order ASC, start_date ASC NULLS LAST, created_at ASC);
