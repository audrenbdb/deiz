CREATE TABLE clinician_booking (
                                   id SERIAL PRIMARY KEY,
                                   created_at TIMESTAMP DEFAULT NOW(),
                                   confirmed BOOLEAN NOT NULL DEFAULT false,
                                   blocked BOOLEAN NOT NULL DEFAULT false,
                                   address_id INT REFERENCES address(id) ON DELETE SET NULL,
                                   clinician_person_id INT NOT NULL REFERENCES person(id) ON DELETE CASCADE,
                                   patient_id INT REFERENCES patient(id) ON DELETE CASCADE,
                                   delete_id uuid DEFAULT uuid_generate_v4 (),
                                   booking_motive_id INT REFERENCES booking_motive(id) ON DELETE SET NULL,
                                   during TSRANGE NOT NULL,
                                   paid BOOLEAN DEFAULT FALSE,
                                   remote BOOLEAN DEFAULT FALSE,
                                   note TEXT,
                                   FOREIGN KEY (clinician_person_id, patient_id) REFERENCES patient(clinician_person_id, id) ON DELETE CASCADE,
                                   EXCLUDE USING gist (clinician_person_id WITH =, during WITH &&)
);
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";