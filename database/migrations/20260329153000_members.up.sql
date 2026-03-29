SET statement_timeout = 0;

--bun:split

create table members (
    id uuid primary key default gen_random_uuid(),
    gender_id uuid references genders(id),
    prefix_id uuid references prefixes(id),
    first_name varchar,
    last_name varchar,
    display_name varchar,
    phone varchar,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    last_login timestamptz,
    deleted_at timestamptz
);
