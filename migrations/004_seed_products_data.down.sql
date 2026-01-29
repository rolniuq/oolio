-- Remove seed data (only if it matches our seed data)
DELETE FROM products 
WHERE name IN ('Chicken Waffle', 'Classic Waffle', 'Chocolate Waffle', 'Berry Waffle', 'Sausage Waffle')
AND price IN (15.99, 8.99, 10.99, 12.99, 14.99)
AND category = 'Waffle';