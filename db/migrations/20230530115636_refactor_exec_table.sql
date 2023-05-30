-- +goose Up
-- +goose StatementBegin
alter table unweave.exec
    add column provider text;

update unweave.exec
set provider = 'unweave'
where exec.provider is null;

alter table unweave.exec
    alter column provider set not null;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd
