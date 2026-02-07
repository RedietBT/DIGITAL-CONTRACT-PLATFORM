-- 1. Create the Profile Schema
CREATE SCHEMA IF NOT EXISTS profile_schema;

-- 2. Create the Profiles Table inside that schema
CREATE TABLE IF NOT EXISTS profile_schema.profiles (
    -- Linked to Auth Schema User ID
    user_id            UUID PRIMARY KEY,
    display_name       VARCHAR(100) NOT NULL,
    bio                TEXT,
    skill_level        INT DEFAULT 1,
    is_template_seller BOOLEAN DEFAULT FALSE,
    rating_avg         NUMERIC(3, 2) DEFAULT 0.00,
    rating_count       INT DEFAULT 0,
    created_at         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    -- 3. The Cross-Schema Foreign Key
    CONSTRAINT fk_user
      FOREIGN KEY(user_id) 
      REFERENCES auth_schema.users(id) 
      ON DELETE CASCADE
);

-- 4. Updated_at Trigger for this specific schema/table
CREATE OR REPLACE FUNCTION profile_schema.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_profiles_modtime
    BEFORE UPDATE ON profile_schema.profiles
    FOR EACH ROW
    EXECUTE PROCEDURE profile_schema.update_updated_at_column();