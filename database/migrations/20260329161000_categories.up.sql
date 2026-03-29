SET statement_timeout = 0;

--bun:split

create type category_type as enum ('income', 'expense');

--bun:split

create table categories (
    id uuid primary key default gen_random_uuid(),
    member_id uuid references members(id),
    name varchar,
    type category_type not null,
    icon_name varchar,
    color_code varchar,
    created_at timestamptz not null default now()
);
