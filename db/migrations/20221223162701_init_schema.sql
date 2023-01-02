-- migrate:up
create schema unweave;
grant all on schema unweave to postgres;

-- Minimal users table to allow for constraints
create table unweave.users
(
    id integer primary key generated always as identity
);

create table unweave.projects
(
    id       uuid primary key default gen_random_uuid(),
    name     text                                  not null,
    owner_id integer references unweave.users (id) not null
);

create table unweave.sessions
(
    id         uuid primary key     default gen_random_uuid(),
    created_by integer     not null references unweave.users (id),
    created_at timestamptz not null default now(),
    ready_at   timestamptz,
    exited_at  timestamptz,
    project_id uuid        not null references unweave.projects (id),
    runtime    text        not null -- we don't want to constrain this to a specific set of values
);

create table unweave.ssh_keys
(
    id          uuid primary key     default gen_random_uuid(),
    name        text        not null,
    owner_id    integer     not null references unweave.users (id),
    created_at  timestamptz not null default now(),
    public_key text        not null unique
);

-- migrate:down

