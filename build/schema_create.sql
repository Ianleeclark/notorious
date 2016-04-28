DROP DATABASE IF EXISTS testdb;
CREATE DATABASE testdb;

USE testdb;

GRANT USAGE ON *.* TO 'testuser'@'%';
DROP USER 'testuser'@'%';
CREATE USER 'testuser'@'%' IDENTIFIED BY 'testuser';
GRANT ALL PRIVILEGES ON testuser.* TO testuser@localhost WITH GRANT OPTION;
FLUSH PRIVILEGES;
