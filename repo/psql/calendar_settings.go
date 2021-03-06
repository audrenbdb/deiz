package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
)

func getCalendarSettingsByPersonID(ctx context.Context, db db, personID int) (deiz.CalendarSettings, error) {
	const query = `SELECT s.id, s.remote_allowed, s.new_patient_allowed,
	COALESCE(m.id, 0), COALESCE(m.duration, 60), COALESCE(m.price, 5000), COALESCE(m.name, 'Défaut'), COALESCE(m.public, false),
	t.id, t.name
	FROM calendar_settings s
	LEFT JOIN booking_motive m ON s.default_booking_motive_id = m.id
	INNER JOIN timezone t ON t.id = s.timezone_id
	WHERE s.person_id = $1`
	row := db.QueryRow(ctx, query, personID)
	var s deiz.CalendarSettings
	err := row.Scan(&s.ID, &s.RemoteAllowed, &s.NewPatientAllowed,
		&s.DefaultMotive.ID, &s.DefaultMotive.Duration, &s.DefaultMotive.Price, &s.DefaultMotive.Name, &s.DefaultMotive.Public,
		&s.Timezone.ID, &s.Timezone.Name)
	if err != nil {
		return deiz.CalendarSettings{}, err
	}
	return s, nil
}

func insertCalendarSettings(ctx context.Context, db db, s *deiz.CalendarSettings, personID int) error {
	const query = `INSERT INTO calendar_settings(person_id, default_booking_motive_id, timezone_id) VALUES($1, NULLIF($2, 0), COALESCE(NULLIF($3, 0), 1)) RETURNING id`
	row := db.QueryRow(ctx, query, personID, s.DefaultMotive.ID, s.Timezone.ID)
	return row.Scan(&s.ID)
}

func (r *Repo) UpdateCalendarSettings(ctx context.Context, s *deiz.CalendarSettings, clinicianID int) error {
	const query = `UPDATE calendar_settings SET default_booking_motive_id = NULLIF($1, 0), remote_allowed = $2, new_patient_allowed = $3 WHERE person_id = $4`
	tag, err := r.conn.Exec(ctx, query, s.DefaultMotive.ID, s.RemoteAllowed, s.NewPatientAllowed, clinicianID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errNoRowsUpdated
	}
	return nil
}

func (r *Repo) GetClinicianCalendarSettings(ctx context.Context, clinicianID int) (deiz.CalendarSettings, error) {
	return getCalendarSettingsByPersonID(ctx, r.conn, clinicianID)
}

func (r *Repo) GetClinicianTimezone(ctx context.Context, clinicianID int) (deiz.Timezone, error) {
	const query = `SELECT t.id, t.name FROM timezone t INNER JOIN calendar_settings c ON t.id = c.timezone_id WHERE c.person_id = $1`
	row := r.conn.QueryRow(ctx, query, clinicianID)
	var tz deiz.Timezone
	err := row.Scan(&tz.ID, &tz.Name)
	if err != nil {
		return deiz.Timezone{}, err
	}
	return tz, nil

}
