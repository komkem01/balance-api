SET statement_timeout = 0;

--bun:split

ALTER TABLE goals
    DROP COLUMN IF EXISTS deposit_wallet_id;
