CREATE TABLE if not exists users
(
    id            serial primary key not null unique,
    name          varchar(255)       not null unique,
    email         varchar(255)       not null unique,
    password_hash varchar(255)       not null,
    verification  boolean,
    zone          int

);

CREATE TABLE if not exists codes
(
    id   serial primary key not null unique,
    code int

);

CREATE TABLE if not exists admins
(
    id            serial primary key not null unique,
    name          varchar(255)       not null unique,
    email         varchar(255)       not null unique,
    password_hash varchar(255)       not null
);

CREATE TABLE if not exists events
(
    id   serial primary key not null unique,
    time date               not null
);

CREATE TABLE if not exists sessions
(
    id           int  not null unique,
    refreshtoken varchar(255),
    expiresat    date,
    lastvisitat  date

);

CREATE TABLE if not exists zones
(
    eventId int                not null references events (id) on delete cascade,
    id      serial primary key not null,
    taken   int                not null,
    price   int                not null
);



INSERT INTO events (id, time)
VALUES (0, date('2021-07-30'));

INSERT INTO zones (eventId, id, taken, price)
VALUES (0, 1, 1, 100),
       (0, 2, 0, 200),
       (0, 3, 0, 300),
       (0, 4, 0, 400),
       (0, 5, 1, 500),
       (0, 6, 0, 600);

