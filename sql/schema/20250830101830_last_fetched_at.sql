-- +goose Up
-- +goose StatementBegin
ALTER TABLE feeds RENAME COLUMN last_fetched to last_fetched_at;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE feeds RENAME COLUMN last_fetched_at to last_fetched;
-- +goose StatementEnd
