-- +goose Up
-- +goose StatementBegin
create table if not exists unweave.node
(
    id         text primary key not null,
    provider   text             not null,
    region     text             not null,
    spec       jsonb            not null,
    status     text,
    created_at timestamptz,
    ready_at   timestamptz,
    owner_id   text references unweave.account (id)
);

alter table unweave.session
    drop column if exists provider,
    add column if not exists spec jsonb not null default '{}'::jsonb;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd
