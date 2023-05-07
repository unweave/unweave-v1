-- +goose Up
-- +goose StatementBegin

-- Drop constraints

alter table only unweave.session
    drop constraint session_build_fkey;
alter table only unweave.session
    drop constraint session_created_by_fkey;
alter table only unweave.session
    drop constraint session_project_id_fkey;
alter table only unweave.session
    drop constraint session_ssh_key_id_fkey;

-- Rename table
alter table unweave.session
    rename to exec;

-- Add back constraints

alter table only unweave.exec
    add constraint exec_build_id_fkey foreign key (build_id) references unweave.build (id);
alter table only unweave.exec
    add constraint exec_created_by_fkey foreign key (created_by) references unweave.account (id);
alter table only unweave.exec
    add constraint exec_project_id_fkey foreign key (project_id) references unweave.project (id);
alter table only unweave.exec
    add constraint exec_ssh_key_id_fkey foreign key (ssh_key_id) references unweave.ssh_key (id);

set search_path to unweave;
alter type session_status rename to exec_status;
reset search_path;

-- change type of unweave.exec.status to exec_status
alter table unweave.exec
    alter column status type unweave.exec_status using status::unweave.exec_status;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd
