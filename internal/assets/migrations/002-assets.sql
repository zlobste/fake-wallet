-- +migrate Up
insert into assets (symbol, name, fee)
values ('fETH', 'Fake ethereum', 0.05),
       ('fBTC', 'Fake bitcoin', 0.2);
