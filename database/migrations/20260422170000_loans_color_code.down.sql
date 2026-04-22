SET statement_timeout = 0;

--bun:split

ALTER TABLE loans
    DROP COLUMN IF EXISTS color_code;
