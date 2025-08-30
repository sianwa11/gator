-- +goose Up
-- +goose StatementBegin
ALTER TABLE posts
ADD CONSTRAINT posts_url_key UNIQUE (url);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE posts
DROP CONSTRAINT posts_url_key;
-- +goose StatementEnd
