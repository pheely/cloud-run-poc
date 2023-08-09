USE hr;
CREATE TABLE IF NOT EXISTS employees (
    id INT PRIMARY KEY AUTO_INCREMENT,
    first_name VARCHAR(10) NOT NULL,
    last_name VARCHAR(10) NOT NULL,
    department VARCHAR(10) NOT NULL,
    salary INT NOT NULL,
    age INT NOT NULL
);
INSERT INTO employees VALUES (
    NULL, 
    'John',
    'Doe',
    'Products',
    200000,
    25
);
INSERT INTO employees VALUES (
    NULL, 
    'Jane',
    'Doe',
    'Sales',
    100000,
    22
)