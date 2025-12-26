DELETE b
FROM bookings b
JOIN users u ON u.id = b.user_id
JOIN resources r ON r.id = b.resource_id
WHERE u.email = 'user@bookinghub.local'
  AND r.title = 'Переговорная А';
