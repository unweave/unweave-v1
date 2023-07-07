-- +goose Up
-- +goose StatementBegin
CREATE TABLE unweave.eval (
    id text NOT NULL PRIMARY KEY,
    exec_id text NOT NULL REFERENCES unweave.exec (id),
    project_id text NOT NULL REFERENCES unweave.project (id),
    created_at timestamp with time zone DEFAULT now() NOT NULL
);

CREATE TABLE unweave.endpoint (
    id text NOT NULL PRIMARY KEY,
    exec_id text NOT NULL REFERENCES unweave.exec (id),
    project_id text NOT NULL REFERENCES unweave.project (id),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone
);

CREATE TABLE unweave.endpoint_eval (
    endpoint_id text NOT NULL REFERENCES unweave.endpoint (id),
    eval_id text NOT NULL REFERENCES unweave.eval (id),
    PRIMARY KEY (endpoint_id, eval_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE unweave.eval;
DROP TABLE unweave.endpoint;
DROP TABLE unweave.endpoint_eval;
-- +goose StatementEnd
