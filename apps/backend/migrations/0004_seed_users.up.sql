INSERT INTO users (email, name, role, password_hash) VALUES
  ('admin@bookinghub.local',   'Администратор', 'ADMIN',   'TEMP'),
  ('manager@bookinghub.local', 'Менеджер',      'MANAGER', 'TEMP'),
  ('user@bookinghub.local',    'Пользователь',  'USER',    'TEMP')
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  role = VALUES(role);
