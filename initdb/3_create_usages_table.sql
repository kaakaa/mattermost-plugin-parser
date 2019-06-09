CREATE TABLE IF NOT EXISTS `usages` (
    `commit_id` VARCHAR(40),
    `api` VARCHAR(40) NOT NULL,
    `path` VARCHAR(128) NOT NULL,
    `line` INT NOT NULL,
    `type` VARCHAR(40) NOT NULL,
    PRIMARY KEY (`commit_id`, `api`, `path`, `line`),
    CONSTRAINT
        FOREIGN KEY (`commit_id`)
        REFERENCES `repositories` (`commit_id`)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4