INSERT INTO resource_categories (name) VALUES
  ('Meeting room'),
  ('Workplace'),
  ('Studio'),
  ('Car'),
  ('Equipment')
ON DUPLICATE KEY UPDATE name = name;