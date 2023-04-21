create database api_go

create table tb_user (id serial primary key, email varchar unique not null, password varchar not null, name varchar unique)