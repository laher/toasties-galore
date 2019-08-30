set timezone='UTC';
CREATE TABLE IF NOT EXISTS ingredients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT NULL,
    quantity INTEGER DEFAULT 0 NOT NULL
);

CREATE UNIQUE INDEX ingredients_name ON ingredients (name);
