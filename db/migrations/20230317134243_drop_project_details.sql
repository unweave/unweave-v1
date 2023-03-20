-- +goose Up
-- +goose StatementBegin
alter table unweave.project
    drop column name,
    drop column icon,
    drop column owner_id,
    drop column created_at,
    drop column default_build;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd
