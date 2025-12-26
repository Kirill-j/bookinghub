INSERT INTO bookings (resource_id, user_id, start_at, end_at, status)
SELECT
  r.id,
  u.id,
  DATE_ADD(CURDATE(), INTERVAL 1 DAY) + INTERVAL 10 HOUR,
  DATE_ADD(CURDATE(), INTERVAL 1 DAY) + INTERVAL 11 HOUR,
  'PENDING'
FROM users u
JOIN resources r ON r.title = 'Переговорная А'
WHERE u.email = 'user@bookinghub.local'
AND NOT EXISTS (
  SELECT 1 FROM bookings b
  WHERE b.user_id = u.id
    AND b.resource_id = r.id
    AND b.start_at = DATE_ADD(CURDATE(), INTERVAL 1 DAY) + INTERVAL 10 HOUR
    AND b.end_at   = DATE_ADD(CURDATE(), INTERVAL 1 DAY) + INTERVAL 11 HOUR
);
