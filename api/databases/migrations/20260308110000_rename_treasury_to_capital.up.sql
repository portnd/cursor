-- Rename treasury_balance to capital_balance on teams (idempotent: skip if already renamed)
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'teams' AND column_name = 'treasury_balance'
  ) AND NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'teams' AND column_name = 'capital_balance'
  ) THEN
    ALTER TABLE teams RENAME COLUMN treasury_balance TO capital_balance;
  END IF;
END $$;
