SET statement_timeout = 0;

--bun:split

alter table transactions
    alter column amount type numeric
    using amount::numeric;

--bun:split

alter table budgets
    alter column amount type numeric
    using amount::numeric;
