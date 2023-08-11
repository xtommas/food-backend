ALTER TABLE dishes ADD CONSTRAINT dishes_price_check CHECK (price >= 0);

ALTER TABLE dishes ADD CONSTRAINT categories_length_check CHECK (array_length(categories, 1) BETWEEN 1 AND 5);