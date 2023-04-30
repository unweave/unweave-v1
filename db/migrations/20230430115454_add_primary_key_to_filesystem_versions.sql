-- +goose Up
-- +goose StatementBegin
alter table unweave.filesystem_version
    add primary key (filesystem_id, exec_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd
