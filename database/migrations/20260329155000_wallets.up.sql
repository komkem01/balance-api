SET statement_timeout = 0;

--bun:split

create table wallets (
    id uuid primary key default gen_random_uuid(),
    member_id uuid references members(id),
    name varchar,
    balance numeric not null default 0,
    currency varchar not null default 'THB',
    color_code varchar,
    is_active boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);
