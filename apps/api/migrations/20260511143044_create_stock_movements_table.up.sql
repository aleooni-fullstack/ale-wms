CREATE TYPE movement_type AS ENUM ('IN', 'OUT', 'ADJUSTMENT');

CREATE TABLE stock_movements (
    id          TEXT PRIMARY KEY,
    product_id  TEXT NOT NULL REFERENCES products(id),
    location_id TEXT NOT NULL REFERENCES locations(id),
    type        movement_type NOT NULL,
    quantity    NUMERIC(15, 4) NOT NULL,
    note        TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT now()
);