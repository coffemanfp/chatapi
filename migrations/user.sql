create table if not exists users (
    id serial unique not null,
    nickname varchar unique not null,
    password varchar not null,
    created_at timestamp not null,

    primary key (id)
);

create table if not exists user_session (
    id varchar unique not null,
    user_id integer not null,
    logged_at timestamp not null,
    last_seen_at timestamp not null,

    primary key (id),
    foreign key (user_id) references users(id)
);