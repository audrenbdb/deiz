CREATE TABLE adeli (
                       id SERIAL PRIMARY KEY,
                       person_id INT UNIQUE NOT NULL REFERENCES person(id) ON DELETE CASCADE,
                       identifier VARCHAR(50) UNIQUE DEFAULT NULL
);