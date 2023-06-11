-- +goose Up
-- +goose StatementBegin
drop table unweave.filesystem_version;
drop table unweave.filesystem;

create table unweave.volume
(
    id         text primary key                              default ('vol_'::text || public.nanoid()) not null,
    name       text                                 not null,
    project_id text references unweave.project (id) not null,
    provider   text                                 not null,
    created_at timestamp                            not null default now(),
    updated_at timestamp                            not null default now()
);

alter table unweave.exec
    drop column persist_fs;

create table unweave.exec_volume
(
    exec_id    text references unweave.exec (id)   not null,
    volume_id  text references unweave.volume (id) not null,
    mount_path text                                not null,

    primary key (exec_id, volume_id, mount_path)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd
