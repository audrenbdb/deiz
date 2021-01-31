CREATE TABLE person (
                        id SERIAL PRIMARY KEY,
                        role INT NOT NULL REFERENCES role(level),
                        address_id INT UNIQUE REFERENCES address(id) ON DELETE SET NULL,
                        profession VARCHAR(80) DEFAULT NULL,
                        name VARCHAR(80) NOT NULL
                            CONSTRAINT name_length CHECK (CHAR_LENGTH(name) > 1),
                        surname VARCHAR(80) NOT NULL
                            CONSTRAINT surname_length CHECK (CHAR_LENGTH(surname) > 1),
                        phone VARCHAR(20) NOT NULL
                            CONSTRAINT phone_length CHECK (CHAR_LENGTH(phone) > 9),
                        email VARCHAR(254) UNIQUE NOT NULL
                            CONSTRAINT email_length CHECK (CHAR_LENGTH(email) > 5)
                            CONSTRAINT valid_email CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
                        created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE role (
                      id SERIAL PRIMARY KEY,
                      level INT UNIQUE NOT NULL,
                      name VARCHAR(50) NOT NULL
);
