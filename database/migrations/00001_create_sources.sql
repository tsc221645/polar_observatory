-- +goose Up
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE SCHEMA polar;

CREATE TABLE polar.sources (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    slug text NOT NULL UNIQUE,
    provider text NOT NULL,
    product text NOT NULL,
    version text NOT NULL,
    citation text NOT NULL,
    terms_url text NOT NULL,
    redistribution_status text NOT NULL,
    expected_interval interval NOT NULL CHECK (expected_interval > INTERVAL '0')
);

-- +goose Down
DROP TABLE polar.sources;
DROP SCHEMA polar;
