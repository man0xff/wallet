create database if not exists example default character set utf8;
use example;

create table if not exists wallets (
    id bigint unsigned not null primary key auto_increment,
    name varchar(255) not null unique key,
    balance bigint unsigned not null default 0
) engine=innodb default charset=utf8;

create table if not exists operations (
    id bigint unsigned not null primary key auto_increment,
    time datetime not null default current_timestamp,
    wallet_id bigint unsigned not null,
    direction tinyint unsigned not null,
    amount bigint unsigned not null,
    index (wallet_id, time)
) engine=innodb default charset=utf8;
