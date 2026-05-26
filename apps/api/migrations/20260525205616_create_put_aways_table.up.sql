CREATE TYPE put_away_status AS ENUM ('PENDING', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED');

CREATE TABLE put_aways (
    id          TEXT PRIMARY KEY,
    receipt_id  TEXT NOT NULL REFERENCES receipts(id),
    status      put_away_status NOT NULL DEFAULT 'PENDING',
    note        TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE put_away_items (
    id              TEXT PRIMARY KEY,
    put_away_id     TEXT NOT NULL REFERENCES put_aways(id),
    product_id      TEXT NOT NULL REFERENCES products(id),
    location_id     TEXT NOT NULL REFERENCES locations(id),
    quantity        NUMERIC(15, 4) NOT NULL,
    put_away        BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMP NOT NULL DEFAULT now(),
    updated_at      TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (put_away_id, product_id)
);