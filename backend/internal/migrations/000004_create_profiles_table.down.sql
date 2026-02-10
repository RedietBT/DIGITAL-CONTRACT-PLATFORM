-- 1. Drop the trigger first
DROP TRIGGER IF EXISTS update_profiles_modtime ON profile_schema.profiles;

-- 2. Drop the function
DROP FUNCTION IF EXISTS profile_schema.update_updated_at_column();

-- 3. Drop the profiles table (This also removes the FK constraint)
DROP TABLE IF EXISTS profile_schema.profiles;

-- 4. Finally, drop the schema
DROP SCHEMA IF EXISTS profile_schema;