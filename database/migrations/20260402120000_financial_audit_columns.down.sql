SET statement_timeout = 0;

--bun:split

ALTER TABLE budgets
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS updated_at;

--bun:split

ALTER TABLE categories
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS updated_at;

--bun:split

ALTER TABLE transactions
    DROP COLUMN IF EXISTS deleted_at;

--bun:split

ALTER TABLE wallets
    DROP COLUMN IF EXISTS deleted_at;
