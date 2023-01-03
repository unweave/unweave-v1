-- +goose Up
-- +goose StatementBegin

with u as (
    insert into unweave.users (id)
        values ('00000000-0000-0000-0000-000000000001')
        on conflict (id) do nothing returning id)
insert
into unweave.projects(id, name, owner_id)
select '00000000-0000-0000-0000-000000000002','default-project', (select id from u)
where not exists(select 1 from unweave.projects where name = 'default-project')

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd

