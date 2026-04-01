SET statement_timeout = 0;

--bun:split

alter table wallets
    alter column balance type numeric(18,2)
    using round(balance::numeric, 2);
