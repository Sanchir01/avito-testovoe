-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN coin INT NOT NULL DEFAULT 1000 CHECK(coin >= 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
