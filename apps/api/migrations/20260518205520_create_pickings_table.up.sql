CREATE TYPE picking_status AS ENUM ('PENDING', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED');

CREATE TABLE pickings (
    id          TEXT PRIMARY KEY,
    order_id    TEXT NOT NULL REFERENCES orders(id),
    status      picking_status NOT NULL DEFAULT 'PENDING',
    note        TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE picking_items (
    id          TEXT PRIMARY KEY,
    picking_id  TEXT NOT NULL REFERENCES pickings(id),
    product_id  TEXT NOT NULL REFERENCES products(id),
    location_id TEXT NOT NULL REFERENCES locations(id),
    quantity    NUMERIC(15, 4) NOT NULL,
    picked      BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now()
);