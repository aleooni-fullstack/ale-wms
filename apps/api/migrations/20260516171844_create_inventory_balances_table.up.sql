CREATE TYPE inventory_balance_status AS ENUM ('DRAFT', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED');

CREATE TABLE inventory_balances (
    id           TEXT PRIMARY KEY,
    location_id  TEXT NOT NULL REFERENCES locations(id),
    status       inventory_balance_status NOT NULL DEFAULT 'DRAFT',
    note         TEXT,
    created_at   TIMESTAMP NOT NULL DEFAULT now(),
    updated_at   TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE inventory_balance_items (
    id                   TEXT PRIMARY KEY,
    inventory_balance_id TEXT NOT NULL REFERENCES inventory_balances(id),
    product_id           TEXT NOT NULL REFERENCES products(id),
    system_quantity      NUMERIC(15, 4) NOT NULL,
    counted_quantity     NUMERIC(15, 4),
    created_at           TIMESTAMP NOT NULL DEFAULT now(),
    updated_at           TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (inventory_balance_id, product_id)
);