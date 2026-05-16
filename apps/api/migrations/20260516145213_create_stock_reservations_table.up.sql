CREATE TYPE reservation_status AS ENUM ('PENDING', 'CONFIRMED', 'FULFILLED', 'RELEASED', 'CANCELLED');

CREATE TABLE stock_reservations (
    id          TEXT PRIMARY KEY,
    product_id  TEXT NOT NULL REFERENCES products(id),
    location_id TEXT NOT NULL REFERENCES locations(id),
    quantity    NUMERIC(15, 4) NOT NULL,
    status      reservation_status NOT NULL DEFAULT 'PENDING',
    reference   VARCHAR(255),
    note        TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now()
);