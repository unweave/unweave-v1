-- +goose Up
-- +goose StatementBegin

create schema unweave;

grant all on schema unweave to postgres;
grant all on all tables in schema unweave to postgres;


create table unweave.account
(
    id uuid primary key not null default gen_random_uuid()
);


create table unweave.project
(
    id         text primary key                              default 'pr_' || nanoid() check ( length(id) > 11 ),
    name       text                                 not null check ( name <> '' ),
    icon       text                                 not null default '📦',
    owner_id   uuid references unweave.account (id) not null,
    created_at timestamptz                          not null default now()
);

create type unweave.session_status as enum ('initializing', 'running', 'terminated', 'error');


create table unweave.ssh_key
(
    id         text primary key     default 'ss_' || nanoid() check ( length(id) > 11 ),
    name       text        not null,
    owner_id   uuid        not null references unweave.account (id),
    created_at timestamptz not null default now(),
    public_key text        not null,
    is_active  bool        not null default true,

    unique (name, owner_id)
);


create table unweave.session
(
    id              text primary key                default 'se_' || nanoid() check ( length(id) > 11 ),
    name            text                   not null default '',
    -- node_id is provider specific identifier of the compute resource assigned to this session.
    node_id         text                   not null,
    region          text                   not null,
    created_by      uuid                   not null references unweave.account (id),
    created_at      timestamptz            not null default now(),
    ready_at        timestamptz,
    exited_at       timestamptz,
    status          unweave.session_status not null default 'initializing',
    project_id      text                   not null references unweave.project (id),
    -- We don't want to constrain this to an enum to allow users to register their own
    -- providers without having to update the database schema.
    provider        text                   not null,
    ssh_key_id      text                   not null references unweave.ssh_key (id),
    connection_info jsonb                  not null default '{}',
    error           text
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd