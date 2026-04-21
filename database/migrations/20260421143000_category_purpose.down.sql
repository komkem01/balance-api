SET statement_timeout = 0;

--bun:split

DROP INDEX IF EXISTS idx_categories_purpose;

--bun:split

ALTER TABLE categories
    DROP CONSTRAINT IF EXISTS categories_purpose_allowed;

--bun:split

ALTER TABLE categories
    DROP COLUMN IF EXISTS purpose;