CREATE TABLE IF NOT EXISTS `manifest` (
    `commit_id` VARCHAR(40),
    `id` VARCHAR(190) NOT NULL,
    `name` VARCHAR(128) NOT NULL,
    `version` VARCHAR(8),
    `min_server_version` VARCHAR(8) NOT NULL,
    PRIMARY KEY (`commit_id`),
    CONSTRAINT
        FOREIGN KEY (`commit_id`)
        REFERENCES `repositories` (`commit_id`)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4
