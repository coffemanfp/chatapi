create table if not exists users (
    id serial unique not null,
    nickname varchar unique not null,
    password varchar not null,
    created_at timestamp not null,

    primary key (id)
);