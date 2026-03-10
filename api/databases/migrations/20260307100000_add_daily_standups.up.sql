CREATE TABLE IF NOT EXISTS daily_standups (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         BIGINT NOT NULL,
    date            DATE NOT NULL,
    yesterday_summary TEXT NOT NULL,
    today_task_ids  TEXT[] NOT NULL DEFAULT '{}',
    blocker         TEXT NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_daily_standups_user_date UNIQUE (user_id, date),
    CONSTRAINT fk_daily_standups_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_daily_standups_date ON daily_standups (date);
