CREATE TYPE shipping_status AS ENUM ('PENDING', 'SHIPPED', 'CANCELLED');

CREATE TABLE shippings (
    id              TEXT PRIMARY KEY,
    order_id        TEXT NOT NULL REFERENCES orders(id),
    packing_id      TEXT NOT NULL REFERENCES packings(id),
    status          shipping_status NOT NULL DEFAULT 'PENDING',
    tracking_code   VARCHAR(255),
    note            TEXT,
    created_at      TIMESTAMP NOT NULL DEFAULT now(),
    updated_at      TIMESTAMP NOT NULL DEFAULT now()
);