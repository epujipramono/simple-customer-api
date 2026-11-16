CREATE DATABASE p2gc01_customerapi;
USE p2gc01_customerapi;

CREATE TABLE IF NOT EXISTS customers(
    id int primary key auto_increment,
    name varchar(255) not null,
    email varchar(255) unique not null,
    phone varchar(15) not null,
    created_at timestamp,
    updated_at timestamp
);
