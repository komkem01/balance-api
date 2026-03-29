SET statement_timeout = 0;

--bun:split

create table prefixes (
    id uuid primary key default gen_random_uuid(),
    gender_id uuid not null references genders(id) on delete cascade,
    name text not null,
    is_active boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
)
