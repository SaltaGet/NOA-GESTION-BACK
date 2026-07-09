ALTER TABLE income_sales ADD COLUMN IF NOT EXISTS delivered boolean DEFAULT false NOT NULL;
ALTER TABLE income_sales ADD COLUMN IF NOT EXISTS is_budget boolean DEFAULT false NOT NULL;
