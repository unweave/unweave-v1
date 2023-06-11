-- +goose Up
-- +goose StatementBegin
alter table unweave.exec
    add column provider text;

create table unweave.exec_ssh_key (
    exec_id text not null references unweave.exec(id),
    ssh_key_id text references unweave.ssh_key(id),

    primary key (exec_id, ssh_key_id)
);

update unweave.exec
set provider = 'unweave'
where exec.provider is null;

alter table unweave.exec
    alter column provider set not null,
    drop column node_id,
    drop column ssh_key_id;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd
