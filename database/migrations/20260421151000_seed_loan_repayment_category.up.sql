SET statement_timeout = 0;

--bun:split

-- Seed system-level loan repayment category (no member_id = applies to all members as default)
INSERT INTO categories (id, member_id, name, type, purpose, icon_name, color_code)
VALUES (
    'a1b2c3d4-0000-4000-8000-000000000001',
    NULL,
    'ชำระเงินกู้',
    'expense',
    'loan_repayment',
    'wallet',
    '#6366f1'
)
ON CONFLICT (id) DO NOTHING;