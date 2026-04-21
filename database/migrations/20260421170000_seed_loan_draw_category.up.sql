SET statement_timeout = 0;

--bun:split

-- Seed system-level loan draw category (no member_id = applies to all members as default)
INSERT INTO categories (id, member_id, name, type, purpose, icon_name, color_code)
VALUES (
    'a1b2c3d4-0000-4000-8000-000000000002',
    NULL,
    'Withdraw a loan',
    'income',
    NULL,
    'wallet',
    '#0ea5e9'
)
ON CONFLICT (id) DO NOTHING;
