-- +goose Up
-- +goose StatementBegin
ALTER TABLE songs ADD CONSTRAINT songs_group_name_fkey FOREIGN KEY (group_name) REFERENCES groups(group_name) ON DELETE RESTRICT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE songs DROP CONSTRAINT songs_group_id_fkey;
-- +goose StatementEnd
