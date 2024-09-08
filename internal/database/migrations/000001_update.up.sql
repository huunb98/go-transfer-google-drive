
-- +migrate Up

CREATE TABLE IF NOT EXISTS oauth (
    id SERIAL PRIMARY KEY,
    avatar_url VARCHAR(255),
    provider VARCHAR(50),
    email VARCHAR(255) UNIQUE,
    refresh_token TEXT,
    access_token TEXT,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);