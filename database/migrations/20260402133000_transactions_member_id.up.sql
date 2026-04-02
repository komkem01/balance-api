SET statement_timeout = 0;

--bun:split

ALTER TABLE transactions
    ADD COLUMN IF NOT EXISTS member_id uuid;

--bun:split

UPDATE transactions AS t
SET member_id = w.member_id
FROM wallets AS w
WHERE t.wallet_id = w.id
  AND t.member_id IS NULL;

--bun:split

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'transactions_member_id_fkey'
    ) THEN
        ALTER TABLE transactions
            ADD CONSTRAINT transactions_member_id_fkey
            FOREIGN KEY (member_id) REFERENCES members(id);
    END IF;
END $$;

--bun:split

CREATE INDEX IF NOT EXISTS idx_transactions_member_id
    ON transactions (member_id);

--bun:split

CREATE OR REPLACE FUNCTION set_transactions_member_id_from_wallet()
RETURNS trigger AS $$
BEGIN
    IF NEW.wallet_id IS NULL THEN
        NEW.member_id := NULL;
        RETURN NEW;
    END IF;

    SELECT w.member_id
    INTO NEW.member_id
    FROM wallets AS w
    WHERE w.id = NEW.wallet_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

--bun:split

DROP TRIGGER IF EXISTS trg_transactions_set_member_id ON transactions;

--bun:split

CREATE TRIGGER trg_transactions_set_member_id
BEFORE INSERT OR UPDATE OF wallet_id ON transactions
FOR EACH ROW
EXECUTE FUNCTION set_transactions_member_id_from_wallet();
