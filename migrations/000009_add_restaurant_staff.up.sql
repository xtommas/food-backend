CREATE TABLE IF NOT EXISTS restaurant_staff (
    user_id BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
    restaurant_id BIGINT NOT NULL REFERENCES restaurants ON DELETE CASCADE,
    role TEXT NOT NULL DEFAULT 'staff',
    PRIMARY KEY (user_id, restaurant_id)
);
