CREATE TYPE transfer_status AS ENUM ('PENDING', 'COMPLETED', 'CANCELLED');

CREATE TABLE stock_transfers (
    id                  TEXT PRIMARY KEY,
    product_id          TEXT NOT NULL REFERENCES products(id),
    from_location_id    TEXT NOT NULL REFERENCES locations(id),
    to_location_id      TEXT NOT NULL REFERENCES locations(id),
    quantity            NUMERIC(15, 4) NOT NULL,
    status              transfer_status NOT NULL DEFAULT 'PENDING',
    note                TEXT,
    created_at          TIMESTAMP NOT NULL DEFAULT now(),
    updated_at          TIMESTAMP NOT NULL DEFAULT now()
);