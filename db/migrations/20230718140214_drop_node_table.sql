-- +goose Up
-- +goose StatementBegin
drop table if exists unweave.node_ssh_key;
drop table if exists unweave.node;
drop function if exists unweave.insert_node;

alter table unweave.exec
    alter column id drop default;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
