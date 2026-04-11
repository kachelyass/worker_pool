create table if not exists users_projection (
    id uuid primary key,
    name varchar(25) not null,
    lastname varchar(25) not null
);

create table if not exists user_processed (
    event_id uuid primary key,
    created_at timestamptz not null default now()
);