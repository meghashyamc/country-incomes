CREATE TABLE IF NOT EXISTS country_parity_factors (
    id uuid DEFAULT uuid_generate_v4(),
    country_from_iso VARCHAR(100) NOT NULL,
    country_to_iso VARCHAR(100) NOT NULL,
    parity_factor DECIMAL NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);