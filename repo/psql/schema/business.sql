CREATE TABLE business (
                          id SERIAL PRIMARY KEY,
                          person_id INT UNIQUE REFERENCES person(id) ON DELETE CASCADE,
                          name VARCHAR(100),
                          identifier VARCHAR(100),
                          tax_exemption_id INT REFERENCES tax_exemption(id)
);

CREATE TABLE tax_exemption (
                               id SERIAL PRIMARY KEY,
                               discount_percent INT NOT NULL,
                               code VARCHAR(100)
);