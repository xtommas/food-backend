-- =============================================================================
-- Re-point dishes.restaurant_id → restaurants (was → users)
-- =============================================================================
ALTER TABLE dishes
    DROP CONSTRAINT dishes_restaurant_id_fkey;

ALTER TABLE dishes
    ADD CONSTRAINT dishes_restaurant_id_fkey
        FOREIGN KEY (restaurant_id) REFERENCES restaurants (id) ON DELETE CASCADE;

-- =============================================================================
-- Re-point orders.restaurant_id → restaurants (was → users)
-- =============================================================================
ALTER TABLE orders
    DROP CONSTRAINT orders_restaurant_id_fkey;

ALTER TABLE orders
    ADD CONSTRAINT orders_restaurant_id_fkey
        FOREIGN KEY (restaurant_id) REFERENCES restaurants (id) ON DELETE CASCADE;

-- =============================================================================
-- FLOAT → BIGINT for all money columns (store as cents)
-- =============================================================================
ALTER TABLE dishes      ALTER COLUMN price    TYPE bigint USING (price * 100)::bigint;
ALTER TABLE orders      ALTER COLUMN total    TYPE bigint USING (total * 100)::bigint;
ALTER TABLE order_items ALTER COLUMN subtotal TYPE bigint USING (subtotal * 100)::bigint;

-- =============================================================================
-- Snapshot dish name + unit price in order_items (preserve order history)
--    New rows will have these populated by the application.
--    Existing rows get a best-effort backfill from the dishes table.
--    Note: dishes.price is already in cents at this point.
-- =============================================================================
ALTER TABLE order_items
    ADD COLUMN IF NOT EXISTS dish_name  TEXT,
    ADD COLUMN IF NOT EXISTS unit_price bigint;

UPDATE order_items oi
SET
    dish_name  = d.name,
    unit_price = d.price
FROM dishes d
WHERE oi.dish_id = d.id;

-- Now that backfill is done, enforce NOT NULL going forward
ALTER TABLE order_items
    ALTER COLUMN dish_name  SET NOT NULL,
    ALTER COLUMN unit_price SET NOT NULL;

-- =============================================================================
-- Add updated_at to dishes and orders
-- =============================================================================
ALTER TABLE dishes ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW();
ALTER TABLE orders ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW();

-- =============================================================================
-- orders.status CHECK constraint
-- =============================================================================
ALTER TABLE orders
    ADD CONSTRAINT orders_status_check
        CHECK (status IN ('pending', 'confirmed', 'preparing', 'ready', 'delivered', 'cancelled'));

-- =============================================================================
-- Indexes for common query patterns
-- =============================================================================

-- Orders filtered by user or restaurant
CREATE INDEX IF NOT EXISTS orders_user_id_idx       ON orders (user_id);
CREATE INDEX IF NOT EXISTS orders_restaurant_id_idx ON orders (restaurant_id);

-- Orders filtered by status (e.g. "show all pending orders for restaurant X")
CREATE INDEX IF NOT EXISTS orders_status_idx ON orders (status);

-- Dishes filtered by restaurant
CREATE INDEX IF NOT EXISTS dishes_restaurant_id_idx ON dishes (restaurant_id);

-- Order items looked up by order (most common join path)
CREATE INDEX IF NOT EXISTS order_items_order_id_idx ON order_items (order_id);

-- =============================================================================
-- Shared trigger function (reusable across any table)
-- =============================================================================
CREATE OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- Attach trigger to dishes and orders
-- =============================================================================
CREATE TRIGGER dishes_set_updated_at
    BEFORE UPDATE ON dishes
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER orders_set_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
