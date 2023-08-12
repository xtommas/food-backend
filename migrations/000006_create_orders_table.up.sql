CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    restaurant_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    total FLOAT NOT NULL,
    address TEXT NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    status TEXT NOT NULL
);