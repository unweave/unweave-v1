-- +goose Up
-- +goose StatementBegin

alter table unweave.build
    alter column created_by type text;
alter table unweave.session
    alter column created_by type text;
alter table unweave.ssh_key
    alter column owner_id type text;
alter table unweave.account
    alter column id type text;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
