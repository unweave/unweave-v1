-- +goose Up
-- +goose StatementBegin

alter type unweave.session_status add value 'snapshotting';
alter type unweave.build_status add value 'syncing_snapshot';

create table unweave.filesystem
(
    id         text primary key                              default 'fs_' || nanoid() check (length(id) > 11),
    name       text                                 not null,
    project_id text references unweave.project (id) not null,
    exec_id    text references unweave.session (id) not null,
    owner_id   text references unweave.account (id) not null,
    created_at timestamptz                          not null default now(),

    constraint filesystem_unique_name_per_project unique (project_id, name)
);

create table unweave.filesystem_version
(
    filesystem_id text        not null references unweave.filesystem (id),
    exec_id       text        not null references unweave.session (id),
    version       int         not null,
    created_at    timestamptz not null default now(),
    build_id      text references unweave.build (id),

    constraint filesystem_version_unique unique (filesystem_id, version)
);

create or replace function unweave.insert_filesystem_version(p_filesystem_id text, p_exec_id text)
    returns unweave.filesystem_version as
$$
declare
    v_next_version           int;
    v_new_filesystem_version unweave.filesystem_version;
begin
    select coalesce(max(version), -1) + 1
    into v_next_version
    from unweave.filesystem_version
    where filesystem_id = p_filesystem_id;

    -- Insert a new row with the incremented version
    insert into unweave.filesystem_version (filesystem_id, exec_id, version)
    values (p_filesystem_id, p_exec_id, v_next_version)
    returning * into v_new_filesystem_version;

    return v_new_filesystem_version;
end;
$$ language plpgsql;


create or replace function unweave.auto_insert_version_zero() returns trigger as
$$
begin
    perform unweave.insert_filesystem_version(new.id, new.exec_id);
    return new;
end;
$$ language plpgsql;
comment on function unweave.auto_insert_version_zero() is 'Automatically add version 0 when a new filesystem is created.';

create trigger auto_insert_version_zero_trigger
    after insert
    on unweave.filesystem
    for each row
execute function unweave.auto_insert_version_zero();


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd
