CREATE TABLE IF NOT EXISTS dishes (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price FLOAT NOT NULL,
    description TEXT NOT NULL,
    category TEXT[] NOT NULL,
    photo TEXT,
    available BOOLEAN NOT NULL DEFAULT TRUE
);