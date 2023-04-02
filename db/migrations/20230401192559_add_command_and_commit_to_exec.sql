-- +goose Up
-- +goose StatementBegin
alter table unweave.session
    add column commit_id      text,
    add column git_remote_url text,
    add column command        text[];

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd
