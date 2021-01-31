CREATE TABLE invoice (
                         id SERIAL PRIMARY KEY,
                         person_id INT REFERENCES person(id) ON DELETE CASCADE,
                         created_at TIMESTAMP NOT NULL DEFAULT timezone('utc', NOW()),
                         identifier VARCHAR(100) NOT NULL
                             CONSTRAINT identifier_length CHECK (CHAR_LENGTH(identifier) > 8),
                         sender VARCHAR(50)[] NOT NULL,
                         recipient VARCHAR(50)[] NOT NULL,
                         city_and_date VARCHAR(100) NOT NULL,
                         label VARCHAR(50) NOT NULL,
                         price_before_tax INT NOT NULL
                             CONSTRAINT amount_min CHECK (price_before_tax >= 0),
                         price_after_tax INT NOT NULL
                             CONSTRAINT price_after_tax_min CHECK (price_after_tax >= 0 AND price_after_tax >= price_before_tax),
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

CREATE TABLE clinician_booking_invoice (
                                           id SERIAL PRIMARY KEY,
                                           person_id INT REFERENCES person(id) ON DELETE CASCADE,
                                           clinician_booking_id INT REFERENCES clinician_booking(id) ON DELETE SET NULL,
                                           invoice_id INT REFERENCES invoice(id) ON DELETE CASCADE
);