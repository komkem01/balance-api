SET statement_timeout = 0;

--bun:split

create type budget_period as enum ('daily', 'weekly', 'monthly');

--bun:split

create table budgets (
    id uuid primary key default gen_random_uuid(),
    member_id uuid references members(id),
    category_id uuid references categories(id),
    amount numeric not null default 0,
    period budget_period not null,
    start_date date,
    end_date date,
    created_at timestamptz not null default now()
);
