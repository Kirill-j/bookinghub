CREATE TABLE IF NOT EXISTS resources (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  owner_user_id BIGINT UNSIGNED NOT NULL,
  category_id BIGINT UNSIGNED NOT NULL,
  title VARCHAR(150) NOT NULL,
  description TEXT NULL,
  location VARCHAR(150) NULL,
  price_per_hour INT NOT NULL DEFAULT 0,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_resources_owner (owner_user_id),
  CONSTRAINT fk_resources_owner
  FOREIGN KEY (owner_user_id) REFERENCES users(id)
  ON DELETE RESTRICT ON UPDATE CASCADE,
  KEY idx_resources_category_id (category_id),
  UNIQUE KEY uq_resources_category_title (category_id, title),
  CONSTRAINT fk_resources_category
    FOREIGN KEY (category_id)
    REFERENCES resource_categories(id)
    ON DELETE RESTRICT
    ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;