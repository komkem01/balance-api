SET statement_timeout = 0;

--bun:split

create type transaction_type as enum ('income', 'expense');

--bun:split

create table transactions (
    id uuid primary key default gen_random_uuid(),
    wallet_id uuid references wallets(id),
    category_id uuid references categories(id),
    amount numeric not null default 0,
    type transaction_type not null,
    transaction_date date,
    note text,
    image_url varchar,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);
