-- +goose Up
-- +goose StatementBegin
CREATE TABLE unweave.eval (
    id text NOT NULL PRIMARY KEY,
    exec_id text NOT NULL,
    project_id text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);

CREATE TABLE unweave.endpoint (
    id text NOT NULL PRIMARY KEY,
    exec_id text NOT NULL,
    project_id text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone
);

CREATE TABLE unweave.endpoint_eval (
    endpoint_id text NOT NULL,
    eval_id text NOT NULL,
    PRIMARY KEY (endpoint_id, eval_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE unweave.eval;
DROP TABLE unweave.endpoint;
DROP TABLE unweave.endpoint_eval;
-- +goose StatementEnd
