-- +goose Up
-- +goose StatementBegin
ALTER TABLE unweave.volume ADD COLUMN size integer NOT NULL DEFAULT 0;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd
