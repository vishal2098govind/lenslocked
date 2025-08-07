-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS
    "password_resets" (
        id SERIAL PRIMARY KEY,
        user_id INT UNIQUE,
        token_hash TEXT UNIQUE NOT NULL, -- a hash not actual token, similar reasons as password
        expires_at TIMESTAMPTZ NOT NULL
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "password_resets";
-- +goose StatementEnd
