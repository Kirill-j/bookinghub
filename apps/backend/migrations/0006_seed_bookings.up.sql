-- Одна демо-бронь на завтра 10:00-11:00 для user@bookinghub.local, ресурс id=1
INSERT INTO bookings (resource_id, user_id, start_at, end_at, status)
SELECT
  1,
  u.id,
  DATE_ADD(CURDATE(), INTERVAL 1 DAY) + INTERVAL 10 HOUR,
  DATE_ADD(CURDATE(), INTERVAL 1 DAY) + INTERVAL 11 HOUR,
  'PENDING'
FROM users u
WHERE u.email='user@bookinghub.local'
LIMIT 1;