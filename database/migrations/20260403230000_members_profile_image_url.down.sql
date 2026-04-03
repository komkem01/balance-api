SET statement_timeout = 0;

--bun:split

ALTER TABLE members
    DROP COLUMN IF EXISTS profile_image_url;
