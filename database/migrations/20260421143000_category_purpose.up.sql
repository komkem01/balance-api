SET statement_timeout = 0;

--bun:split

ALTER TABLE categories
    ADD COLUMN IF NOT EXISTS purpose varchar;

--bun:split

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'categories_purpose_allowed'
    ) THEN
        ALTER TABLE categories
            ADD CONSTRAINT categories_purpose_allowed
            CHECK (purpose IS NULL OR purpose IN ('loan_repayment'));
    END IF;
END $$;

--bun:split

CREATE INDEX IF NOT EXISTS idx_categories_purpose ON categories(purpose);