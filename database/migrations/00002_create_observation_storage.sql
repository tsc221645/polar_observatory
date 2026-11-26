-- +goose Up
CREATE TABLE polar.artifacts (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    source_id bigint NOT NULL REFERENCES polar.sources (id),
    url text NOT NULL,
    retrieved_at timestamptz NOT NULL,
    sha256 text NOT NULL CHECK (sha256 ~ '^[0-9a-f]{64}$'),
    path text NOT NULL,
    bytes bigint NOT NULL CHECK (bytes >= 0),
    etag text,
    last_modified text,
    UNIQUE (source_id, sha256)
);

CREATE TABLE polar.ingestion_runs (
    id uuid PRIMARY KEY,
    source_id bigint NOT NULL REFERENCES polar.sources (id),
    status text NOT NULL CHECK (status IN ('running', 'succeeded', 'failed')),
    started_at timestamptz NOT NULL,
    finished_at timestamptz,
    accepted_count bigint NOT NULL DEFAULT 0 CHECK (accepted_count >= 0),
    rejected_count bigint NOT NULL DEFAULT 0 CHECK (rejected_count >= 0),
    error_code text,
    git_sha text
);

CREATE TABLE polar.observations (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    source_id bigint NOT NULL REFERENCES polar.sources (id),
    artifact_id bigint NOT NULL REFERENCES polar.artifacts (id),
    entity_key text NOT NULL,
    metric text NOT NULL,
    observed_at timestamptz NOT NULL,
    value double precision,
    unit text NOT NULL,
    original_value double precision,
    original_unit text,
    quality_flag text NOT NULL DEFAULT 'accepted',
    revision integer NOT NULL CHECK (revision > 0),
    ingested_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (source_id, entity_key, metric, observed_at, revision)
);

CREATE INDEX observations_entity_metric_time_idx
    ON polar.observations (entity_key, metric, observed_at, id);

CREATE INDEX observations_source_time_idx
    ON polar.observations (source_id, observed_at);

-- +goose Down
DROP TABLE polar.observations;
DROP TABLE polar.ingestion_runs;
DROP TABLE polar.artifacts;
