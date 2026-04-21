SET statement_timeout = 0;

--bun:split

ALTER TABLE goals
    ADD COLUMN deposit_wallet_id uuid REFERENCES wallets(id);
