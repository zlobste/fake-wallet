-- +migrate Up
create table users
(
    id       bigserial primary key not null,
    name     varchar(50)           not null,
    surname  varchar(50)           not null,
    email    varchar(50)           not null,
    password varchar(200)          not null
);

create table assets
(
    id     bigserial primary key         not null,
    symbol varchar(10)                   not null,
    name   varchar(50)                   not null,
    fee    numeric(5, 2) check (fee > 0) not null
);

create table wallets
(
    address  text primary key                 not null,
    asset_id bigint references assets (id)    not null,
    balance  numeric(32) check (balance >= 0) not null,
    owner_id bigint references users (id)     not null
);

create table transactions
(
    id       bigserial primary key             not null,
    sender   text references wallets (address) not null,
    receiver text references wallets (address) not null,
    amount   numeric(32) check (amount > 0)    not null,
    fee      numeric(32) check (fee > 0)       not null,
    time     timestamp without time zone       not null
);

-- +migrate Down
drop table users cascade;
drop table assets cascade;
drop table wallets cascade;
drop table transactions cascade;