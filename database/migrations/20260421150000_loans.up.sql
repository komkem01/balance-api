SET statement_timeout = 0;

--bun:split

CREATE TABLE loans (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    member_id uuid REFERENCES members(id),
    name varchar NOT NULL,
    lender varchar,
    total_amount numeric NOT NULL DEFAULT 0,
    remaining_balance numeric NOT NULL DEFAULT 0,
    monthly_payment numeric NOT NULL DEFAULT 0,
    interest_rate numeric NOT NULL DEFAULT 0,
    start_date date,
    end_date date,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);