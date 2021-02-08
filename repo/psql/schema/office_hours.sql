CREATE TABLE office_hours (
                              id SERIAL PRIMARY KEY,
                              person_id INT NOT NULL REFERENCES person(id),
                              start_mn INT NOT NULL
                                  CONSTRAINT start_mn_val CHECK(start_mn < end_mn),
                              end_mn INT NOT NULL
                                  CONSTRAINT end_mn_val CHECK(end_mn <= 1440),
                              week_day INT NOT NULL
                                    CONSTRAINT day_of_week_min_max CHECK(week_day >= 0 AND week_day < 7),
                              address_id INT REFERENCES address(id) ON DELETE SET NULL,
                              remote_allowed BOOL NOT NULL DEFAULT FALSE,
                              EXCLUDE using gist (person_id WITH =, week_day WITH =, (array[start_mn, end_mn]) WITH &&)
);
CREATE EXTENSION IF NOT EXISTS intarray;
CREATE EXTENSION IF NOT EXISTS btree_gist;