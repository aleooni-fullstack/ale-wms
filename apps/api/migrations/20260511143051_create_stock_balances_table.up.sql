CREATE TABLE stock_balances (
    id          TEXT PRIMARY KEY,
    product_id  TEXT NOT NULL REFERENCES products(id),
    location_id TEXT NOT NULL REFERENCES locations(id),
    quantity    NUMERIC(15, 4) NOT NULL DEFAULT 0,
    updated_at  TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (product_id, location_id)
);