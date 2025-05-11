CREATE TABLE IF NOT EXISTS `users` (
  `id`                 INT(11)       NOT NULL AUTO_INCREMENT,
  `telegram_username`  VARCHAR(255)  NOT NULL,
  `telegram_chat_id`   VARCHAR(255)  NULL,
  `telegram_first_name`VARCHAR(255)   NULL,
  `telegram_last_name` VARCHAR(255)  NULL,
  `created_at`         TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `ip_addresses` (
  `id`         INT(11)       NOT NULL AUTO_INCREMENT,
  `user_id`    INT(11)       NOT NULL,
  `ip_address` VARCHAR(255)  NOT NULL,
  `note`       VARCHAR(255)  DEFAULT NULL,
  `created_at` TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
