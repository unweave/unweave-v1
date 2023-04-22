-- +goose Up
-- +goose StatementBegin
alter table unweave.filesystem
    add column src_path text;

update unweave.filesystem
    set src_path = '/tmp' where src_path is null;

alter table unweave.filesystem
    alter column src_path set not null;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd
