CREATE TABLE IF NOT EXISTS `advertistments` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `title` VARCHAR(255) NULL,
    `start_at` DATETIME(6) NULL,
    `end_at` DATETIME(6) NULL,
    `gender` set('M', 'F') NULL,
    `country` set('TW', 'JP') NULL,
    `platform` set('ANDROID', 'IOS', 'WEB') NULL,
    `age_start` INT NULL,
    `age_end` INT NULL,
    `created_at` DATETIME(6) NULL,
    `updated_at` DATETIME(6) NULL,
    PRIMARY KEY (`id`)
);