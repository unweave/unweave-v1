-- +goose Up
-- +goose StatementBegin
ALTER TABLE unweave.endpoint ADD COLUMN name text;
ALTER TABLE unweave.endpoint ADD CONSTRAINT unweave_endpoint_unique_name_by_project UNIQUE (project_id, name);
ALTER TABLE unweave.endpoint ADD CONSTRAINT unweave_endpoint_unique_http_address UNIQUE (http_address);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE unweave.endpoint DROP COLUMN name;
ALTER TABLE unweave.endpoint DROP CONSTRAINT unweave_endpoint_unique_name_by_project;
ALTER TABLE unweave.endpoint DROP CONSTRAINT unweave_endpoint_unique_http_address;

-- +goose StatementEnd
