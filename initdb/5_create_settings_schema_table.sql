CREATE TABLE IF NOT EXISTS `settings_schema` (
    `commit_id` VARCHAR(40),
    `settings_header` BOOLEAN,
    `settings_footer` BOOLEAN,
    PRIMARY KEY (`commit_id`),
    CONSTRAINT
        FOREIGN KEY (`commit_id`)
        REFERENCES `repositories` (`commit_id`)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4
