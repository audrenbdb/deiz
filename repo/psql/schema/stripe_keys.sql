CREATE TABLE stripe_keys (
                             id SERIAL PRIMARY KEY,
                             person_id INT UNIQUE NOT NULL REFERENCES person(id) ON DELETE CASCADE,
                             public VARCHAR(120) DEFAULT NULL,
                             secret BYTEA UNIQUE DEFAULT NULL
);
/* Requirement */
CREATE EXTENSION "pgcrypto";
CREATE EXTENSION pg_trgm;