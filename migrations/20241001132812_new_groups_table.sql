-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS groups (
    group_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    group_name TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS groups;
-- +goose StatementEnd
