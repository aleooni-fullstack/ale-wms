CREATE TABLE warehouses (
    id         TEXT PRIMARY KEY,
    code       VARCHAR(50) NOT NULL UNIQUE,
    name       VARCHAR(255) NOT NULL,
    address    TEXT,
    active     BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);