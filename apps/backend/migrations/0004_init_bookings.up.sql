CREATE TABLE IF NOT EXISTS bookings (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  resource_id BIGINT UNSIGNED NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,

  start_at DATETIME NOT NULL,
  end_at   DATETIME NOT NULL,

  status ENUM('PENDING','APPROVED','REJECTED','CANCELED') NOT NULL DEFAULT 'PENDING',
  manager_comment VARCHAR(255) NULL,

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,

  PRIMARY KEY (id),

  KEY idx_bookings_resource_time (resource_id, start_at, end_at),
  KEY idx_bookings_user (user_id),
  KEY idx_bookings_status (status),

  CONSTRAINT fk_bookings_resource
    FOREIGN KEY (resource_id) REFERENCES resources(id)
    ON DELETE RESTRICT ON UPDATE CASCADE,

  CONSTRAINT fk_bookings_user
    FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE RESTRICT ON UPDATE CASCADE,

  CONSTRAINT chk_booking_time CHECK (end_at > start_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
