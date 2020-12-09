CREATE DATABASE user;
CREATE TABLE userinfo (id int not null PRIMARY KEY AUTO_INCREMENT, name varchar(32) not null, age int not null );
INSERT INTO userinfo (name,age) VALUES ('Lili', 18), ('Tom', 20), ('Grace', 30);
