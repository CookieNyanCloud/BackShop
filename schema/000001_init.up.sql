CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_cron";

CREATE TABLE if not exists users
(
    id            varchar(255) NOT NULL PRIMARY KEY,
    email         varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE if not exists admins
(
    id            varchar(255) NOT NULL PRIMARY KEY,
    email         varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE if not exists verification
(
    id    varchar(255) primary key references users (id) on delete cascade,
    code  int,
    state bool
);

CREATE TABLE if not exists sessions
(
    id           varchar(255) primary key references users (id) on delete cascade,
    refreshtoken varchar(255),
    expiresat    date,
    lastvisitat  date
);

CREATE TABLE if not exists events
(
    id          serial primary key not null unique,
    time        date               not null,
    zoneAmount  int,
    description text,
    mapfile     varchar(255)
);

CREATE TABLE if not exists zones
(
    id      int          not null,
    eventId int          not null references events (id) on delete cascade,
    taken   varchar(255) not null references users (id) on delete cascade,
    price   int          not null
);

CREATE TABLE if not exists orders
(
    id      uuid         not null,
    userId  varchar(255) not null references users (id) on DELETE cascade,
    eventId int          not null references events (id) on DELETE cascade,
    zonesId int[]        not null references zones (id) on DELETE cascade,
    expire  date         not null
);

CREATE OR REPLACE FUNCTION untake_zones()
    RETURNS TRIGGER
    LANGUAGE PLPGSQL
AS
$$
BEGIN

END;
$$



