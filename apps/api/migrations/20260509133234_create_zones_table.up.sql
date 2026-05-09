CREATE TABLE zones (
    id           TEXT PRIMARY KEY,
    warehouse_id TEXT NOT NULL REFERENCES warehouses(id),
    code         VARCHAR(50) NOT NULL,
    name         VARCHAR(255) NOT NULL,
    active       BOOLEAN NOT NULL DEFAULT true,
    created_at   TIMESTAMP NOT NULL DEFAULT now(),
    updated_at   TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (warehouse_id, code)
);