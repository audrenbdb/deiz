package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
)

func deleteBookingMotive(ctx context.Context, db db, motiveID int, clinicianID int) error {
	const query = `DELETE FROM booking_motive WHERE id = $1 AND person_id = $2`
	tag, err := db.Exec(ctx, query, motiveID, clinicianID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errNothingDeleted
	}
	return nil
}

func getBookingMotivesByPersonID(ctx context.Context, db db, clinicianID int) ([]deiz.BookingMotive, error) {
	const query = `SELECT id, name, duration, price, public FROM booking_motive WHERE person_id = $1`
	rows, err := db.Query(ctx, query, clinicianID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	motives := []deiz.BookingMotive{}
	for rows.Next() {
		var m deiz.BookingMotive
		err := rows.Scan(&m.ID, &m.Name, &m.Duration, &m.Price, &m.Public)
		if err != nil {
			return nil, err
		}
		motives = append(motives, m)
	}
	return motives, nil
}

func (r *Repo) GetBookingMotiveByID(ctx context.Context, motiveID int, clinicianID int) (deiz.BookingMotive, error) {
	const query = `SELECT id, name, duration, price, public FROM booking_motive WHERE person_id = $1 AND id = $2`
	row := r.conn.QueryRow(ctx, query, motiveID, clinicianID)
	var m deiz.BookingMotive
	return m, row.Scan(&m.ID, &m.Name, &m.Duration, &m.Price, &m.Public)
}

func (r *Repo) UpdateBookingMotive(ctx context.Context, m *deiz.BookingMotive, clinicianID int) error {
	const query = `UPDATE booking_motive SET duration = $1, price = $2, name = $3, public = $4 WHERE id = $5 AND person_id = $6`
	tag, err := r.conn.Exec(ctx, query, m.Duration, m.Price, m.Name, m.Public, m.ID, clinicianID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errNoRowsUpdated
	}
	return nil
}

func (r *Repo) CreateBookingMotive(ctx context.Context, b *deiz.BookingMotive, clinicianID int) error {
	const query = `INSERT INTO booking_motive(person_id, duration, price, name, public) VALUES($1, $2, $3, $4, $5) RETURNING id`
	row := r.conn.QueryRow(ctx, query, clinicianID, b.Duration, b.Price, b.Name, b.Public)
	return row.Scan(&b.ID)
}

func (r *Repo) DeleteBookingMotive(ctx context.Context, mID, clinicianID int) error {
	return deleteBookingMotive(ctx, r.conn, mID, clinicianID)
}
