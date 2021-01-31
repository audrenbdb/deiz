package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
	"time"
)

func insertBooking(ctx context.Context, db db, b *deiz.Booking) error {
	const query = `INSERT INTO clinician_booking(address_id, blocked, remote, clinician_person_id, patient_id, booking_motive_id, during, paid, note)
	VALUES(NULLIF($1, 0), $2, $3, $4, $5, NULLIF($6, 0), tsrange($7, $8, '()'), $9, NULLIF($10, ''))
	RETURNING id, delete_id`
	row := db.QueryRow(ctx, query, b.Address.ID, b.Blocked, b.Remote, b.Clinician.ID, b.Patient.ID, b.Motive.ID, b.Start, b.End, b.Paid, b.Note)
	err := row.Scan(&b.ID, &b.DeleteID)
	if err != nil {
		return err
	}
	return nil
}

func deleteBooking(ctx context.Context, db db, bookingID, clinicianID int) error {
	const query = `DELETE FROM clinician_booking WHERE clinician_person_id = $1 AND id = $2`
	cmdTag, err := db.Exec(ctx, query, clinicianID, bookingID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errNothingDeleted
	}
	return nil
}

func updateBookingPaidStatus(ctx context.Context, db db, paid bool, bookingID int, clinicianID int) error {
	const query = `UPDATE clinician_booking SET paid = $1 WHERE clinician_person_id = $2 AND id = $3`
	cmdTag, err := db.Exec(ctx, query, paid, clinicianID, bookingID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errNoRowsUpdated
	}
	return nil
}

//Create a booking slot and a patient if he does not exist
func (r *repo) FillFreeBookingSlot(ctx context.Context, b *deiz.Booking) error {
	tx, err := r.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}
	//create patient if not exist
	if b.Patient.ID == 0 {
		err := insertPatient(ctx, tx, &b.Patient, b.Clinician.ID)
		if err != nil {
			return err
		}
	}

	err = insertBooking(ctx, tx, b)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *repo) RemoveBookingSlot(ctx context.Context, s *deiz.Booking) error {
	return deleteBooking(ctx, r.conn, s.ID, s.Clinician.ID)
}

func (r *repo) GetBookingsInTimeRange(ctx context.Context, from, to time.Time, clinicianID int) ([]deiz.Booking, error) {
	const query = `SELECT b.id, b.delete_id, lower(b.during), upper(b.during),
	COALESCE(m.id, 0), COALESCE(m.duration, 0), COALESCE(m.price, 0), COALESCE(m.public, false),
	c.id, c.surname, c.name, c.phone,
	COALESCE(p.id, 0), COALESCE(p.surname, ''), COALESCE(p.name, ''), COALESCE(p.phone, ''), COALESCE(p.email, ''),
	COALESCE(a.id, 0), COALESCE(a.line, ''), COALESCE(a.post_code, 0), COALESCE(a.city, ''),
	b.remote, b.paid, b.blocked, COALESCE(b.note, '')
	FROM clinician_booking b
	LEFT JOIN booking_motive m ON b.booking_motive_id = m.id
	LEFT JOIN patient p ON b.patient_id = p.id
	LEFT JOIN person c ON b.clinician_person_id = c.id
	LEFT JOIN address a ON b.address_id = a.id
	WHERE b.clinician_person_id = $1 AND $2 <= upper(b.during) AND lower(b.during) <= $3`
	rows, err := r.conn.Query(ctx, query, clinicianID, from, to)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var bookingSlots []deiz.Booking
	for rows.Next() {
		var b deiz.Booking
		err := rows.Scan(&b.ID, &b.DeleteID, &b.Start, &b.End, &b.Motive.ID, &b.Motive.Duration, &b.Motive.Price, &b.Motive.Public,
			&b.Clinician.ID, &b.Clinician.Surname, &b.Clinician.Name, &b.Clinician.Phone,
			&b.Patient.ID, &b.Patient.Surname, &b.Patient.Name, &b.Patient.Phone, &b.Patient.Email,
			&b.Address.ID, &b.Address.Line, &b.Address.PostCode, &b.Address.City,
			&b.Remote, &b.Paid, &b.Blocked, &b.Note)
		if err != nil {
			return nil, err
		}
		bookingSlots = append(bookingSlots, b)
	}
	return bookingSlots, nil
}
