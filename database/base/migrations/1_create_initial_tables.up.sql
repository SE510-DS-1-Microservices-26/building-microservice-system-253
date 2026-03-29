-- create a table for item categories
CREATE TABLE IF NOT EXISTS item_categories
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- create table for items
CREATE TABLE IF NOT EXISTS items
(
    id          SERIAL PRIMARY KEY,
    category_id INT            REFERENCES item_categories (id) ON DELETE SET NULL,
    name        VARCHAR(255)   NOT NULL,
    description TEXT           NOT NULL DEFAULT '',
    image_url   TEXT           NOT NULL DEFAULT '',
    price       NUMERIC(10, 2) NOT NULL,
    quantity    INT            NOT NULL DEFAULT 0,
    created_at  TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- create index for items' category, price and quantity
CREATE INDEX IF NOT EXISTS idx_items_category_id ON items (category_id);
CREATE INDEX IF NOT EXISTS idx_items_price ON items (price);
CREATE INDEX IF NOT EXISTS idx_items_quantity ON items (quantity);

-- create table for orders
CREATE TABLE IF NOT EXISTS orders
(
    id         SERIAL PRIMARY KEY,
    user_id    INT,
    status     SMALLINT  NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- order statuses comment
COMMENT ON COLUMN orders.status IS '1=pending, 2=confirmed, 3=ready, 4=delivered, 5=cancelled';

-- create indexes for orders' user and status
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders (user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders (status);

-- create a table for order items
CREATE TABLE IF NOT EXISTS order_items
(
    order_id   INT            NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    item_id    INT            NOT NULL REFERENCES items (id) ON DELETE CASCADE,
    quantity   INT            NOT NULL,
    unit_price NUMERIC(10, 2) NOT NULL,
    PRIMARY KEY (order_id, item_id)
);