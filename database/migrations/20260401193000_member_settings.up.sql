SET statement_timeout = 0;

--bun:split

ALTER TABLE members
ADD COLUMN preferred_currency varchar NOT NULL DEFAULT 'THB',
ADD COLUMN preferred_language varchar NOT NULL DEFAULT 'EN',
ADD COLUMN notify_budget boolean NOT NULL DEFAULT true,
ADD COLUMN notify_security boolean NOT NULL DEFAULT true,
ADD COLUMN notify_weekly boolean NOT NULL DEFAULT false;
