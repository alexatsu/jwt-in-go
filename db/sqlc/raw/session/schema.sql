CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL UNIQUE,
    refresh_token TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);