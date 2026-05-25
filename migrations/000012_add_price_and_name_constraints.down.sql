ALTER TABLE dishes DROP CONSTRAINT dishes_price_check;
ALTER TABLE dishes ADD CONSTRAINT dishes_price_check CHECK (price >= 0);
ALTER TABLE dishes DROP CONSTRAINT dishes_name_check;
