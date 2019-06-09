CREATE TABLE IF NOT EXISTS `props` (
    `commit_id` VARCHAR(40) NOT NULL,
    `key` VARCHAR(128) NOT NULL,
    `type` VARCHAR(128) NOT NULL,
    PRIMARY KEY (`commit_id`, `key`),
    CONSTRAINT
        FOREIGN KEY (`commit_id`)
        REFERENCES `repositories` (`commit_id`)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4
