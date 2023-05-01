-- +goose Up
-- +goose StatementBegin

alter table unweave.session
    add unique (project_id, name);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd
