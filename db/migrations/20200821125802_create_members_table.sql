-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS members
(
    id       serial PRIMARY KEY,
    username VARCHAR(255),
    email    VARCHAR(255)
);

ALTER TABLE members
    ADD UNIQUE (username, email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS members;
-- +goose StatementEnd
