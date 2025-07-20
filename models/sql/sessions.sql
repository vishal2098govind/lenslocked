CREATE TABLE IF NOT EXISTS
    sessions (
        id SERIAL PRIMARY KEY,
        user_id INT UNIQUE,
        token_hash TEXT UNIQUE NOT NULL -- a hash not actual token, similar reasons as password
    );


ALTER TABLE sessions
ADD CONSTRAINT sessions_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;