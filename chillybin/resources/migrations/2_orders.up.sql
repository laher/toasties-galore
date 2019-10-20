BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    customer VARCHAR(255) NOT NULL,
    ingredient VARCHAR(255) NOT NULL,
    quantity INTEGER DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE INDEX orders_customer ON orders (customer);

COMMIT;
