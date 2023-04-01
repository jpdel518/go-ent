-- Create "groups" table
CREATE TABLE `groups` (`id` bigint NOT NULL AUTO_INCREMENT, `created_at` timestamp NOT NULL, `updated_at` timestamp NOT NULL, `name` varchar(255) NOT NULL, PRIMARY KEY (`id`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "users" table
CREATE TABLE `users` (`id` bigint NOT NULL AUTO_INCREMENT, `created_at` timestamp NOT NULL, `updated_at` timestamp NOT NULL, `first_name` varchar(20) NOT NULL, `last_name` varchar(20) NOT NULL, `email` varchar(50) NOT NULL, `age` bigint NULL, PRIMARY KEY (`id`), UNIQUE INDEX `email` (`email`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "cars" table
CREATE TABLE `cars` (`id` bigint NOT NULL AUTO_INCREMENT, `created_at` timestamp NOT NULL, `updated_at` timestamp NOT NULL, `name` varchar(255) NOT NULL, `model` varchar(255) NOT NULL, `registered_at` timestamp NOT NULL, `user_cars` bigint NULL, PRIMARY KEY (`id`), CONSTRAINT `cars_users_cars` FOREIGN KEY (`user_cars`) REFERENCES `users` (`id`) ON DELETE SET NULL) CHARSET utf8mb4 COLLATE utf8mb4_bin;
