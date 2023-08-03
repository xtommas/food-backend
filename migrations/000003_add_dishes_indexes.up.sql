CREATE INDEX IF NOT EXISTS dishes_names_idx ON dishes USING GIN (to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS dishes_categories_idx ON dishes USING GIN (category);