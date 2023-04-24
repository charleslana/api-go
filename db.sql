create database api_go

create table tb_user (id serial primary key, email varchar unique not null, password varchar not null, name varchar unique)

create table tb_permission (id serial primary key, name varchar, user_id int not null, constraint fk_user_permission foreign key (user_id) references tb_user(id))