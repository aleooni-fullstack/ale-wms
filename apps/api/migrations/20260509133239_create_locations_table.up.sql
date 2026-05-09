CREATE TABLE locations (
    id           TEXT PRIMARY KEY,
    zone_id      TEXT NOT NULL REFERENCES zones(id),
    code         VARCHAR(50) NOT NULL,
    name         VARCHAR(255) NOT NULL,
    active       BOOLEAN NOT NULL DEFAULT true,
    created_at   TIMESTAMP NOT NULL DEFAULT now(),
    updated_at   TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (zone_id, code)
);