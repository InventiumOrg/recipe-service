-- Sample data for recipe-service database
-- This file contains test data for recipe table

-- Insert sample recipes
INSERT INTO recipe (name, ingredients, cost) VALUES
('Classic Omelette', ARRAY['egg', 'salt', 'pepper', 'butter'], 12000),
('Pancakes', ARRAY['flour', 'egg', 'milk', 'sugar', 'butter'], 18000),
('Tomato Pasta', ARRAY['pasta', 'tomato', 'garlic', 'olive oil', 'basil'], 25000),
('Chicken Salad', ARRAY['chicken', 'lettuce', 'tomato', 'cucumber', 'olive oil'], 22000),
('Iced Coffee', ARRAY['coffee', 'ice', 'milk', 'sugar'], 15000);
