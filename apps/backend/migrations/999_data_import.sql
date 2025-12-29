-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Хост: 127.0.0.1:3306
-- Время создания: Дек 29 2025 г., 16:54
-- Версия сервера: 9.1.0
-- Версия PHP: 8.3.14

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- База данных: `bookinghub`
--

-- --------------------------------------------------------

--
-- Структура таблицы `bookings`
--

DROP TABLE IF EXISTS `bookings`;
CREATE TABLE IF NOT EXISTS `bookings` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `resource_id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `start_at` datetime NOT NULL,
  `end_at` datetime NOT NULL,
  `status` enum('PENDING','APPROVED','REJECTED','CANCELED') NOT NULL DEFAULT 'PENDING',
  `manager_comment` varchar(255) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_bookings_resource_time` (`resource_id`,`start_at`,`end_at`),
  KEY `idx_bookings_user` (`user_id`),
  KEY `idx_bookings_status` (`status`)
) ;

--
-- Дамп данных таблицы `bookings`
--

INSERT INTO `bookings` (`id`, `resource_id`, `user_id`, `start_at`, `end_at`, `status`, `manager_comment`, `created_at`, `updated_at`) VALUES
(1, 1, 3, '2025-12-30 10:00:00', '2025-12-30 11:00:00', 'APPROVED', NULL, '2025-12-29 10:34:29', '2025-12-29 10:35:37'),
(2, 5, 3, '2025-12-30 10:00:00', '2025-12-30 11:00:00', 'APPROVED', NULL, '2025-12-29 10:50:19', '2025-12-29 11:09:57'),
(3, 8, 1, '2025-12-30 10:00:00', '2025-12-30 11:00:00', 'CANCELED', NULL, '2025-12-29 10:52:54', '2025-12-29 11:37:36'),
(4, 4, 2, '2025-12-30 10:00:00', '2025-12-30 11:00:00', 'PENDING', NULL, '2025-12-29 11:11:39', NULL);

-- --------------------------------------------------------

--
-- Структура таблицы `resources`
--

DROP TABLE IF EXISTS `resources`;
CREATE TABLE IF NOT EXISTS `resources` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `owner_user_id` bigint UNSIGNED NOT NULL,
  `category_id` bigint UNSIGNED NOT NULL,
  `title` varchar(150) NOT NULL,
  `description` text,
  `location` varchar(150) DEFAULT NULL,
  `price_per_hour` int NOT NULL DEFAULT '0',
  `is_active` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_resources_category_title` (`category_id`,`title`),
  KEY `idx_resources_owner` (`owner_user_id`),
  KEY `idx_resources_category_id` (`category_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Дамп данных таблицы `resources`
--

INSERT INTO `resources` (`id`, `owner_user_id`, `category_id`, `title`, `description`, `location`, `price_per_hour`, `is_active`, `created_at`) VALUES
(1, 1, 1, 'Переговорная А', 'Небольшая переговорная до 6 человек', '1 этаж', 500, 1, '2025-12-29 10:34:29'),
(2, 1, 1, 'Переговорная Б', 'Большая переговорная до 12 человек', '2 этаж', 800, 1, '2025-12-29 10:34:29'),
(3, 1, 2, 'Коворкинг у окна', 'Рабочее место, розетки рядом', 'Open-space', 200, 1, '2025-12-29 10:34:29'),
(4, 1, 3, 'Студия звукозаписи', 'Шумоизоляция, базовый комплект', 'Блок S', 1200, 1, '2025-12-29 10:34:29'),
(5, 1, 5, 'Проектор Epson', 'FullHD проектор для презентаций', 'Склад', 150, 1, '2025-12-29 10:34:29'),
(8, 3, 4, 'awdawd', 'adawdda', 'awdasd', 123, 1, '2025-12-29 10:52:43');

-- --------------------------------------------------------

--
-- Структура таблицы `resource_categories`
--

DROP TABLE IF EXISTS `resource_categories`;
CREATE TABLE IF NOT EXISTS `resource_categories` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_resource_categories_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Дамп данных таблицы `resource_categories`
--

INSERT INTO `resource_categories` (`id`, `name`, `created_at`) VALUES
(1, 'Переговорная', '2025-12-29 10:34:29'),
(2, 'Рабочее место', '2025-12-29 10:34:29'),
(3, 'Студия', '2025-12-29 10:34:29'),
(4, 'Автомобиль', '2025-12-29 10:34:29'),
(5, 'Оборудование', '2025-12-29 10:34:29');

-- --------------------------------------------------------

--
-- Структура таблицы `schema_migrations`
--

DROP TABLE IF EXISTS `schema_migrations`;
CREATE TABLE IF NOT EXISTS `schema_migrations` (
  `version` bigint NOT NULL,
  `dirty` tinyint(1) NOT NULL,
  PRIMARY KEY (`version`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Дамп данных таблицы `schema_migrations`
--

INSERT INTO `schema_migrations` (`version`, `dirty`) VALUES
(8, 0);

-- --------------------------------------------------------

--
-- Структура таблицы `users`
--

DROP TABLE IF EXISTS `users`;
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `email` varchar(190) NOT NULL,
  `name` varchar(120) NOT NULL,
  `role` enum('INDIVIDUAL','COMPANY','ADMIN') NOT NULL DEFAULT 'INDIVIDUAL',
  `password_hash` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_users_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Дамп данных таблицы `users`
--

INSERT INTO `users` (`id`, `email`, `name`, `role`, `password_hash`, `created_at`) VALUES
(1, 'admin@bookinghub.local', 'Администратор', 'ADMIN', '$2a$10$dRDFKQZZUM4GKIn7YcpKNOEuZ8CYxKFOrOh4p6VOU4orTZTD2aJ9O', '2025-12-29 10:34:29'),
(2, 'manager@bookinghub.local', 'Компания', 'COMPANY', '$2a$10$dRDFKQZZUM4GKIn7YcpKNOEuZ8CYxKFOrOh4p6VOU4orTZTD2aJ9O', '2025-12-29 10:34:29'),
(3, 'user@bookinghub.local', 'Пользователь', 'INDIVIDUAL', '$2a$10$dRDFKQZZUM4GKIn7YcpKNOEuZ8CYxKFOrOh4p6VOU4orTZTD2aJ9O', '2025-12-29 10:34:29');

--
-- Ограничения внешнего ключа сохраненных таблиц
--

--
-- Ограничения внешнего ключа таблицы `bookings`
--
ALTER TABLE `bookings`
  ADD CONSTRAINT `fk_bookings_resource` FOREIGN KEY (`resource_id`) REFERENCES `resources` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  ADD CONSTRAINT `fk_bookings_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

--
-- Ограничения внешнего ключа таблицы `resources`
--
ALTER TABLE `resources`
  ADD CONSTRAINT `fk_resources_category` FOREIGN KEY (`category_id`) REFERENCES `resource_categories` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  ADD CONSTRAINT `fk_resources_owner` FOREIGN KEY (`owner_user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
