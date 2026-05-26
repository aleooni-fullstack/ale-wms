CREATE TYPE receipt_status AS ENUM ('PENDING', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED');

CREATE TABLE receipts (
    id                  TEXT PRIMARY KEY,
    purchase_order_id   TEXT NOT NULL REFERENCES purchase_orders(id),
    status              receipt_status NOT NULL DEFAULT 'PENDING',
    note                TEXT,
    created_at          TIMESTAMP NOT NULL DEFAULT now(),
    updated_at          TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE receipt_items (
    id          TEXT PRIMARY KEY,
    receipt_id  TEXT NOT NULL REFERENCES receipts(id),
    product_id  TEXT NOT NULL REFERENCES products(id),
    expected_quantity   NUMERIC(15, 4) NOT NULL,
    received_quantity   NUMERIC(15, 4),
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (receipt_id, product_id)
);