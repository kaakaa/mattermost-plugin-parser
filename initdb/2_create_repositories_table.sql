CREATE TABLE IF NOT EXISTS `repositories` (
    `url` TEXT NOT NULL,
    `commit_id` VARCHAR(40) NOT NULL,
    `created_at` DATETIME,
    PRIMARY KEY (`commit_id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4
