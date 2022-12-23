-- migrate:up
create schema unweave;
grant all on schema unweave to postgres;

-- Create tables
create table unweave.project
(
    id   uuid primary key default gen_random_uuid(),
    name text not null
);

create table unweave.project_token
(
    id           uuid primary key    default gen_random_uuid(),
    name         text       not null,
    hash         text       not null,
    display_name text       not null,
    project_id   uuid       not null references unweave.project (id),
    created_at   timestamptz not null default now(),
    expires_at   timestamptz not null default now() + interval '7 day'
);

create table unweave.session
(
    id         uuid primary key    default gen_random_uuid(),
    created_at timestamptz not null default now(),
    ready_at   timestamptz,
    exited_at  timestamptz,
    project_id uuid       not null references unweave.project (id),
    runtime    text       not null -- we don't want to constrain this to a specific set of values
)

-- migrate:down

