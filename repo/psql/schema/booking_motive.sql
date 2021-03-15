CREATE TABLE booking_motive (
                                id SERIAL PRIMARY KEY,
                                person_id INT NOT NULL REFERENCES person(id) ON DELETE CASCADE,
                                duration INT NOT NULL DEFAULT 30,
                                price INT NOT NULL DEFAULT 50,
                                name VARCHAR(80) NOT NULL
                                    CONSTRAINT name_length CHECK (CHAR_LENGTH(name) > 1),
                                public BOOLEAN NOT NULL default false
);
CREATE UNIQUE index id_clinician_id_booking ON booking_motive(id, person_id);