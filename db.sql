create database api_go

create table tb_user (id serial primary key, email varchar unique not null, password varchar not null, name varchar unique);

create table tb_permission (id serial primary key, name varchar, user_id int not null, constraint fk_user_permission foreign key (user_id) references tb_user(id));

create table tb_character (id serial primary key, name varchar unique not null, hp int not null, type varchar not null check(type IN ('dark', 'diamond', 'earth', 'electric', 'fire', 'water')));

insert into tb_character (name, hp, type) values ('Kikflick', 100, 'earth');
insert into tb_character (name, hp, type) values ('Menza', 100, 'fire');
insert into tb_character (name, hp, type) values ('Snorky', 100, 'water');
insert into tb_character (name, hp, type) values ('Amuranther', 150, 'dark');

create table tb_user_character (id serial primary key, level int not null default 1, hp_min int not null, slot int not null default 0, user_id int not null, character_id int not null, constraint fk_user_character foreign key (user_id) references tb_user(id), constraint fk_character foreign key (character_id) references tb_character(id));