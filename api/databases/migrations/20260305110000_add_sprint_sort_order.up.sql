-- Sprint drag-and-drop: display order
ALTER TABLE sprints ADD COLUMN IF NOT EXISTS sort_order INT NOT NULL DEFAULT 0;
