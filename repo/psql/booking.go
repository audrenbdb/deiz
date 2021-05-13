package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/jackc/pgx/v4"
	"time"
)

const bookingSelect = `SELECT b.id, b.delete_id, lower(b.during), upper(b.during), b.booking_type_id,
	COALESCE(m.id, 0), COALESCE(m.name, ''), COALESCE(m.duration, 0), COALESCE(m.price, 0), COALESCE(m.public, false),
	c.id, c.surname, c.name, c.phone, c.email,
	COALESCE(p.id, 0), COALESCE(p.surname, ''), COALESCE(p.name, ''), COALESCE(p.phone, ''), COALESCE(p.email, ''),
	COALESCE(pa.id, 0), COALESCE(pa.line, ''), COALESCE(pa.post_code, 0), COALESCE(pa.city, ''),
	COALESCE(a.id, 0), COALESCE(a.line, ''), COALESCE(a.post_code, 0), COALESCE(a.city, ''),
	b.paid, b.blocked, COALESCE(b.note, ''), b.confirmed
	FROM clinician_booking b
	LEFT JOIN booking_motive m ON b.booking_motive_id = m.id
	LEFT JOIN patient p ON b.patient_id = p.id
	LEFT JOIN address pa ON p.address_id = pa.id
	LEFT JOIN person c ON b.clinician_person_id = c.id
	LEFT JOIN address a ON b.address_id = a.id `

func scanBookingRow(row pgx.Row) (deiz.Booking, error) {
	var b deiz.Booking
	err := row.Scan(&b.ID, &b.DeleteID, &b.Start, &b.End, &b.BookingType, &b.Motive.ID, &b.Motive.Name, &b.Motive.Duration, &b.Motive.Price, &b.Motive.Public,
		&b.Clinician.ID, &b.Clinician.Surname, &b.Clinician.Name, &b.Clinician.Phone, &b.Clinician.Email,
		&b.Patient.ID, &b.Patient.Surname, &b.Patient.Name, &b.Patient.Phone, &b.Patient.Email,
		&b.Patient.Address.ID, &b.Patient.Address.Line, &b.Patient.Address.PostCode, &b.Patient.Address.City,
		&b.Address.ID, &b.Address.Line, &b.Address.PostCode, &b.Address.City,
		&b.Paid, &b.Blocked, &b.Note, &b.Confirmed)
	return b, err
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

func (r *Repo) GetUnpaidBookings(ctx context.Context, clinicianID int) ([]deiz.Booking, error) {
	const query = `SELECT b.id, b.delete_id, lower(b.during), upper(b.during), b.booking_type_id,
	COALESCE(m.id, 0), COALESCE(m.name, ''), COALESCE(m.duration, 0), COALESCE(m.price, 0), COALESCE(m.public, false),
	c.id, c.surname, c.name, c.phone,
	COALESCE(p.id, 0), COALESCE(p.surname, ''), COALESCE(p.name, ''), COALESCE(p.phone, ''), COALESCE(p.email, ''),
	COALESCE(pa.id, 0), COALESCE(pa.line, ''), COALESCE(pa.post_code, 0), COALESCE(pa.city, ''),
	COALESCE(a.id, 0), COALESCE(a.line, ''), COALESCE(a.post_code, 0), COALESCE(a.city, ''),
	b.paid, b.blocked, COALESCE(b.note, ''), b.confirmed
	FROM clinician_booking b
	LEFT JOIN booking_motive m ON b.booking_motive_id = m.id
	LEFT JOIN patient p ON b.patient_id = p.id
	LEFT JOIN person c ON b.clinician_person_id = c.id
	LEFT JOIN address a ON b.address_id = a.id
	LEFT JOIN address pa ON p.address_id = pa.id
    WHERE b.paid = false AND b.confirmed = true AND LOWER(during) <= NOW() AND b.clinician_person_id = $1`
	rows, err := r.conn.Query(ctx, query, clinicianID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var bookings []deiz.Booking
	for rows.Next() {
		var b deiz.Booking
		err := rows.Scan(&b.ID, &b.DeleteID, &b.Start, &b.End, &b.BookingType, &b.Motive.ID, &b.Motive.Name, &b.Motive.Duration, &b.Motive.Price, &b.Motive.Public,
			&b.Clinician.ID, &b.Clinician.Surname, &b.Clinician.Name, &b.Clinician.Phone,
			&b.Patient.ID, &b.Patient.Surname, &b.Patient.Name, &b.Patient.Phone, &b.Patient.Email,
			&b.Patient.Address.ID, &b.Patient.Address.Line, &b.Patient.Address.PostCode, &b.Patient.Address.City,
			&b.Address.ID, &b.Address.Line, &b.Address.PostCode, &b.Address.City,
			&b.Paid, &b.Blocked, &b.Note, &b.Confirmed)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *Repo) DeleteBooking(ctx context.Context, bookingID int, clinicianID int) error {
	const query = `DELETE FROM clinician_booking WHERE clinician_person_id = $1 AND id = $2`
	cmdTag, err := r.conn.Exec(ctx, query, clinicianID, bookingID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errNothingDeleted
	}
	return nil
}

func (r *Repo) CreateBooking(ctx context.Context, b *deiz.Booking) error {
	const query = `INSERT INTO clinician_booking(address_id, blocked, booking_type_id, clinician_person_id, patient_id, booking_motive_id, during, paid, note, confirmed)
	VALUES(NULLIF($1, 0), $2, $3, $4, NULLIF($5, 0), NULLIF($6, 0), tsrange($7, $8, '()'), $9, NULLIF($10, ''), $11)
	RETURNING id, delete_id`
	row := r.conn.QueryRow(ctx, query, b.Address.ID, b.Blocked, b.BookingType, b.Clinician.ID, b.Patient.ID, b.Motive.ID, b.Start, b.End, b.Paid, b.Note, b.Confirmed)
	err := row.Scan(&b.ID, &b.DeleteID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) DeleteBlockedBookingsInTimeRange(ctx context.Context, start, end time.Time, clinicianID int) error {
	const query = `DELETE FROM clinician_booking
	WHERE clinician_person_id = $1 AND $2 < upper(during) AND lower(during) < $3 AND patient_id IS NULL`
	_, err := r.conn.Exec(ctx, query, clinicianID, start, end)
	return err
}

func (r *Repo) GetBookingByDeleteID(ctx context.Context, deleteID string) (deiz.Booking, error) {
	const query = bookingSelect + `WHERE b.delete_id = $1`
	row := r.conn.QueryRow(ctx, query, deleteID)
	b, err := scanBookingRow(row)
	if err != nil {
		return deiz.Booking{}, err
	}
	return b, nil
}

func (r *Repo) GetBookingByID(ctx context.Context, bookingID int) (deiz.Booking, error) {
	const query = bookingSelect + `WHERE b.id = $1`
	row := r.conn.QueryRow(ctx, query, bookingID)
	b, err := scanBookingRow(row)
	if err != nil {
		return deiz.Booking{}, err
	}
	return b, nil
}

func (r *Repo) GetBookingsInTimeRange(ctx context.Context, from, to time.Time) ([]deiz.Booking, error) {
	const query = bookingSelect + `WHERE lower(b.during) >= $1 AND lower(b.during) < $2`
	rows, err := r.conn.Query(ctx, query, from, to)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var bookingSlots []deiz.Booking
	for rows.Next() {
		b, err := scanBookingRow(rows)
		if err != nil {
			return nil, err
		}
		bookingSlots = append(bookingSlots, b)
	}
	return bookingSlots, nil
}

func (r *Repo) GetClinicianBookingsInTimeRange(ctx context.Context, from, to time.Time, clinicianID int) ([]deiz.Booking, error) {
	const query = bookingSelect + `WHERE b.clinician_person_id = $1 AND $2 <= upper(b.during) AND lower(b.during) <= $3 ORDER BY lower(b.during) ASC`
	rows, err := r.conn.Query(ctx, query, clinicianID, from, to)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var bookingSlots []deiz.Booking
	for rows.Next() {
		b, err := scanBookingRow(rows)
		if err != nil {
			return nil, err
		}
		bookingSlots = append(bookingSlots, b)
	}
	return bookingSlots, nil
}

func (r *Repo) GetPatientBookings(ctx context.Context, clinicianID int, patientID int) ([]deiz.Booking, error) {
	const query = bookingSelect + `WHERE b.clinician_person_id = $1 AND b.patient_id = $2`
	rows, err := r.conn.Query(ctx, query, clinicianID, patientID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var bookings []deiz.Booking
	for rows.Next() {
		b, err := scanBookingRow(rows)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *Repo) UpdateBooking(ctx context.Context, b *deiz.Booking) error {
	const query = `UPDATE clinician_booking 
	SET address_id = NULLIF($1, 0), blocked = $2, booking_type_id = $3, clinician_person_id = $4, patient_id = $5,
	booking_motive_id = NULLIF($6, 0), during = tsrange($7, $8, '()'), paid = $9, note = NULLIF($10, ''), confirmed = $11 WHERE id = $12`
	cmdTag, err := r.conn.Exec(ctx, query, b.Address.ID, b.Blocked, b.BookingType, b.Clinician.ID, b.Patient.ID, b.Motive.ID,
		b.Start, b.End, b.Paid, b.Note, b.Confirmed, b.ID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errNoRowsUpdated
	}
	return nil
}

func (r *Repo) IsBookingSlotAvailable(ctx context.Context, from, to time.Time, clinicianID int) (bool, error) {
	const query = `SELECT EXISTS(
		SELECT 1 FROM clinician_booking b 
		WHERE clinician_person_id = $1 
		AND $2 < upper(b.during) AND lower(b.during) < $3 AND b.patient_id IS NOT NULL)`
	var slotTaken bool
	row := r.conn.QueryRow(ctx, query, clinicianID, from, to)
	err := row.Scan(&slotTaken)
	if err != nil {
		return false, err
	}
	return !slotTaken, nil
}

func (r *Repo) DeleteBlockedBookingPrior(ctx context.Context, d time.Time) error {
	const query = `DELETE FROM clinician_booking WHERE blocked IS TRUE AND upper(during) < $1`
	_, err := r.conn.Exec(ctx, query, d)
	return err
}
