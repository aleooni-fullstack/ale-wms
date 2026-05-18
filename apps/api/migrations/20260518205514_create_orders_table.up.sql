CREATE TYPE order_status AS ENUM ('DRAFT', 'CONFIRMED', 'PICKING', 'PACKING', 'SHIPPED', 'CANCELLED');

CREATE TABLE orders (
    id          TEXT PRIMARY KEY,
    reference   VARCHAR(255) NOT NULL UNIQUE,
    status      order_status NOT NULL DEFAULT 'DRAFT',
    note        TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE order_items (
    id          TEXT PRIMARY KEY,
    order_id    TEXT NOT NULL REFERENCES orders(id),
    product_id  TEXT NOT NULL REFERENCES products(id),
    location_id TEXT NOT NULL REFERENCES locations(id),
    quantity    NUMERIC(15, 4) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (order_id, product_id, location_id)
);