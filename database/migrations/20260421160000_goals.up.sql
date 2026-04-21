SET statement_timeout = 0;

--bun:split

CREATE TABLE goals (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    member_id uuid REFERENCES members(id),
    name varchar NOT NULL,
    type varchar NOT NULL,
    target_amount numeric NOT NULL DEFAULT 0,
    start_amount numeric NOT NULL DEFAULT 0,
    current_amount numeric NOT NULL DEFAULT 0,
    start_date date,
    target_date date,
    status varchar NOT NULL DEFAULT 'active',
    auto_tracking boolean NOT NULL DEFAULT true,
    tracking_source_type varchar,
    tracking_source_id uuid,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);
