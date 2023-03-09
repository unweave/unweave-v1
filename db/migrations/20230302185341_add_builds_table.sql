-- +goose Up
-- +goose StatementBegin

create type unweave.build_status as enum ('initializing', 'building', 'success', 'failed', 'error', 'canceled');

create table unweave.build
(
    id           text primary key                              default 'bld_' || nanoid() check ( length(id) > 11 ),
    project_id   text references unweave.project (id) not null,
    builder_type text                                 not null,
    status       unweave.build_status                 not null default 'initializing',
    created_at   timestamptz                          not null default now(),
    updated_at   timestamptz                          not null default now(),
    meta_data    jsonb                                not null default '{}'::jsonb
);

alter table unweave.project
    add column default_build text references unweave.build (id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd
