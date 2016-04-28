DROP DATABASE IF EXISTS testdb;
CREATE DATABASE testdb;

USE testdb;

CREATE USER 'testuser'@'localhost' IDENTIFIED BY 'testuser';
GRANT ALL PRIVILEGES ON *.* TO 'testuser'@'localhost' WITH GRANT OPTION;
FLUSH PRIVILEGES;
