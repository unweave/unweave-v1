-- +goose Up
-- +goose StatementBegin
alter table unweave.node rename column spec to metadata;

drop function if exists unweave.insert_node(text, text, text, jsonb, text, text, text[]);

-- Rename spec parameter

create function unweave.insert_node(v_node_id text, v_provider text, v_region text, v_metadata jsonb, v_status text, v_owner_id text, v_ssh_key_ids text[]) returns void
    language plpgsql
as
$$
begin
    insert into unweave.node (id, provider, region, metadata, status, owner_id)
    values (v_node_id, v_provider, v_region, v_metadata, v_status, v_owner_id);

    if v_ssh_key_ids is not null and array_upper(v_ssh_key_ids, 1) is not null then
        for i in 1 .. array_upper(v_ssh_key_ids, 1)
            loop
                insert into unweave.node_ssh_key (node_id, ssh_key_id)
                values (v_node_id, v_ssh_key_ids[i]);
            end loop;
    end if;
end;
$$;

comment on function unweave.insert_node(text, text, text, jsonb, text, text, text[]) is 'Insert a new node and associated ssh key';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
