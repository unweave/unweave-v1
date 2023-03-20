-- +goose Up
-- +goose StatementBegin

with u as (
    insert into unweave.account (id)
        values ('00000000-0000-0000-0000-000000000001')
        on conflict (id) do nothing returning id)
insert
into unweave.project(id)
select '00000000-0000-0000-0000-000000000002'
where not exists(select 1 from unweave.project where id = '00000000-0000-0000-0000-000000000002');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd

