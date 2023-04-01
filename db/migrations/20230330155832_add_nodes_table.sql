-- +goose Up
-- +goose StatementBegin
create table if not exists unweave.node
(
    id         text primary key not null,
    provider   text             not null,
    region     text             not null,
    spec       jsonb            not null,
    status     text             not null,
    created_at timestamptz      not null default now(),
    ready_at   timestamptz,
    owner_id   text             not null references unweave.account (id)
);

create table if not exists unweave.node_ssh_key
(
    node_id    text             not null references unweave.node (id),
    ssh_key_id text             not null references unweave.ssh_key (id),
    created_at timestamptz      not null default now(),

    constraint node_ssh_key_pk primary key (node_id, ssh_key_id)
);

create or replace function unweave.insert_node(
    v_node_id text,
    v_provider text,
    v_region text,
    v_spec jsonb,
    v_status text,
    v_owner_id text,
    v_ssh_key_ids text[]
)
    returns void
    language plpgsql
as $$
begin
    insert into unweave.node (id, provider, region, spec, status, owner_id)
    values (v_node_id, v_provider, v_region, v_spec, v_status, v_owner_id);

    for i in 1 .. array_upper(v_ssh_key_ids, 1)
    loop
        insert into unweave.node_ssh_key (node_id, ssh_key_id)
        values (v_node_id, v_ssh_key_ids[i]);
    end loop;
end;
$$;
comment on function unweave.insert_node(text, text, text, jsonb, text, text, text[]) is
    'Insert a new node and associated ssh key';

alter table unweave.session
    drop column if exists provider,
    alter column ssh_key_id drop not null,
    add column if not exists spec jsonb not null default '{}'::jsonb;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd
