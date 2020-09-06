-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS twitter_streams
(
    id               serial PRIMARY KEY,
    twitter_id       VARCHAR(255),
    twitter_username VARCHAR(255),
    channel_id       VARCHAR(255)
);

ALTER TABLE twitter_streams
    ADD UNIQUE (twitter_id, channel_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS twitter_streams;
-- +goose StatementEnd
