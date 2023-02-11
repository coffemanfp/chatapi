create table if not exists account (
    id serial unique not null,
    name varchar not null,
    last_name varchar,
    nickname varchar unique,
    email varchar unique,
    password varchar,
    picture_url varchar,
    created_at timestamptz not null,
    updated_at timestamptz,
    deleted_at timestamptz,

    primary key (id)
);

create table if not exists blocked_account (
	id serial unique not null,
	from_account_id integer unique not null,
	to_account_id integer unique not null,
	created_at timestamptz not null,

	primary key (id),
	foreign key (from_account_id) references account(id),
	foreign key (to_account_id) references account(id)
);

create table if not exists contact (
	id serial unique not null,
	from_account_id integer not null,
	to_account_id integer not null,
	name varchar not null,
	last_name varchar,
	created_at timestamptz not null,

	primary key (id),
	foreign key (from_account_id) references account(id),
	foreign key (to_account_id) references account(id)
);

create table if not exists conversation (
	id serial unique not null,
	capacity_members integer not null,
	name varchar not null,
	picture_url varchar not null,
	created_at timestamptz not null,
	deleted_at timestamptz,

	primary key (id)
);

create table if not exists conversation_role_permissions (
	id serial unique not null,
	write boolean not null,
	kick_account boolean not null,
	add_account boolean not null,
	change_role boolean not null,
	change_conversation_detail boolean not null,

	primary key (id)
);

create table if not exists convesation_members (
	id serial unique not null,
	account_id integer not null,
	conversation_id integer not null,
	joined_at timestamptz not null,
	left_at timestamptz,
	role_id integer not null,

	primary key (id),
	foreign key (account_id) references account(id),
	foreign key (conversation_id) references conversation(id)
);

create table if not exists conversation_role (
	id serial unique not null,
	permissions_id integer not null,
	name varchar not null,
	description varchar,

	primary key (id),
	foreign key (permissions_id) references conversation_role_permissions(id)
);

create table if not exists account_session (
    id varchar unique not null,
    account_id integer not null,
    logged_at timestamptz not null,
    last_seen_at timestamptz not null,
    logged_with varchar,
    actived boolean,

    primary key (id),
    foreign key (account_id) references account(id)
);

create unique index idx_account_id_actived on account_session(account_id, actived);