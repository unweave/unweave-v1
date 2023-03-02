-- +goose Up
-- +goose StatementBegin

create type unweave.build_status as enum ('initializing', 'building', 'success', 'error', 'canceled');

create table unweave.build
(
    id           text primary key                              default 'bl_' || nanoid() check ( length(id) > 11 ),
    project_id   text references unweave.project (id) not null,
    builder_type text                                 not null,
    status       unweave.build_status                 not null default 'initializing',
    created_at   timestamptz                          not null default now(),
    meta_data    jsonb                                not null default '{}'::jsonb
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd
