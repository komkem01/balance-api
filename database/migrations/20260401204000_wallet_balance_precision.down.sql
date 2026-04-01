SET statement_timeout = 0;

--bun:split

alter table wallets
    alter column balance type numeric
    using balance::numeric;
