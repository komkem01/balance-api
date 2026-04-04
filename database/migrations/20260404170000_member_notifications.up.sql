CREATE TABLE IF NOT EXISTS member_notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_id UUID NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('budget', 'security', 'weekly')),
    level VARCHAR(20) NOT NULL CHECK (level IN ('info', 'warning', 'critical')),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    dedupe_key VARCHAR(255),
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    read_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_member_notifications_member_created
    ON member_notifications(member_id, created_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_member_notifications_member_unread
    ON member_notifications(member_id, is_read)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_member_notifications_member_type_dedupe
    ON member_notifications(member_id, type, dedupe_key)
    WHERE deleted_at IS NULL AND dedupe_key IS NOT NULL;
