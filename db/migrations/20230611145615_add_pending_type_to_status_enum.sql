-- +goose Up
-- +goose StatementBegin
alter type unweave.exec_status add value 'pending';
commit; -- require for enum type before they're used (below)

-- set default status to pending
alter table unweave.exec
    alter column status drop not null,
    alter column status drop default;

alter table unweave.exec
    alter column status set default 'pending'::unweave.exec_status,
    alter column status set not null;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd
