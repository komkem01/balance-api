SET statement_timeout = 0;

--bun:split

ALTER TABLE members
    ADD COLUMN profile_image_url varchar NOT NULL DEFAULT '';
