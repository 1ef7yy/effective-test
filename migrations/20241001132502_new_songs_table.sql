-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS songs (
    song_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    group_id UUID NOT NULL,
    song_name TEXT NOT NULL,
    release_date DATE NOT NULL,
    song_text TEXT NOT NULL,
    link TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS songs;
-- +goose StatementEnd
