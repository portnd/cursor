DROP TABLE IF EXISTS project_transactions;

ALTER TABLE projects
    DROP COLUMN IF EXISTS capital_balance,
    DROP COLUMN IF EXISTS bonus_percentage;
