-- Add previous_sprint_id column to tasks table
-- This tracks which sprint a task was in before being moved to backlog
ALTER TABLE tasks ADD COLUMN previous_sprint_id UUID;

-- Add index for faster lookups
CREATE INDEX idx_tasks_previous_sprint_id ON tasks(previous_sprint_id);

-- Add foreign key constraint
ALTER TABLE tasks ADD CONSTRAINT fk_tasks_previous_sprint 
    FOREIGN KEY (previous_sprint_id) REFERENCES sprints(id) ON DELETE SET NULL;
