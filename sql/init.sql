create table member(
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