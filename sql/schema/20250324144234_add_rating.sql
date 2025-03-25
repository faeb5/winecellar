-- +goose Up
-- +goose StatementBegin
ALTER table ratings
ADD COLUMN rating TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER table ratings
DROP COLUMN rating;
-- +goose StatementEnd
