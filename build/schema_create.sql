DROP DATABASE IF EXISTS testdb;
CREATE DATABASE testdb;

USE testdb;

 GRANT USAGE ON *.* TO 'testuser'@'localhost';
DROP USER 'testuser'@'localhost';
CREATE USER 'testuser'@'localhost' IDENTIFIED BY 'testuser';
GRANT ALL PRIVILEGES ON testuser.* TO testuser@localhost;
FLUSH PRIVILEGES;
