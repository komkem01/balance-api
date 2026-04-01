SET statement_timeout = 0;

--bun:split

ALTER TABLE wallets
    ADD COLUMN IF NOT EXISTS deleted_at timestamptz;

--bun:split

ALTER TABLE transactions
    ADD COLUMN IF NOT EXISTS deleted_at timestamptz;

--bun:split

ALTER TABLE categories
    ADD COLUMN IF NOT EXISTS updated_at timestamptz NOT NULL DEFAULT now(),
    ADD COLUMN IF NOT EXISTS deleted_at timestamptz;

--bun:split

ALTER TABLE budgets
    ADD COLUMN IF NOT EXISTS updated_at timestamptz NOT NULL DEFAULT now(),
    ADD COLUMN IF NOT EXISTS deleted_at timestamptz;
