
CREATE TABLE patient (
                         id SERIAL PRIMARY KEY,
                         clinician_person_id INT NOT NULL REFERENCES person(id) ON DELETE CASCADE,
                         email VARCHAR(254) NOT NULL
                             CONSTRAINT email_length CHECK (CHAR_LENGTH(email) > 5)
                             CONSTRAINT valid_email CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
                         name VARCHAR(50) NOT NULL
                             CONSTRAINT name_length CHECK (CHAR_LENGTH(name) > 1),
                         surname VARCHAR(50) NOT NULL
                             CONSTRAINT surname_length CHECK (CHAR_LENGTH(surname) > 1),
                         phone VARCHAR(50) NOT NULL
                             CONSTRAINT phone_length CHECK (CHAR_LENGTH(phone) >= 10),
                         note VARCHAR(255) DEFAULT NULL,
                         address_id INT REFERENCES address(id) ON DELETE SET NULL DEFAULT NULL,
                         UNIQUE (id, clinician_person_id),
                         UNIQUE (email, clinician_person_id)
);
CREATE UNIQUE index clinician_patient_unique ON patient(id, clinician_person_id);
CREATE INDEX trgm_idx_patient ON patient USING GIST (name gist_trgm_ops);