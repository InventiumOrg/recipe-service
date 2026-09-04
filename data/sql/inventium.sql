-- Sample data for recipe-service database
-- This file contains test data for recipe table

-- Insert sample recipes
INSERT INTO recipe (name, ingredients, cost) VALUES
('Classic Omelette', '[{"name":"egg","quantity":2,"measure":"pcs","unit":"egg"},{"name":"butter","quantity":10,"measure":"g","unit":"slice"},{"name":"salt","quantity":1,"measure":"pinch","unit":"pinch"}]'::jsonb, 12000),
('Pancakes', '[{"name":"flour","quantity":200,"measure":"g","unit":"bag"},{"name":"milk","quantity":250,"measure":"ml","unit":"cup"},{"name":"egg","quantity":1,"measure":"pcs","unit":"egg"}]'::jsonb, 18000),
('Tomato Pasta', '[{"name":"pasta","quantity":200,"measure":"g","unit":"pack"},{"name":"tomato","quantity":2,"measure":"pcs","unit":"tomato"},{"name":"garlic","quantity":3,"measure":"clove","unit":"clove"}]'::jsonb, 25000),
('Chicken Salad', '[{"name":"chicken","quantity":150,"measure":"g","unit":"breast"},{"name":"lettuce","quantity":100,"measure":"g","unit":"bunch"},{"name":"olive oil","quantity":1,"measure":"tbsp","unit":"spoon"}]'::jsonb, 22000),
('Iced Coffee', '[{"name":"coffee","quantity":18,"measure":"g","unit":"shot"},{"name":"milk","quantity":200,"measure":"ml","unit":"cup"},{"name":"ice","quantity":6,"measure":"pcs","unit":"cube"}]'::jsonb, 15000);
