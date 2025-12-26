INSERT INTO resources (category_id, title, description, location, price_per_hour, is_active)
SELECT c.id, x.title, x.description, x.location, x.price_per_hour, TRUE
FROM resource_categories c
JOIN (
  SELECT 'Переговорная' AS cat, 'Переговорная А' AS title,
         'Небольшая переговорная до 6 человек' AS description,
         '1 этаж' AS location, 500 AS price_per_hour
  UNION ALL
  SELECT 'Переговорная', 'Переговорная Б',
         'Большая переговорная до 12 человек',
         '2 этаж', 800
  UNION ALL
  SELECT 'Рабочее место', 'Коворкинг у окна',
         'Рабочее место, розетки рядом',
         'Open-space', 200
  UNION ALL
  SELECT 'Студия', 'Студия звукозаписи',
         'Шумоизоляция, базовый комплект',
         'Блок S', 1200
  UNION ALL
  SELECT 'Оборудование', 'Проектор Epson',
         'FullHD проектор для презентаций',
         'Склад', 150
) x ON x.cat = c.name
ON DUPLICATE KEY UPDATE
  description = VALUES(description),
  location = VALUES(location),
  price_per_hour = VALUES(price_per_hour),
  is_active = VALUES(is_active);