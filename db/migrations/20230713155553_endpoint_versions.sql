-- +goose Up
-- +goose StatementBegin
ALTER TABLE unweave.endpoint_eval DROP CONSTRAINT endpoint_eval_endpoint_id_fkey;
ALTER TABLE unweave.endpoint_check DROP CONSTRAINT endpoint_check_endpoint_id_fkey;

ALTER TABLE unweave.endpoint_check_step DROP CONSTRAINT endpoint_check_step_check_id_fkey;

TRUNCATE unweave.endpoint_eval;
TRUNCATE unweave.endpoint_check;

ALTER TABLE unweave.endpoint_check_step 
ADD CONSTRAINT endpoint_check_step_check_id_fkey 
FOREIGN KEY (check_id) 
REFERENCES unweave.endpoint_check (id);

DROP TABLE unweave.endpoint;

CREATE TABLE unweave.endpoint (
    id text NOT NULL PRIMARY KEY,
    name text NOT NULL,
    icon text NOT NULL DEFAULT '🚀',
    project_id text NOT NULL,
    http_address text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone
);

CREATE INDEX unweave_endpoint_name_idx ON unweave.endpoint USING btree (name);

CREATE TABLE unweave.endpoint_version (
    id text NOT NULL PRIMARY KEY,
    endpoint_id text NOT NULL REFERENCES unweave.endpoint (id),
    exec_id text NOT NULL REFERENCES unweave.exec (id),
    project_id text NOT NULL,
    http_address text NOT NULL,
    primary_version boolean NOT NULL DEFAULT FALSE,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone
);

ALTER TABLE unweave.endpoint_eval 
ADD CONSTRAINT endpoint_check_endpoint_id_fkey 
FOREIGN KEY (endpoint_id) 
REFERENCES unweave.endpoint (id);

ALTER TABLE unweave.endpoint_check 
ADD CONSTRAINT endpoint_check_endpoint_id_fkey
FOREIGN KEY (endpoint_id) REFERENCES unweave.endpoint (id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE unweave.endpoint_version;
DROP TABLE unweave.endpoint;

CREATE TABLE unweave.endpoint (
    id text NOT NULL PRIMARY KEY,
    exec_id text NOT NULL REFERENCES unweave.exec (id),
    project_id text NOT NULL,
    http_address text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    name text
);

CREATE INDEX unweave_endpoint_name_idx ON unweave.endpoint USING btree (name);
-- +goose StatementEnd
