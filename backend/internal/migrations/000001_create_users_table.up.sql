--1. Enable UUID Extension (Required for gen_random_uuid())
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

--2. Create the Auth Schema
CREATE SCHEMA IF NOT EXISTS auth_schema;

--3. Create the User Table inside that schema
CREATE TABLE IF NOT EXISTS auth_schema.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMPTZ -- Nullable until first login
);