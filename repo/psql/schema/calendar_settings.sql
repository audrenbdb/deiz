CREATE TABLE calendar_settings (
                                   id SERIAL PRIMARY KEY,
                                   person_id INT UNIQUE NOT NULL REFERENCES person(id) ON DELETE CASCADE,
                                   default_booking_motive_id INT REFERENCES booking_motive(id) ON DELETE SET NULL,
                                   timezone_id INT NOT NULL REFERENCES timezone(id) DEFAULT 1,
                                   remote_allowed BOOL NOT NULL DEFAULT false
);

CREATE TABLE timezone (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(80) NOT NULL
);

ALTER TABLE OFFICE_HOURS ADD COLUMN allow_remote bool NOT NULL DEFAULT FALSE;
ALTER TABLE OFFICE_HOURS ADD COLUMN allow_booking_to_patient_home bool NOT NULL DEFAULT FALSE;
ALTER TABLE OFFICE_HOURS ADD COLUMN allow_new_patient bool NOT NULL DEFAULT TRUE;