-- +goose Up
-- +goose StatementBegin
CREATE TABLE ratings (
    id text NOT NULL PRIMARY KEY,
    wine_id text NOT NULL,
    user_id text NOT NULL,
    rating TEXT NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (wine_id, user_id),
    FOREIGN KEY (wine_id) REFERENCES wines(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ratings;
-- +goose StatementEnd
