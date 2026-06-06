-- =============================================================================
-- Drop triggers first (before removing columns they reference)
-- =============================================================================
DROP TRIGGER IF EXISTS orders_set_updated_at ON orders;
DROP TRIGGER IF EXISTS dishes_set_updated_at ON dishes;

-- =============================================================================
-- Drop trigger function
-- =============================================================================
DROP FUNCTION IF EXISTS set_updated_at();

-- =============================================================================
-- Drop indexes
-- =============================================================================
DROP INDEX IF EXISTS order_items_order_id_idx;
DROP INDEX IF EXISTS dishes_restaurant_id_idx;
DROP INDEX IF EXISTS orders_status_idx;
DROP INDEX IF EXISTS orders_restaurant_id_idx;
DROP INDEX IF EXISTS orders_user_id_idx;

-- =============================================================================
-- Drop orders.status CHECK constraint
-- =============================================================================
ALTER TABLE orders
    DROP CONSTRAINT IF EXISTS orders_status_check;

-- =============================================================================
-- Drop updated_at columns
-- =============================================================================
ALTER TABLE orders DROP COLUMN IF EXISTS updated_at;
ALTER TABLE dishes DROP COLUMN IF EXISTS updated_at;

-- =============================================================================
-- Drop snapshot columns from order_items
-- =============================================================================
ALTER TABLE order_items
    DROP COLUMN IF EXISTS unit_price,
    DROP COLUMN IF EXISTS dish_name;

-- =============================================================================
-- BIGINT → FLOAT for money columns (convert cents back to decimal)
-- =============================================================================
ALTER TABLE order_items
    ALTER COLUMN subtotal TYPE FLOAT USING (subtotal / 100.0)::FLOAT;

ALTER TABLE orders
    ALTER COLUMN total TYPE FLOAT USING (total / 100.0)::FLOAT;

ALTER TABLE dishes
    ALTER COLUMN price TYPE FLOAT USING (price / 100.0)::FLOAT;

-- =============================================================================
-- Re-point orders.restaurant_id back to users
-- =============================================================================
ALTER TABLE orders
    DROP CONSTRAINT orders_restaurant_id_fkey;

ALTER TABLE orders
    ADD CONSTRAINT orders_restaurant_id_fkey
        FOREIGN KEY (restaurant_id) REFERENCES users (id) ON DELETE CASCADE;

-- =============================================================================
-- Re-point dishes.restaurant_id back to users
-- =============================================================================
ALTER TABLE dishes
    DROP CONSTRAINT dishes_restaurant_id_fkey;

ALTER TABLE dishes
    ADD CONSTRAINT dishes_restaurant_id_fkey
        FOREIGN KEY (restaurant_id) REFERENCES users (id) ON DELETE CASCADE;
