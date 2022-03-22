CREATE TABLE IF NOT EXISTS countries (
    id uuid DEFAULT uuid_generate_v4(),
    country_name VARCHAR NOT NULL,
    iso VARCHAR(100) NOT NULL,
    gdp_percapita_ppp INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);