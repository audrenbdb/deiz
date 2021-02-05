package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
)

func getCalendarSettingsByPersonID(ctx context.Context, db db, personID int) (deiz.CalendarSettings, error) {
	const query = `SELECT s.id, s.step,
	COALESCE(m.id, 0), COALESCE(m.duration, 30), COALESCE(m.price, 5000), COALESCE(m.name, 'défaut'), COALESCE(m.public, false),
	t.id, t.name
	FROM calendar_settings s
	LEFT JOIN booking_motive m ON s.default_booking_motive_id = m.id
	INNER JOIN timezone t ON t.id = s.timezone_id
	WHERE s.person_id = $1`
	row := db.QueryRow(ctx, query, personID)
	var s deiz.CalendarSettings
	err := row.Scan(&s.ID, &s.Step,
		&s.DefaultMotive.ID, &s.DefaultMotive.Duration, &s.DefaultMotive.Price, &s.DefaultMotive.Name, &s.DefaultMotive.Public,
		&s.Timezone.ID, &s.Timezone.Name)
	if err != nil {
		return deiz.CalendarSettings{}, err
	}
	return s, nil
}

func insertCalendarSettings(ctx context.Context, db db, s *deiz.CalendarSettings, personID int) error {
	const query = `INSERT INTO calendar_settings(person_id, default_booking_motive_id, timezone_id, step) VALUES($1, NULLIF($2, 0), COALESCE(NULLIF($3, 0), 1), COALESCE(NULLIF($4, 0), 30)) RETURNING id`
	row := db.QueryRow(ctx, query, personID, s.DefaultMotive.ID, s.Timezone.ID, s.Step)
	return row.Scan(&s.ID)
}

func updateCalendarSettings(ctx context.Context, db db, s *deiz.CalendarSettings, personID int) error {
	const query = `UPDATE calendar_settings SET default_booking_motive_id = NULLIF($1, 0), step = $2 WHERE person_id = $3`
	tag, err := db.Exec(ctx, query, s.DefaultMotive.ID, s.Step, personID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errNoRowsUpdated
	}
	return nil

}

func (r *repo) EditCalendarSettings(ctx context.Context, s *deiz.CalendarSettings, clinicianID int) error {
	return updateCalendarSettings(ctx, r.conn, s, clinicianID)
}

func (r *repo) GetCalendarSettings(ctx context.Context, clinicianID int) (deiz.CalendarSettings, error) {
	return getCalendarSettingsByPersonID(ctx, r.conn, clinicianID)
}

func (r *repo) GetClinicianTimezone(ctx context.Context, clinicianID int) (deiz.Timezone, error) {
	const query = `SELECT t.id, t.name FROM timezone t INNER JOIN calendar_settings c ON t.id = c.timezone_id WHERE c.person_id = $1`
	row := r.conn.QueryRow(ctx, query, clinicianID)
	var tz deiz.Timezone
	err := row.Scan(&tz.ID, &tz.Name)
	if err != nil {
		return deiz.Timezone{}, err
	}
	return tz, nil

}
