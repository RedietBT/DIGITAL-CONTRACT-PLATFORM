-- Migration: Add refresh_tokens table (CORRECTED)
CREATE TABLE IF NOT EXISTS auth_schema.refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Added UNIQUE here so ON CONFLICT (user_id) works
    user_id UUID NOT NULL UNIQUE REFERENCES auth_schema.users(id) ON DELETE CASCADE, 
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index for faster lookups (Optional but good)
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON auth_schema.refresh_tokens(token);