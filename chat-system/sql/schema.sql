CREATE DATABASE IF NOT EXISTS chat_system DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE chat_system;

CREATE TABLE IF NOT EXISTS `user` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `account` VARCHAR(32) NOT NULL,
  `nickname` VARCHAR(64) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `create_time` DATETIME(3) NULL,
  `last_login_time` DATETIME(3) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_account` (`account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `friend` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `friend_id` BIGINT UNSIGNED NOT NULL,
  `status` VARCHAR(16) NOT NULL DEFAULT 'accepted',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_friend` (`user_id`, `friend_id`),
  KEY `idx_friend_friend_id` (`friend_id`),
  KEY `idx_friend_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `message` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `sender_id` BIGINT UNSIGNED NOT NULL,
  `receiver_id` BIGINT UNSIGNED NOT NULL,
  `content` TEXT NOT NULL,
  `status` VARCHAR(16) NOT NULL DEFAULT 'created',
  `send_time` DATETIME(3) NULL,
  PRIMARY KEY (`id`),
  KEY `idx_message_sender_id` (`sender_id`),
  KEY `idx_message_receiver_id` (`receiver_id`),
  KEY `idx_message_status` (`status`),
  KEY `idx_message_send_time` (`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `conversation` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `peer_id` BIGINT UNSIGNED NOT NULL,
  `last_message_id` BIGINT UNSIGNED NOT NULL,
  `status` VARCHAR(16) NOT NULL DEFAULT 'normal',
  `update_time` DATETIME(3) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_peer` (`user_id`, `peer_id`),
  KEY `idx_conversation_last_message_id` (`last_message_id`),
  KEY `idx_conversation_status` (`status`),
  KEY `idx_conversation_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
