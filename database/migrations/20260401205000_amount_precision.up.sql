SET statement_timeout = 0;

--bun:split

alter table transactions
    alter column amount type numeric(18,2)
    using round(amount::numeric, 2);

--bun:split

alter table budgets
    alter column amount type numeric(18,2)
    using round(amount::numeric, 2);
