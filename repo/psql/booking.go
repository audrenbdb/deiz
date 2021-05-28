package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/jackc/pgx/v4"
	"time"
)

const bookingSelect = `SELECT b.id, COALESCE(b.description, ''), b.delete_id, lower(b.during), upper(b.during), b.booking_type_id, COALESCE(b.availability_id, 0),
	c.id, c.surname, c.name, c.phone, c.email,
	COALESCE(p.id, 0), COALESCE(p.surname, ''), COALESCE(p.name, ''), COALESCE(p.phone, ''), COALESCE(p.email, ''),
	COALESCE(b.address, ''), COALESCE(b.price, 0),
	b.paid, COALESCE(b.note, ''), b.confirmed, b.recurrence_id
	FROM clinician_booking b
	LEFT JOIN patient p ON b.patient_id = p.id
	LEFT JOIN person c ON b.clinician_person_id = c.id `

func scanBookingRow(row pgx.Row) (deiz.Booking, error) {
	var b deiz.Booking
	err := row.Scan(&b.ID, &b.Description, &b.DeleteID, &b.Start, &b.End, &b.BookingType, &b.AvailabilityType,
		&b.Clinician.ID, &b.Clinician.Surname, &b.Clinician.Name, &b.Clinician.Phone, &b.Clinician.Email,
		&b.Patient.ID, &b.Patient.Surname, &b.Patient.Name, &b.Patient.Phone, &b.Patient.Email,
		&b.Address, &b.Price,
		&b.Paid, &b.Note, &b.Confirmed, &b.Recurrence)
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
	const query = bookingSelect + `WHERE b.paid = false AND b.confirmed = true AND LOWER(during) <= NOW() AND b.booking_type_id = 1 AND b.clinician_person_id = $1`
	return r.queryBookingRows(ctx, query, clinicianID)
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
	const query = `INSERT INTO clinician_booking(address, price, description, booking_type_id, availability_id, clinician_person_id, patient_id, during, paid, note, confirmed, recurrence_id)
	VALUES(NULLIF($1, ''), $2, NULLIF($3, ''), $4, NULLIF($5, 0), $6, NULLIF($7, 0), tsrange($8, $9, '()'), $10, NULLIF($11, ''), $12, $13)
	RETURNING id, delete_id`
	row := r.conn.QueryRow(ctx, query, b.Address, b.Price, b.Description, b.BookingType, b.AvailabilityType, b.Clinician.ID, b.Patient.ID, b.Start, b.End, b.Paid, b.Note, b.Confirmed, b.Recurrence)
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
	return r.queryBookingRow(ctx, query, deleteID)
}

func (r *Repo) GetBookingByID(ctx context.Context, bookingID int) (deiz.Booking, error) {
	return r.queryBookingRow(ctx, bookingSelect+`WHERE b.id = $1`, bookingID)
}

func (r *Repo) queryBookingRow(ctx context.Context, query string, args ...interface{}) (deiz.Booking, error) {
	row := r.conn.QueryRow(ctx, query, args...)
	return scanBookingRow(row)
}

func (r *Repo) queryBookingRows(ctx context.Context, query string, args ...interface{}) ([]deiz.Booking, error) {
	rows, err := r.conn.Query(ctx, query, args...)
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

func (r *Repo) GetClinicianWeeklyRecurrentBookings(ctx context.Context, clinicianID int) ([]deiz.Booking, error) {
	const query = bookingSelect + `WHERE b.recurrence_id = 2 AND b.clinician_person_id = $1`
	return r.queryBookingRows(ctx, query, clinicianID)
}

func (r *Repo) GetBookingsInTimeRange(ctx context.Context, from, to time.Time) ([]deiz.Booking, error) {
	return r.queryBookingRows(
		ctx, bookingSelect+`WHERE lower(b.during) >= $1 AND lower(b.during) < $2`,
		from, to)
}

func (r *Repo) GetNonRecurrentClinicianBookingsInTimeRange(ctx context.Context, from, to time.Time, clinicianID int) ([]deiz.Booking, error) {
	const query = bookingSelect + `WHERE b.clinician_person_id = $1 AND b.recurrence_id = 0 AND  $2 <= upper(b.during) AND lower(b.during) <= $3 ORDER BY lower(b.during) ASC`
	return r.queryBookingRows(ctx, query, clinicianID, from, to)
}

func (r *Repo) GetPatientBookings(ctx context.Context, clinicianID int, patientID int) ([]deiz.Booking, error) {
	const query = bookingSelect + `WHERE b.clinician_person_id = $1 AND b.patient_id = $2`
	return r.queryBookingRows(ctx, query, clinicianID, patientID)
}

func (r *Repo) UpdateBooking(ctx context.Context, b *deiz.Booking) error {
	const query = `UPDATE clinician_booking 
	SET address = NULLIF($1, ''), price = COALESCE($2, 0), description = NULLIF($3, ''), booking_type_id = $4, clinician_person_id = $5, patient_id = $6,
	during = tsrange($7, $8, '()'), paid = $9, note = NULLIF($10, ''), confirmed = $11, availability_id = $12, recurrence_id = $13 WHERE id = $14`
	cmdTag, err := r.conn.Exec(ctx, query, b.Address, b.Price, b.Description, b.BookingType, b.Clinician.ID, b.Patient.ID,
		b.Start, b.End, b.Paid, b.Note, b.Confirmed, b.AvailabilityType, b.Recurrence, b.ID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errNoRowsUpdated
	}
	return nil
}

/*
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
*/

func (r *Repo) DeleteBlockedBookingPrior(ctx context.Context, d time.Time) error {
	const query = `DELETE FROM clinician_booking WHERE availability_id = 0 AND upper(during) < $1`
	_, err := r.conn.Exec(ctx, query, d)
	return err
}
