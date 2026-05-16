CREATE TABLE IF NOT EXISTS restaurants (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    photo TEXT,
    address TEXT NOT NULL,
    city TEXT NOT NULL,
    state TEXT,
    province TEXT,
    country TEXT NOT NULL,
    latitude NUMERIC(9, 6),
    longitude NUMERIC(9, 6),
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1
);
