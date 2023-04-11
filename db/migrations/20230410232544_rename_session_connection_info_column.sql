-- +goose Up
-- +goose StatementBegin

alter table unweave.session
    add column if not exists metadata jsonb default '{}'::jsonb not null,
    drop column if exists connection_info;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd
