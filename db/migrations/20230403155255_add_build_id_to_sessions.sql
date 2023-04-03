-- +goose Up
-- +goose StatementBegin
alter table unweave.session
    rename column build to build_id;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
