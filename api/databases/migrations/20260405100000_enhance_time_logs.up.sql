-- Enhance time_logs table with work_type, logged_date, is_timer_session
-- This allows backfilling logs for past dates (≤7 days) and categorizing work types

ALTER TABLE time_logs
    ADD COLUMN IF NOT EXISTS work_type VARCHAR(20) DEFAULT 'DEV'
        CHECK (work_type IN ('DEV','REVIEW','TESTING','MEETING','RESEARCH','OTHER')),
    ADD COLUMN IF NOT EXISTS logged_date DATE DEFAULT CURRENT_DATE,
    ADD COLUMN IF NOT EXISTS is_timer_session BOOLEAN DEFAULT FALSE;

-- Index for daily user reports (Discipline page, personal summary)
CREATE INDEX IF NOT EXISTS idx_time_logs_user_date
    ON time_logs(user_id, logged_date DESC);

-- Index for analytics by work type
CREATE INDEX IF NOT EXISTS idx_time_logs_task_date
    ON time_logs(task_id, logged_date DESC);

-- Backfill existing rows: set logged_date from logged_at
UPDATE time_logs SET logged_date = logged_at::DATE WHERE logged_date = CURRENT_DATE AND logged_at < NOW() - INTERVAL '1 minute';
