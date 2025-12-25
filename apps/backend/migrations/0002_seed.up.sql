-- ===== Категории (RU) =====
INSERT INTO resource_categories (name) VALUES
  ('Переговорная'),
  ('Рабочее место'),
  ('Студия'),
  ('Автомобиль'),
  ('Оборудование')
ON DUPLICATE KEY UPDATE name = name;

-- ===== Тестовые ресурсы =====
-- Чтобы seed был устойчивым, вставляем через SELECT id по имени категории
INSERT INTO resources (category_id, title, description, location)
SELECT id, 'Переговорная А', 'Небольшая переговорная до 6 человек', '1 этаж'
FROM resource_categories WHERE name = 'Переговорная';

INSERT INTO resources (category_id, title, description, location)
SELECT id, 'Переговорная Б', 'Большая переговорная до 12 человек', '2 этаж'
FROM resource_categories WHERE name = 'Переговорная';
