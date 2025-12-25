CREATE TABLE IF NOT EXISTS resource_categories (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_category_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS resources (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  category_id BIGINT UNSIGNED NOT NULL,
  title VARCHAR(150) NOT NULL,
  description TEXT NULL,
  location VARCHAR(150) NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_resources_category_id (category_id),
  CONSTRAINT fk_resources_category
    FOREIGN KEY (category_id) REFERENCES resource_categories(id)
    ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO resource_categories (name) VALUES
  ('Meeting room'),
  ('Workplace'),
  ('Studio')
ON DUPLICATE KEY UPDATE name = name;

INSERT INTO resources (category_id, title, description, location)
SELECT c.id, 'Room A', 'Small meeting room for up to 6 people', '1st floor'
FROM resource_categories c WHERE c.name='Meeting room'
UNION ALL
SELECT c.id, 'Room B', 'Large meeting room for up to 12 people', '2nd floor'
FROM resource_categories c WHERE c.name='Meeting room';