-- Only insert seed data if products table is empty
INSERT INTO products (name, price, category, thumbnail_url, mobile_url, tablet_url, desktop_url)
SELECT 'Chicken Waffle', 15.99, 'Waffle', 
       'https://example.com/images/chicken-waffle-thumb.jpg',
       'https://example.com/images/chicken-waffle-mobile.jpg',
       'https://example.com/images/chicken-waffle-tablet.jpg',
       'https://example.com/images/chicken-waffle-desktop.jpg'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE name = 'Chicken Waffle' AND price = 15.99)

UNION ALL

SELECT 'Classic Waffle', 8.99, 'Waffle',
       'https://example.com/images/classic-waffle-thumb.jpg',
       'https://example.com/images/classic-waffle-mobile.jpg',
       'https://example.com/images/classic-waffle-tablet.jpg',
       'https://example.com/images/classic-waffle-desktop.jpg'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE name = 'Classic Waffle' AND price = 8.99)

UNION ALL

SELECT 'Chocolate Waffle', 10.99, 'Waffle',
       'https://example.com/images/chocolate-waffle-thumb.jpg',
       'https://example.com/images/chocolate-waffle-mobile.jpg',
       'https://example.com/images/chocolate-waffle-tablet.jpg',
       'https://example.com/images/chocolate-waffle-desktop.jpg'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE name = 'Chocolate Waffle' AND price = 10.99)

UNION ALL

SELECT 'Berry Waffle', 12.99, 'Waffle',
       'https://example.com/images/berry-waffle-thumb.jpg',
       'https://example.com/images/berry-waffle-mobile.jpg',
       'https://example.com/images/berry-waffle-tablet.jpg',
       'https://example.com/images/berry-waffle-desktop.jpg'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE name = 'Berry Waffle' AND price = 12.99)

UNION ALL

SELECT 'Sausage Waffle', 14.99, 'Waffle',
       'https://example.com/images/sausage-waffle-thumb.jpg',
       'https://example.com/images/sausage-waffle-mobile.jpg',
       'https://example.com/images/sausage-waffle-tablet.jpg',
       'https://example.com/images/sausage-waffle-desktop.jpg'
WHERE NOT EXISTS (SELECT 1 FROM products WHERE name = 'Sausage Waffle' AND price = 14.99);