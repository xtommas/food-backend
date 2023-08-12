CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id bigint NOT NULL REFERENCES orders ON DELETE CASCADE,
    dish_id bigint NOT NULL REFERENCES dishes ON DELETE CASCADE,
    quantity int NOT NULL,
    subtotal FLOAT NOT NULL
);