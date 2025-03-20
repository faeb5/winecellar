-- +goose Up
-- +goose StatementBegin
CREATE TABLE wines (
    id TEXT NOT NULL PRIMARY KEY,
    color TEXT NOT NULL,
    name TEXT NOT NULL,
    producer TEXT NOT NULL,
    country TEXT NOT NULL,
    vintage INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (producer, name, vintage)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE wines;
-- +goose StatementEnd
