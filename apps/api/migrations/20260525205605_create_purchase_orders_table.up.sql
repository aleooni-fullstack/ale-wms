CREATE TYPE purchase_order_status AS ENUM ('DRAFT', 'CONFIRMED', 'RECEIVING', 'COMPLETED', 'CANCELLED');

CREATE TABLE purchase_orders (
    id          TEXT PRIMARY KEY,
    reference   VARCHAR(255) NOT NULL UNIQUE,
    supplier    VARCHAR(255),
    status      purchase_order_status NOT NULL DEFAULT 'DRAFT',
    note        TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE purchase_order_items (
    id                  TEXT PRIMARY KEY,
    purchase_order_id   TEXT NOT NULL REFERENCES purchase_orders(id),
    product_id          TEXT NOT NULL REFERENCES products(id),
    quantity            NUMERIC(15, 4) NOT NULL,
    created_at          TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (purchase_order_id, product_id)
);