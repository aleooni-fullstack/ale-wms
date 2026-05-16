ALTER TABLE stock_balances
ADD COLUMN reserved_quantity NUMERIC(15, 4) NOT NULL DEFAULT 0;