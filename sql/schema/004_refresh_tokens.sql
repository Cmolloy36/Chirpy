-- +goose Up
CREATE TABLE refresh_tokens(
    token VARCHAR(64) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULl,
    revoked_at TIMESTAMP DEFAULT(NULL)
);

-- +goose Down
DROP TABLE refresh_tokens;