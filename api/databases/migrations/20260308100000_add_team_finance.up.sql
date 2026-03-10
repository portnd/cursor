ALTER TABLE teams
    ADD COLUMN IF NOT EXISTS treasury_balance DECIMAL(15,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS bonus_percentage DECIMAL(5,2)  NOT NULL DEFAULT 0;

CREATE TABLE IF NOT EXISTS team_transactions (
    id           BIGSERIAL PRIMARY KEY,
    team_id      BIGINT       NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    type         VARCHAR(20)  NOT NULL,
    amount       DECIMAL(15,2) NOT NULL,
    reference    TEXT         NOT NULL DEFAULT '',
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_team_transaction_type CHECK (type IN ('INJECTION', 'BURN', 'BONUS_PAYOUT', 'ADJUSTMENT'))
);

CREATE INDEX IF NOT EXISTS idx_team_transactions_team_id ON team_transactions(team_id);
CREATE INDEX IF NOT EXISTS idx_team_transactions_created_at ON team_transactions(created_at DESC);
