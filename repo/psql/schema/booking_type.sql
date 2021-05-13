CREATE TABLE booking_type(id INT PRIMARY KEY NOT NULL, description VARCHAR(20) NOT NULL);
INSERT INTO booking_type VALUES(0, 'remote');
INSERT INTO booking_type VALUES(1, 'at clinician address');
INSERT INTO booking_type VALUES(2, 'at patient address');

ALTER TABLE clinician_booking ADD COLUMN booking_type_id INT REFERENCES booking_type(id) NOT NULL DEFAULT 0;
ALTER TABLE office_hours ADD COLUMN booking_type_id INT REFERENCES booking_type(id) NOT NULL DEFAULT 0;

ALTER TABLE business ADD COLUMN address_id INT REFERENCES address(id);

ALTER TABLE calendar_settings ADD COLUMN new_patient_allowed BOOL NOT NULL DEFAULT TRUE;