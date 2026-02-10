CREATE SCHEMA IF NOT EXISTS contract_schema;

-- 1. Main Contract Table
CREATE TABLE contract_schema.contracts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_user_id UUID NOT NULL,          -- User who created and manages the contract
    title VARCHAR(255) NOT NULL,          -- Display name of the document
    description TEXT,                     -- Internal notes or summary
    status VARCHAR(50) DEFAULT 'draft',    -- Lifecycle state: draft, pending, active, completed, cancelled
    current_version INT DEFAULT 1,         -- Tracks the latest version number
    template_id UUID,                      -- Reference to a template if used
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ               -- Timestamp when all required signatures are fulfilled
);

COMMENT ON TABLE contract_schema.contracts IS 'Core metadata for digital contracts';
COMMENT ON COLUMN contract_schema.contracts.status IS 'Determines if the contract is editable (draft) or locked for signing (pending)';

-- 2. Versioning Table
CREATE TABLE contract_schema.contract_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_id UUID REFERENCES contract_schema.contracts(id) ON DELETE CASCADE,
    version_number INT NOT NULL,           -- Sequential versioning (1, 2, 3...)
    content_json JSONB NOT NULL,           -- The actual structured text/data of the contract
    created_by UUID NOT NULL,              -- The user who pushed this specific update
    created_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE contract_schema.contract_versions IS 'Immutable history of contract content changes';

-- 3. Participants Table
CREATE TABLE contract_schema.contract_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_id UUID REFERENCES contract_schema.contracts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,                 -- External reference to auth_schema.users
    role VARCHAR(50) NOT NULL,             -- User role: signer, viewer, approver
    signing_order INT DEFAULT 0,           -- Order in which this user is notified to sign
    is_required BOOLEAN DEFAULT TRUE,      -- If false, the user is an optional signer
    status VARCHAR(50) DEFAULT 'pending',  -- Individual status: pending, signed, declined
    signed_at TIMESTAMPTZ,                 -- When this specific user completed their requirement
    added_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(contract_id, user_id)           -- Prevents adding the same user twice to one contract
);

COMMENT ON TABLE contract_schema.contract_participants IS 'Tracks the involvement and progress of users within a contract';

-- 4. Signature Requirements Table
CREATE TABLE contract_schema.contract_signature_requirements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_id UUID NOT NULL REFERENCES contract_schema.contracts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES auth_schema.users(id) ON DELETE CASCADE,
    required_types TEXT[] NOT NULL,        -- Types allowed: e.g., {'handwritten', 'biometric', 'otp'}
    min_required INT DEFAULT 1,            -- Minimum number of signatures from the types above
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(contract_id, user_id)
);

COMMENT ON TABLE contract_schema.contract_signature_requirements IS 'Specific legal requirements for how a user must sign a specific contract';

-- Index for quick lookups when validating a signature attempt
CREATE INDEX IF NOT EXISTS idx_signature_req_user ON contract_schema.contract_signature_requirements(user_id);