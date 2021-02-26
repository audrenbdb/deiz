CREATE TABLE booking_invoice (
                         id SERIAL PRIMARY KEY,
                         person_id INT REFERENCES person(id) ON DELETE CASCADE,
                         booking_id INT REFERENCES clinician_booking(id) ON DELETE SET NULL,
                         created_at TIMESTAMP NOT NULL DEFAULT timezone('utc', NOW()),
                         identifier VARCHAR(100) NOT NULL
                             CONSTRAINT identifier_length CHECK (CHAR_LENGTH(identifier) > 8),
                         sender VARCHAR(50)[] NOT NULL,
                         recipient VARCHAR(50)[] NOT NULL,
                         city_and_date VARCHAR(100) NOT NULL,
                         label VARCHAR(50) NOT NULL,
                         price_before_tax INT NOT NULL,
                         price_after_tax INT NOT NULL,
                         delivery_date TIMESTAMP NOT NULL,
                         delivery_date_str VARCHAR(50) NOT NULL,
                         tax_fee NUMERIC(6, 2) NOT NULL DEFAULT 20
                             CONSTRAINT tax_fee_min CHECK (tax_fee >= 0),
                         exemption VARCHAR(20),
                         canceled BOOLEAN NOT NULL DEFAULT false,
                         payment_method_id INT REFERENCES payment_method(id)
);

CREATE TABLE payment_method (
                                id SERIAL PRIMARY KEY,
                                name VARCHAR(50) NOT NULL
);