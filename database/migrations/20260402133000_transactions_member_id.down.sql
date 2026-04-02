SET statement_timeout = 0;

--bun:split

DROP TRIGGER IF EXISTS trg_transactions_set_member_id ON transactions;

--bun:split

DROP FUNCTION IF EXISTS set_transactions_member_id_from_wallet();

--bun:split

ALTER TABLE transactions
    DROP CONSTRAINT IF EXISTS transactions_member_id_fkey;

--bun:split

DROP INDEX IF EXISTS idx_transactions_member_id;

--bun:split

ALTER TABLE transactions
    DROP COLUMN IF EXISTS member_id;
