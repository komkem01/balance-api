SET statement_timeout = 0;

--bun:split

ALTER TABLE loans
    ADD COLUMN IF NOT EXISTS color_code varchar;

--bun:split

UPDATE loans
SET color_code = '#6366f1'
WHERE color_code IS NULL OR btrim(color_code) = '';

--bun:split

ALTER TABLE loans
    ALTER COLUMN color_code SET DEFAULT '#6366f1';

--bun:split

ALTER TABLE loans
    ALTER COLUMN color_code SET NOT NULL;
