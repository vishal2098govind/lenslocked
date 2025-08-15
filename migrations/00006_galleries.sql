-- +goose Up
-- +goose StatementBegin
ALTER TABLE galleries DROP CONSTRAINT IF EXISTS galleries_user_id_fkey;
ALTER TABLE galleries
ADD CONSTRAINT galleries_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE galleries DROP CONSTRAINT IF EXISTS galleries_user_id_fkey;
ALTER TABLE galleries
ADD CONSTRAINT galleries_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);
-- +goose StatementEnd