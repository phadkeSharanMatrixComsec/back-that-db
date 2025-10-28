CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    department VARCHAR(50),
    salary INT
);

INSERT INTO employees (name, department, salary) VALUES
('Alice', 'Engineering', 80000),
('Bob', 'Sales', 60000),
('Charlie', 'HR', 50000),
('David', 'Finance', 70000);
