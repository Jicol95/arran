CREATE TABLE IF NOT EXISTS processed_event (
    id              UUID PRIMARY KEY,
    source          TEXT NOT NULL,
    processed_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
