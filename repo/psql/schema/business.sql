CREATE TABLE business (
                          id SERIAL PRIMARY KEY,
                          person_id INT REFERENCES person(id) UNIQUE,
                          name VARCHAR(100),
                          identifier VARCHAR(100),
                          tax_exemption_id INT REFERENCES tax_exemption(id)
);

CREATE TABLE tax_exemption (
                               id SERIAL PRIMARY KEY,
                               discount_percent INT NOT NULL,
                               code VARCHAR(100)
);