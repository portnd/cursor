-- Add default monthly expense fields to company cost config (for quotation defaults)
ALTER TABLE company_cost_configs
  ADD COLUMN IF NOT EXISTS executive_expense DECIMAL(15,2) NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS company_expense DECIMAL(15,2) NOT NULL DEFAULT 0;
