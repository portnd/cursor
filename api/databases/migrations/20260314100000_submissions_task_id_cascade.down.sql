-- Restore FK without CASCADE (rollback)
ALTER TABLE submissions DROP CONSTRAINT IF EXISTS submissions_task_id_fkey;
ALTER TABLE submissions
    ADD CONSTRAINT submissions_task_id_fkey
    FOREIGN KEY (task_id) REFERENCES tasks(id);
