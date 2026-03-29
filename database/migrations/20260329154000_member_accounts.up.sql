SET statement_timeout = 0;

--bun:split

create table member_accounts (
    id uuid primary key default gen_random_uuid(),
    member_id uuid references members(id),
    username varchar,
    password text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);
