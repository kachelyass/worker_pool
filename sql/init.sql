CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    description TEXT
);

INSERT INTO tasks(description) VALUES('task1'), ('task2');
