INSERT INTO users (email, name, role, password_hash) VALUES
  ('admin@bookinghub.local',   'Администратор', 'ADMIN',   '$2a$10$dRDFKQZZUM4GKIn7YcpKNOEuZ8CYxKFOrOh4p6VOU4orTZTD2aJ9O'),
  ('manager@bookinghub.local', 'Менеджер',      'MANAGER', '$2a$10$dRDFKQZZUM4GKIn7YcpKNOEuZ8CYxKFOrOh4p6VOU4orTZTD2aJ9O'),
  ('user@bookinghub.local',    'Пользователь',  'USER',    '$2a$10$dRDFKQZZUM4GKIn7YcpKNOEuZ8CYxKFOrOh4p6VOU4orTZTD2aJ9O')
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  role = VALUES(role),
  password_hash = VALUES(password_hash);