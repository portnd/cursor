-- Add capital tracking to projects (independent of Team capital pool)
ALTER TABLE projects
    ADD COLUMN IF NOT EXISTS capital_balance  DECIMAL(15,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS bonus_percentage DECIMAL(5,2)  NOT NULL DEFAULT 0;

-- Project-scoped capital transaction log
CREATE TABLE IF NOT EXISTS project_transactions (
    id           BIGSERIAL PRIMARY KEY,
    project_id   UUID         NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    type         VARCHAR(20)  NOT NULL,
    amount       DECIMAL(15,2) NOT NULL,
    reference    TEXT         NOT NULL DEFAULT '',
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_project_transaction_type CHECK (type IN ('INJECTION','BURN','BONUS_PAYOUT','ADJUSTMENT'))
);

CREATE INDEX IF NOT EXISTS idx_project_transactions_project_id ON project_transactions(project_id);
CREATE INDEX IF NOT EXISTS idx_project_transactions_created_at ON project_transactions(created_at DESC);
