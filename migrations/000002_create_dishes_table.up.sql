CREATE TABLE IF NOT EXISTS dishes (
    id BIGSERIAL PRIMARY KEY,
    restaurant_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    name TEXT NOT NULL,
    price FLOAT NOT NULL,
    description TEXT NOT NULL,
    categories TEXT[] NOT NULL,
    photo TEXT,
    available BOOLEAN NOT NULL DEFAULT TRUE
);