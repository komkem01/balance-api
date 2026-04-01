SET statement_timeout = 0;

--bun:split

ALTER TABLE members
DROP COLUMN IF EXISTS notify_weekly,
DROP COLUMN IF EXISTS notify_security,
DROP COLUMN IF EXISTS notify_budget,
DROP COLUMN IF EXISTS preferred_language,
DROP COLUMN IF EXISTS preferred_currency;
