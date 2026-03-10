DROP TABLE IF EXISTS team_transactions;

ALTER TABLE teams
    DROP COLUMN IF EXISTS treasury_balance,
    DROP COLUMN IF EXISTS bonus_percentage;
