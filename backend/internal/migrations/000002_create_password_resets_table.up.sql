CREATE TABLE IF NOT EXISTS auth_schema.password_resets (
    email VARCHAR(255) PRIMARY KEY REFERENCES auth_schema.users(email) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL
);