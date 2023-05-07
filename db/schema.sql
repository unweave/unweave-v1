SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

CREATE SCHEMA unweave;

ALTER SCHEMA unweave OWNER TO postgres;

CREATE TYPE unweave.build_status AS ENUM (
    'initializing',
    'building',
    'success',
    'failed',
    'error',
    'canceled',
    'syncing_snapshot'
);

ALTER TYPE unweave.build_status OWNER TO postgres;

CREATE TYPE unweave.exec_status AS ENUM (
    'initializing',
    'running',
    'terminated',
    'error',
    'snapshotting'
);

ALTER TYPE unweave.exec_status OWNER TO postgres;

CREATE FUNCTION unweave.auto_insert_version_zero() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
begin
    perform unweave.insert_filesystem_version(new.id, new.exec_id);
    return new;
end;
$$;

ALTER FUNCTION unweave.auto_insert_version_zero() OWNER TO postgres;

COMMENT ON FUNCTION unweave.auto_insert_version_zero() IS 'Automatically add version 0 when a new filesystem is created.';

SET default_tablespace = '';

SET default_table_access_method = heap;

CREATE TABLE unweave.filesystem_version (
    filesystem_id text NOT NULL,
    exec_id text NOT NULL,
    version integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    build_id text
);

ALTER TABLE unweave.filesystem_version OWNER TO postgres;

CREATE FUNCTION unweave.insert_filesystem_version(p_filesystem_id text, p_exec_id text) RETURNS unweave.filesystem_version
    LANGUAGE plpgsql
    AS $$
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
$$;

ALTER FUNCTION unweave.insert_filesystem_version(p_filesystem_id text, p_exec_id text) OWNER TO postgres;

CREATE FUNCTION unweave.insert_node(v_node_id text, v_provider text, v_region text, v_metadata jsonb, v_status text, v_owner_id text, v_ssh_key_ids text[]) RETURNS void
    LANGUAGE plpgsql
    AS $$
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

ALTER FUNCTION unweave.insert_node(v_node_id text, v_provider text, v_region text, v_metadata jsonb, v_status text, v_owner_id text, v_ssh_key_ids text[]) OWNER TO postgres;

COMMENT ON FUNCTION unweave.insert_node(v_node_id text, v_provider text, v_region text, v_metadata jsonb, v_status text, v_owner_id text, v_ssh_key_ids text[]) IS 'Insert a new node and associated ssh key';

CREATE TABLE unweave.filesystem (
    id text DEFAULT ('fs_'::text || public.nanoid()) NOT NULL,
    name text NOT NULL,
    project_id text NOT NULL,
    exec_id text NOT NULL,
    owner_id text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    src_path text NOT NULL,
    CONSTRAINT filesystem_id_check CHECK ((length(id) > 11))
);

ALTER TABLE unweave.filesystem OWNER TO postgres;

CREATE TABLE unweave.build (
    id text DEFAULT ('bld_'::text || public.nanoid()) NOT NULL,
    name text NOT NULL,
    project_id text NOT NULL,
    builder_type text NOT NULL,
    status unweave.build_status DEFAULT 'initializing'::unweave.build_status NOT NULL,
    created_by text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    started_at timestamp with time zone,
    finished_at timestamp with time zone,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    meta_data jsonb DEFAULT '{}'::jsonb NOT NULL,
    CONSTRAINT build_id_check CHECK ((length(id) > 11))
);

ALTER TABLE unweave.build OWNER TO postgres;

CREATE TABLE unweave.exec (
    id text DEFAULT ('se_'::text || public.nanoid()) NOT NULL,
    name text DEFAULT ''::text NOT NULL,
    node_id text NOT NULL,
    region text NOT NULL,
    created_by text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    ready_at timestamp with time zone,
    exited_at timestamp with time zone,
    status unweave.exec_status DEFAULT 'initializing'::unweave.exec_status NOT NULL,
    project_id text NOT NULL,
    ssh_key_id text,
    error text,
    build_id text,
    spec jsonb DEFAULT '{}'::jsonb NOT NULL,
    commit_id text,
    git_remote_url text,
    command text[],
    metadata jsonb DEFAULT '{}'::jsonb NOT NULL,
    persist_fs boolean DEFAULT false NOT NULL,
    image text DEFAULT 'ubuntu:latest'::text NOT NULL,
    CONSTRAINT session_id_check CHECK ((length(id) > 11))
);

ALTER TABLE unweave.exec OWNER TO postgres;

CREATE TABLE unweave.project (
    id text DEFAULT ('pr_'::text || public.nanoid()) NOT NULL,
    default_build_id text,
    CONSTRAINT project_id_check CHECK ((length(id) > 11))
);

ALTER TABLE unweave.project OWNER TO postgres;

CREATE TABLE unweave.node (
    id text NOT NULL,
    provider text NOT NULL,
    region text NOT NULL,
    metadata jsonb NOT NULL,
    status text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    ready_at timestamp with time zone,
    owner_id text NOT NULL,
    terminated_at timestamp with time zone
);

ALTER TABLE unweave.node OWNER TO postgres;

CREATE TABLE unweave.ssh_key (
    id text DEFAULT ('ss_'::text || public.nanoid()) NOT NULL,
    name text NOT NULL,
    owner_id text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    public_key text NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    CONSTRAINT ssh_key_id_check CHECK ((length(id) > 11))
);

ALTER TABLE unweave.ssh_key OWNER TO postgres;

CREATE TABLE unweave.account (
    id text NOT NULL
);

ALTER TABLE unweave.account OWNER TO postgres;

CREATE TABLE unweave.node_ssh_key (
    node_id text NOT NULL,
    ssh_key_id text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);

ALTER TABLE unweave.node_ssh_key OWNER TO postgres;

ALTER TABLE ONLY unweave.account
    ADD CONSTRAINT account_pkey PRIMARY KEY (id);

ALTER TABLE ONLY unweave.build
    ADD CONSTRAINT build_pkey PRIMARY KEY (id);

ALTER TABLE ONLY unweave.filesystem
    ADD CONSTRAINT filesystem_pkey PRIMARY KEY (id);

ALTER TABLE ONLY unweave.filesystem
    ADD CONSTRAINT filesystem_unique_name_per_project UNIQUE (project_id, name);

ALTER TABLE ONLY unweave.filesystem_version
    ADD CONSTRAINT filesystem_version_pkey PRIMARY KEY (filesystem_id, exec_id);

ALTER TABLE ONLY unweave.filesystem_version
    ADD CONSTRAINT filesystem_version_unique UNIQUE (filesystem_id, version);

ALTER TABLE ONLY unweave.node
    ADD CONSTRAINT node_pkey PRIMARY KEY (id);

ALTER TABLE ONLY unweave.node_ssh_key
    ADD CONSTRAINT node_ssh_key_pk PRIMARY KEY (node_id, ssh_key_id);

ALTER TABLE ONLY unweave.project
    ADD CONSTRAINT project_pkey PRIMARY KEY (id);

ALTER TABLE ONLY unweave.exec
    ADD CONSTRAINT session_pkey PRIMARY KEY (id);

ALTER TABLE ONLY unweave.exec
    ADD CONSTRAINT session_project_id_name_key UNIQUE (project_id, name);

ALTER TABLE ONLY unweave.ssh_key
    ADD CONSTRAINT ssh_key_name_owner_id_key UNIQUE (name, owner_id);

ALTER TABLE ONLY unweave.ssh_key
    ADD CONSTRAINT ssh_key_pkey PRIMARY KEY (id);

CREATE TRIGGER auto_insert_version_zero_trigger AFTER INSERT ON unweave.filesystem FOR EACH ROW EXECUTE FUNCTION unweave.auto_insert_version_zero();

ALTER TABLE ONLY unweave.build
    ADD CONSTRAINT build_created_by_fkey FOREIGN KEY (created_by) REFERENCES unweave.account(id);

ALTER TABLE ONLY unweave.build
    ADD CONSTRAINT build_project_id_fkey FOREIGN KEY (project_id) REFERENCES unweave.project(id);

ALTER TABLE ONLY unweave.exec
    ADD CONSTRAINT exec_build_id_fkey FOREIGN KEY (build_id) REFERENCES unweave.build(id);

ALTER TABLE ONLY unweave.exec
    ADD CONSTRAINT exec_created_by_fkey FOREIGN KEY (created_by) REFERENCES unweave.account(id);

ALTER TABLE ONLY unweave.exec
    ADD CONSTRAINT exec_project_id_fkey FOREIGN KEY (project_id) REFERENCES unweave.project(id);

ALTER TABLE ONLY unweave.exec
    ADD CONSTRAINT exec_ssh_key_id_fkey FOREIGN KEY (ssh_key_id) REFERENCES unweave.ssh_key(id);

ALTER TABLE ONLY unweave.filesystem
    ADD CONSTRAINT filesystem_exec_id_fkey FOREIGN KEY (exec_id) REFERENCES unweave.exec(id);

ALTER TABLE ONLY unweave.filesystem
    ADD CONSTRAINT filesystem_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES unweave.account(id);

ALTER TABLE ONLY unweave.filesystem
    ADD CONSTRAINT filesystem_project_id_fkey FOREIGN KEY (project_id) REFERENCES unweave.project(id);

ALTER TABLE ONLY unweave.filesystem_version
    ADD CONSTRAINT filesystem_version_build_id_fkey FOREIGN KEY (build_id) REFERENCES unweave.build(id);

ALTER TABLE ONLY unweave.filesystem_version
    ADD CONSTRAINT filesystem_version_exec_id_fkey FOREIGN KEY (exec_id) REFERENCES unweave.exec(id);

ALTER TABLE ONLY unweave.filesystem_version
    ADD CONSTRAINT filesystem_version_filesystem_id_fkey FOREIGN KEY (filesystem_id) REFERENCES unweave.filesystem(id);

ALTER TABLE ONLY unweave.node
    ADD CONSTRAINT node_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES unweave.account(id);

ALTER TABLE ONLY unweave.node_ssh_key
    ADD CONSTRAINT node_ssh_key_node_id_fkey FOREIGN KEY (node_id) REFERENCES unweave.node(id);

ALTER TABLE ONLY unweave.node_ssh_key
    ADD CONSTRAINT node_ssh_key_ssh_key_id_fkey FOREIGN KEY (ssh_key_id) REFERENCES unweave.ssh_key(id);

ALTER TABLE ONLY unweave.project
    ADD CONSTRAINT project_default_build_id_fkey FOREIGN KEY (default_build_id) REFERENCES unweave.build(id);

ALTER TABLE ONLY unweave.ssh_key
    ADD CONSTRAINT ssh_key_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES unweave.account(id);

