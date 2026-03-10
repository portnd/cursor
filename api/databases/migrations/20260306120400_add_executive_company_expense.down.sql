ALTER TABLE company_cost_configs
  DROP COLUMN IF EXISTS executive_expense,
  DROP COLUMN IF EXISTS company_expense;
