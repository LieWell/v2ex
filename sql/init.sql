CREATE TABLE `member`
(
    `id`          int NOT NULL AUTO_INCREMENT,
    `number`      int           DEFAULT NULL,
    `name`        varchar(255)  DEFAULT NULL,
    `website`     varchar(2048) DEFAULT NULL,
    `twitter`     varchar(2048) DEFAULT NULL,
    `github`      varchar(2048) DEFAULT NULL,
    `location`    varchar(2048) DEFAULT NULL,
    `tag_line`    varchar(2048) DEFAULT NULL,
    `avatar`      varchar(2048) DEFAULT NULL,
    `status`      varchar(255)  DEFAULT NULL,
    `create_time` datetime      DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_number` (`number`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE `system_config`
(
    `id`    int NOT NULL AUTO_INCREMENT,
    `key`   varchar(255) DEFAULT NULL,
    `value` longtext     DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_key` (`key`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

create table `avatar`
(
    `id`     int NOT NULL AUTO_INCREMENT,
    `name`   varchar(255)  DEFAULT NULL,
    `avatar` varchar(2048) DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

INSERT INTO `system_config` (`id`, `key`, `value`)
VALUES (1, 'last_draw_time', '0001-01-01 00:00:00');