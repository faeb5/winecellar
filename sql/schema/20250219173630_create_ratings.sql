-- +goose Up
-- +goose StatementBegin
CREATE TABLE ratings (
    id TEXT NOT NULL PRIMARY KEY,
    wine_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    rating TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (wine_id, user_id),
    FOREIGN KEY (wine_id) REFERENCES wines(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ratings;
-- +goose StatementEnd
