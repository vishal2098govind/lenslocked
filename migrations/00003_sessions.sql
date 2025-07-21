-- +goose Up
-- +goose StatementBegin
ALTER TABLE sessions
ADD CONSTRAINT sessions_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sessions
DROP CONSTRAINT sessions_user_id_fkey;
-- +goose StatementEnd
