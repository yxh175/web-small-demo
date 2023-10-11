CREATE SCHEMA cache_demo;
use cache_demo;
DROP TABLE if exists products;

CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    p_id Int Not NULL,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    count INT NOT NULL
);