-- +goose Up
-- +goose StatementBegin
alter table unweave.session
    add column persist_fs boolean not null default false;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd
