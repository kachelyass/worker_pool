CREATE TABLE users (
    id uuid PRIMARY KEY,
    name varchar(25) not null,
    lastname varchar(25) not null
);

create table if not exists user_outbox (
    id uuid primary key,
    topic text not null,
    event_key text not null,
    payload jsonb not null,
    status text not null default 'pending'
        check (status in ('pending', 'published', 'failed')),
    attempts integer not null default 0
        check (attempts >= 0),
    available_at timestamptz not null default now(),
    created_at timestamptz not null default now()
);