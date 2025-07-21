-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS
    "sessions" (
        id SERIAL PRIMARY KEY,
        user_id INT UNIQUE,
        token_hash TEXT UNIQUE NOT NULL -- a hash not actual token, similar reasons as password
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "sessions";
-- +goose StatementEnd
