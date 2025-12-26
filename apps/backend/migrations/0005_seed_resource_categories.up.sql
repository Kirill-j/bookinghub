INSERT INTO resource_categories (name) VALUES
  ('Переговорная'),
  ('Рабочее место'),
  ('Студия'),
  ('Автомобиль'),
  ('Оборудование')
ON DUPLICATE KEY UPDATE name = VALUES(name);