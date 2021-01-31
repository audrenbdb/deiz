CREATE TABLE address (
                         id SERIAL PRIMARY KEY,
                         line VARCHAR(100) NOT NULL CONSTRAINT min_line_length CHECK(length(line) > 2),
                         post_code NUMERIC(5) NOT NULL CONSTRAINT min_post_code CHECK(post_code > 10000),
                         city VARCHAR(100) NOT NULL CONSTRAINT min_city_length CHECK(length(city) > 2)
);

CREATE TABLE office_address(
                               id SERIAL PRIMARY KEY,
                               person_id INT REFERENCES person(id),
                               address_id INT REFERENCES address(id) ON DELETE CASCADE
);
CREATE UNIQUE index person_address ON office_address(address_id, person_id);