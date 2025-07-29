-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE rates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    currency VARCHAR(10) NOT NULL,
    ask DECIMAL NOT NULL,
    bid DECIMAL NOT NULL,
    timestamp BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT rates_currency_unique UNIQUE (currency)
);

CREATE INDEX idx_rates_currency ON rates(currency);
CREATE INDEX idx_rates_timestamp ON rates(timestamp);

-- +goose Down
DROP TABLE IF EXISTS rates;
DROP EXTENSION IF EXISTS "uuid-ossp";