/** 支持 emoji 表情 **/
SELECT version();
SET NAMES utf8mb4;
ALTER DATABASE v2ex CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE v2ex.member
    CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
SHOW VARIABLES WHERE Variable_name LIKE 'character%' OR Variable_name LIKE 'collation%';

/** 创建表 **/
create table member
(
    id          int auto_increment primary key,
    number      int           null,
    name        varchar(255)  null,
    website     varchar(2048) null,
    twitter     varchar(2048) null,
    github      varchar(2048) null,
    location    varchar(2048) null,
    tag_line    varchar(2048) null,
    avatar      varchar(2048) null,
    status      varchar(255)  null,
    create_time datetime      null,
    constraint unique_number unique (number)
);