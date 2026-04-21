SET statement_timeout = 0;

--bun:split

CREATE TABLE goal_entries (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    goal_id uuid NOT NULL REFERENCES goals(id),
    member_id uuid REFERENCES members(id),
    source_type varchar NOT NULL,
    source_id uuid,
    amount_before numeric NOT NULL DEFAULT 0,
    amount_after numeric NOT NULL DEFAULT 0,
    amount_delta numeric NOT NULL DEFAULT 0,
    note varchar NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_goal_entries_goal_id_created_at ON goal_entries(goal_id, created_at DESC);
CREATE INDEX idx_goal_entries_member_id_created_at ON goal_entries(member_id, created_at DESC);
