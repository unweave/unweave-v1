-- +goose Up
-- +goose StatementBegin
alter table unweave.session
    add column image text not null default 'ubuntu:latest';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
