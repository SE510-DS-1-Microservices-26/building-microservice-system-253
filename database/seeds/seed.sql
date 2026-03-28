-- clean everything out first so we don't get foreign key errors
TRUNCATE TABLE order_items, orders, items, item_categories RESTART IDENTITY CASCADE;

-- seed item categories
INSERT INTO item_categories (name)
VALUES ('Mains'),
       ('Drinks'),
       ('Desserts'),
       ('Snacks');

-- seed items
INSERT INTO items (category_id, name, description, price, quantity)
VALUES ((SELECT id FROM item_categories WHERE name = 'Mains'), 'Burger', 'Juicy beef burger with lettuce and tomato',
        9.99, 50),
       ((SELECT id FROM item_categories WHERE name = 'Mains'), 'Pasta Carbonara',
        'Classic Italian pasta with egg and bacon', 11.49, 30),
       ((SELECT id FROM item_categories WHERE name = 'Mains'), 'Caesar Salad', 'Romaine lettuce, croutons, parmesan',
        7.99, 40),
       ((SELECT id FROM item_categories WHERE name = 'Drinks'), 'Cola', 'Chilled Coca-Cola 0.5L', 2.99, 100),
       ((SELECT id FROM item_categories WHERE name = 'Drinks'), 'Orange Juice', 'Freshly squeezed orange juice', 3.49,
        60),
       ((SELECT id FROM item_categories WHERE name = 'Drinks'), 'Water', 'Still mineral water 0.5L', 1.49, 200),
       ((SELECT id FROM item_categories WHERE name = 'Desserts'), 'Cheesecake', 'New York style cheesecake', 5.49, 20),
       ((SELECT id FROM item_categories WHERE name = 'Desserts'), 'Brownie', 'Warm chocolate brownie', 4.49, 25),
       ((SELECT id FROM item_categories WHERE name = 'Snacks'), 'French Fries', 'Crispy golden fries', 3.99, 80),
       ((SELECT id FROM item_categories WHERE name = 'Snacks'), 'Onion Rings', 'Battered onion rings', 3.49, 60);

-- seed orders
INSERT INTO orders (user_id, status)
VALUES (1, 1),
       (2, 2),
       (1, 4);

-- seed order items
INSERT INTO order_items (order_id, item_id, quantity, unit_price)
VALUES ((SELECT id FROM orders WHERE user_id = 1 AND status = 1), (SELECT id FROM items WHERE name = 'Burger'), 2,
        9.99),
       ((SELECT id FROM orders WHERE user_id = 1 AND status = 1), (SELECT id FROM items WHERE name = 'Cola'), 2, 2.99),
       ((SELECT id FROM orders WHERE user_id = 2 AND status = 2), (SELECT id FROM items WHERE name = 'Pasta Carbonara'),
        1, 11.49),
       ((SELECT id FROM orders WHERE user_id = 2 AND status = 2), (SELECT id FROM items WHERE name = 'Orange Juice'), 1,
        3.49),
       ((SELECT id FROM orders WHERE user_id = 1 AND status = 4), (SELECT id FROM items WHERE name = 'Caesar Salad'), 1,
        7.99),
       ((SELECT id FROM orders WHERE user_id = 1 AND status = 4), (SELECT id FROM items WHERE name = 'Cheesecake'), 2,
        5.49);