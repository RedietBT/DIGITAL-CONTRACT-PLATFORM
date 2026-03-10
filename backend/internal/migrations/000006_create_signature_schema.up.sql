-- 1. Create the Schema
CREATE SCHEMA IF NOT EXISTS signature_schema;

-- 2. Signature Types Lookup Table
CREATE TABLE signature_schema.signature_types (
    code VARCHAR(50) PRIMARY KEY, -- e.g., 'drawn', 'typed', 'stamp', 'initial'
    name VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Insert Default Types
INSERT INTO signature_schema.signature_types (code name) VALUES
('drawn', 'Hand-drawn Signature'),
('typed', 'Typed Name Signature'),
('stamp', 'Uploaded Stamp/Seal'),
('initial', 'User Initials');

-- 4. Main Signatures Table
CREATE TABLE signature_schema.signatures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_id UUID NOT NULL,
    user_id UUID NOT NULL,
    signature_type VARCHAR(50) REFERENCES signature_schema.signature_types(code),

    -- Content
    file_url TEXT,
    vector_data JSONB,
    hash VARCHAR(255),
    
    -- Placement on PDF
    page_number INT NOT NULL DEFAULT 1,
    pos_x NUMERIC(10, 2),
    pos_y NUMERIC(10, 2),
    width NUMERIC(10, 2),
    height NUMERIC(10, 2),

    -- Styling
    render_style_id UUID,     -- Links to font/color settings

    -- Audit Trail
    signed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(45),   -- Supports IPv4 and IPv6
    device_info TEXT,

    -- Constraints
    CONSTRAINT positive_dimensions CHECK (width > 0 AND height > 0)
);


-- 5. Indexes for performance
CREATE INDEX idx_signatures_contract_id ON signature_schema.signatures(contract_id);
CREATE INDEX idx_signatures_user_id ON signature_schema.signatures(user_id);