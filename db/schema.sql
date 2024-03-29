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
    'snapshotting',
    'pending'
);

ALTER TYPE unweave.exec_status OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

CREATE TABLE unweave.build (
    id text DEFAULT ('bld_'::text || public.nanoid()) NOT NULL,
    project_id text NOT NULL,
    builder_type text NOT NULL,
    status unweave.build_status DEFAULT 'initializing'::unweave.build_status NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    meta_data jsonb DEFAULT '{}'::jsonb NOT NULL,
    name text NOT NULL,
    started_at timestamp with time zone,
    finished_at timestamp with time zone,
    created_by text NOT NULL,
    CONSTRAINT build_id_check CHECK ((length(id) > 11))
);

ALTER TABLE unweave.build OWNER TO postgres;

CREATE TABLE unweave.exec (
    id text NOT NULL,
    name text DEFAULT ''::text NOT NULL,
    region text NOT NULL,
    created_by text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    ready_at timestamp with time zone,
    exited_at timestamp with time zone,
    status unweave.exec_status DEFAULT 'pending'::unweave.exec_status NOT NULL,
    project_id text NOT NULL,
    error text,
    build_id text,
    spec jsonb DEFAULT '{}'::jsonb NOT NULL,
    commit_id text,
    git_remote_url text,
    command text[],
    metadata jsonb DEFAULT '{}'::jsonb NOT NULL,
    image text DEFAULT 'ubuntu:latest'::text NOT NULL,
    provider text NOT NULL,
    CONSTRAINT session_id_check CHECK ((length(id) > 11))
);

ALTER TABLE unweave.exec OWNER TO postgres;

CREATE TABLE unweave.project (
    id text DEFAULT ('pr_'::text || public.nanoid()) NOT NULL,
    default_build_id text,
    CONSTRAINT project_id_check CHECK ((length(id) > 11))
);

ALTER TABLE unweave.project OWNER TO postgres;

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

CREATE TABLE unweave.endpoint (
    id text NOT NULL,
    name text NOT NULL,
    icon text DEFAULT '🚀'::text NOT NULL,
    project_id text NOT NULL,
    http_address text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone
);

ALTER TABLE unweave.endpoint OWNER TO postgres;

CREATE TABLE unweave.endpoint_check (
    id text NOT NULL,
    endpoint_id text NOT NULL,
    project_id text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);

ALTER TABLE unweave.endpoint_check OWNER TO postgres;

CREATE TABLE unweave.endpoint_check_step (
    id text NOT NULL,
    check_id text NOT NULL,
    eval_id text NOT NULL,
    input text,
    output text,
    assertion text
);

ALTER TABLE unweave.endpoint_check_step OWNER TO postgres;

CREATE TABLE unweave.endpoint_eval (
    endpoint_id text NOT NULL,
    eval_id text NOT NULL
);

ALTER TABLE unweave.endpoint_eval OWNER TO postgres;

CREATE TABLE unweave.endpoint_version (
    id text NOT NULL,
    endpoint_id text NOT NULL,
    exec_id text NOT NULL,
    project_id text NOT NULL,
    http_address text NOT NULL,
    primary_version boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone
);

ALTER TABLE unweave.endpoint_version OWNER TO postgres;

CREATE TABLE unweave.eval (
    id text NOT NULL,
    exec_id text NOT NULL,
    project_id text NOT NULL,
    http_address text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);

ALTER TABLE unweave.eval OWNER TO postgres;

CREATE TABLE unweave.exec_ssh_key (
    exec_id text NOT NULL,
    ssh_key_id text NOT NULL
);

ALTER TABLE unweave.exec_ssh_key OWNER TO postgres;

CREATE TABLE unweave.exec_volume (
    exec_id text NOT NULL,
    volume_id text NOT NULL,
    mount_path text NOT NULL
);

ALTER TABLE unweave.exec_volume OWNER TO postgres;

CREATE TABLE unweave.volume (
    id text DEFAULT ('vol_'::text || public.nanoid()) NOT NULL,
    size integer NOT NULL,
    name text NOT NULL,
    project_id text NOT NULL,
    provider text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone
);

ALTER TABLE unweave.volume OWNER TO postgres;

ALTER TABLE ONLY unweave.account
    ADD CONSTRAINT account_pkey PRIMARY KEY (id);

ALTER TABLE ONLY unweave.build
    ADD CONSTRAINT build_pkey PRIMARY KEY (id);

ALTER TABLE ONLY unweave.endpoint_check
    ADD CONSTRAINT endpoint_check_pkey PRIMARY KEY (id);

ALTER TABLE ONLY unweave.endpoint_check_step
    ADD CONSTRAINT endpoint_check_step_pkey PRIMARY KEY (id);

ALTER TABLE ONLY unweave.endpoint_eval
    ADD CONSTRAINT endpoint_eval_pkey PRIMARY KEY (endpoint_id, eval_id);

ALTER TABLE ONLY unweave.endpoint
    ADD CONSTRAINT endpoint_pkey PRIMARY KEY (id);

ALTER TABLE ONLY unweave.endpoint_version
    ADD CONSTRAINT endpoint_version_pkey PRIMARY KEY (id);

ALTER TABLE ONLY unweave.eval
    ADD CONSTRAINT eval_pkey PRIMARY KEY (id);

ALTER TABLE ONLY unweave.exec_ssh_key
    ADD CONSTRAINT exec_ssh_key_pkey PRIMARY KEY (exec_id, ssh_key_id);

ALTER TABLE ONLY unweave.exec_volume
    ADD CONSTRAINT exec_volume_pkey PRIMARY KEY (exec_id, volume_id, mount_path);

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

ALTER TABLE ONLY unweave.volume
    ADD CONSTRAINT volume_pkey PRIMARY KEY (id);

CREATE INDEX unweave_endpoint_name_idx ON unweave.endpoint USING btree (name);

ALTER TABLE ONLY unweave.build
    ADD CONSTRAINT build_created_by_fkey FOREIGN KEY (created_by) REFERENCES unweave.account(id);

ALTER TABLE ONLY unweave.build
    ADD CONSTRAINT build_project_id_fkey FOREIGN KEY (project_id) REFERENCES unweave.project(id);

ALTER TABLE ONLY unweave.endpoint_eval
    ADD CONSTRAINT endpoint_check_endpoint_id_fkey FOREIGN KEY (endpoint_id) REFERENCES unweave.endpoint(id);

ALTER TABLE ONLY unweave.endpoint_check
    ADD CONSTRAINT endpoint_check_endpoint_id_fkey FOREIGN KEY (endpoint_id) REFERENCES unweave.endpoint(id);

ALTER TABLE ONLY unweave.endpoint_check
    ADD CONSTRAINT endpoint_check_project_id_fkey FOREIGN KEY (project_id) REFERENCES unweave.project(id);

ALTER TABLE ONLY unweave.endpoint_check_step
    ADD CONSTRAINT endpoint_check_step_check_id_fkey FOREIGN KEY (check_id) REFERENCES unweave.endpoint_check(id);

ALTER TABLE ONLY unweave.endpoint_check_step
    ADD CONSTRAINT endpoint_check_step_eval_id_fkey FOREIGN KEY (eval_id) REFERENCES unweave.eval(id);

ALTER TABLE ONLY unweave.endpoint_eval
    ADD CONSTRAINT endpoint_eval_eval_id_fkey FOREIGN KEY (eval_id) REFERENCES unweave.eval(id);

ALTER TABLE ONLY unweave.endpoint_version
    ADD CONSTRAINT endpoint_version_endpoint_id_fkey FOREIGN KEY (endpoint_id) REFERENCES unweave.endpoint(id);

ALTER TABLE ONLY unweave.endpoint_version
    ADD CONSTRAINT endpoint_version_exec_id_fkey FOREIGN KEY (exec_id) REFERENCES unweave.exec(id);

ALTER TABLE ONLY unweave.eval
    ADD CONSTRAINT eval_exec_id_fkey FOREIGN KEY (exec_id) REFERENCES unweave.exec(id);

ALTER TABLE ONLY unweave.eval
    ADD CONSTRAINT eval_project_id_fkey FOREIGN KEY (project_id) REFERENCES unweave.project(id);

ALTER TABLE ONLY unweave.exec
    ADD CONSTRAINT exec_build_id_fkey FOREIGN KEY (build_id) REFERENCES unweave.build(id);

ALTER TABLE ONLY unweave.exec
    ADD CONSTRAINT exec_created_by_fkey FOREIGN KEY (created_by) REFERENCES unweave.account(id);

ALTER TABLE ONLY unweave.exec
    ADD CONSTRAINT exec_project_id_fkey FOREIGN KEY (project_id) REFERENCES unweave.project(id);

ALTER TABLE ONLY unweave.exec_ssh_key
    ADD CONSTRAINT exec_ssh_key_exec_id_fkey FOREIGN KEY (exec_id) REFERENCES unweave.exec(id);

ALTER TABLE ONLY unweave.exec_ssh_key
    ADD CONSTRAINT exec_ssh_key_ssh_key_id_fkey FOREIGN KEY (ssh_key_id) REFERENCES unweave.ssh_key(id);

ALTER TABLE ONLY unweave.exec_volume
    ADD CONSTRAINT exec_volume_exec_id_fkey FOREIGN KEY (exec_id) REFERENCES unweave.exec(id);

ALTER TABLE ONLY unweave.exec_volume
    ADD CONSTRAINT exec_volume_volume_id_fkey FOREIGN KEY (volume_id) REFERENCES unweave.volume(id);

ALTER TABLE ONLY unweave.project
    ADD CONSTRAINT project_default_build_id_fkey FOREIGN KEY (default_build_id) REFERENCES unweave.build(id);

ALTER TABLE ONLY unweave.ssh_key
    ADD CONSTRAINT ssh_key_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES unweave.account(id);

ALTER TABLE ONLY unweave.volume
    ADD CONSTRAINT volume_project_id_fkey FOREIGN KEY (project_id) REFERENCES unweave.project(id);

