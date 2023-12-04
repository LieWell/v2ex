/** 支持 emoji 表情 **/
SELECT version();
SET NAMES utf8mb4;
ALTER DATABASE v2ex CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE v2ex.member
    CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
SHOW VARIABLES WHERE Variable_name LIKE 'character%' OR Variable_name LIKE 'collation%';

/** 创建表 **/
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