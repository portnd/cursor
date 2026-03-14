-- Ensure submissions.task_id FK uses ON DELETE CASCADE so deleting a task
-- cascades to submissions (avoids SQLSTATE 23503 on delete).
-- Drop possible constraint names (PG default vs GORM/custom).
ALTER TABLE submissions DROP CONSTRAINT IF EXISTS fk_tasks_submissions;
ALTER TABLE submissions DROP CONSTRAINT IF EXISTS submissions_task_id_fkey;
ALTER TABLE submissions
    ADD CONSTRAINT submissions_task_id_fkey
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE;
