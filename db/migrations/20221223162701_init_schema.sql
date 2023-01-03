-- +goose Up
-- +goose StatementBegin

create schema unweave;
grant all on schema unweave to postgres;

-- Minimal users table to allow for constraints
create table unweave.users
(
    id uuid primary key not null default gen_random_uuid()
);

create table unweave.projects
(
    id         uuid primary key                            default gen_random_uuid(),
    name       text                               not null,
    owner_id   uuid references unweave.users (id) not null,
    created_at timestamptz                        not null default now()
);

create type unweave.session_status as enum ('initializing', 'active', 'terminated');

create table unweave.sessions
(
    id         uuid primary key                default gen_random_uuid(),
    -- node_id is provider specific identifier of the compute resource assigned to this session.
    node_id    text                   not null,
    created_by uuid                   not null references unweave.users (id),
    created_at timestamptz            not null default now(),
    ready_at   timestamptz,
    exited_at  timestamptz,
    status     unweave.session_status not null default 'initializing',
    project_id uuid                   not null references unweave.projects (id),
    -- We don't want to constrain this to an enum to allow users to register their own
    -- providers without having to update the database schema.
    runtime    text                   not null
);

create table unweave.ssh_keys
(
    id         uuid primary key     default gen_random_uuid(),
    name       text        not null,
    owner_id   uuid        not null references unweave.users (id),
    created_at timestamptz not null default now(),
    public_key text        not null unique
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd