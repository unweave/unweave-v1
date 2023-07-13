-- +goose Up
-- +goose StatementBegin
ALTER TABLE unweave.endpoint ADD COLUMN name text UNIQUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE unweave.endpoint DROP COLUMN name;
-- +goose StatementEnd
