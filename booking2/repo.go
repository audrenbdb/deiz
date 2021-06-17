package booking2

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/psql"
	"github.com/jackc/pgx/v4"
)

type getClinicianNonRecurrentBookingsInTimeRange = func(ctx context.Context, tr timeRange, clinicianID int) ([]deiz.Booking, error)
type getClinicianRecurrentBookings = func(ctx context.Context, clinicianID int) ([]deiz.Booking, error)
type getClinicianOfficeHours = func(ctx context.Context, clinicianID int) ([]officeHours, error)

type psqlRepo struct {
	db psql.PGX
}

func (r *psqlRepo) makeGetClinicianOfficeHours() getClinicianOfficeHours {
	return func(ctx context.Context, clinicianID int) ([]officeHours, error) {
		return r.queryOfficeHours(ctx, `WHERE h.person_id = $1`, clinicianID)
	}
}

func (r *psqlRepo) queryOfficeHours(ctx context.Context, queryConditions string, args ...interface{}) ([]officeHours, error) {
	const query = `SELECT h.start_mn, h.end_mn, h.week_day, h.meeting_mode_id, tz.name,
	COALESCE(a.id, 0), COALESCE(a.line, ''), COALESCE(a.post_code, 0), COALESCE(a.city, '')
	FROM office_hours h LEFT JOIN address a ON h.address_id = a.id
	INNER JOIN timezone tz ON tz.id = h.timezone_id
	`
	rows, err := r.db.Query(ctx, query+` `+queryConditions, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return r.scanOfficeHoursRows(rows)
}

func (r *psqlRepo) scanOfficeHoursRows(rows pgx.Rows) ([]officeHours, error) {
	officeHours := []officeHours{}
	for rows.Next() {
		h, err := r.scanSingleOfficeHoursRow(rows)
		if err != nil {
			return nil, err
		}
		officeHours = append(officeHours, h)
	}
	return officeHours, nil
}

func (r *psqlRepo) scanSingleOfficeHoursRow(row pgx.Row) (officeHours, error) {
	var a address
	var h officeHours
	var tz string
	err := row.Scan(&h.tRangeCfg.startMn, &h.tRangeCfg.endMn, &h.tRangeCfg.weekDay, &h.meetingMode, &tz,
		&a.id, &a.line, &a.postCode, &a.city)
	if err != nil {
		return h, err
	}
	h.tRangeCfg.loc = parseTimezone(tz)
	h.address = a.toString()
	return h, nil
}

const psqlBookingSelect = `SELECT b.id, COALESCE(b.description, ''), lower(b.during), upper(b.during), tz.name, b.booking_type_id, COALESCE(b.meeting_mode_id, 0),
	c.id, c.surname, c.name, c.phone, c.email,
	COALESCE(p.id, 0), COALESCE(p.surname, ''), COALESCE(p.name, ''), COALESCE(p.phone, ''), COALESCE(p.email, ''),
	COALESCE(b.address, ''),
	b.confirmed, b.recurrence_id
	FROM clinician_booking b
	INNER JOIN timezone tz ON tz.id = b.timezone_id
	LEFT JOIN patient p ON b.patient_id = p.id
	LEFT JOIN person c ON b.clinician_person_id = c.id`

func (r *psqlRepo) queryBookings(ctx context.Context, queryConditions string, args ...interface{}) ([]deiz.Booking, error) {
	rows, err := r.db.Query(ctx, psqlBookingSelect+` `+queryConditions, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return r.scanBookingRows(rows)
}

func (r *psqlRepo) scanBookingRows(rows pgx.Rows) ([]deiz.Booking, error) {
	bookings := []deiz.Booking{}
	for rows.Next() {
		b, err := r.scanSingleBookingRow(rows)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *psqlRepo) scanSingleBookingRow(row pgx.Row) (b deiz.Booking, err error) {
	return b, row.Scan(&b.ID, &b.Description, &b.Start, &b.End, &b.Timezone,
		&b.BookingType, &b.MeetingMode,
		&b.Clinician.ID, &b.Clinician.Surname, &b.Clinician.Name, &b.Clinician.Phone, &b.Clinician.Email,
		&b.Patient.ID, &b.Patient.Surname, &b.Patient.Name, &b.Patient.Phone, &b.Patient.Email,
		&b.Address, &b.Confirmed, &b.Recurrence)
}

func (r *psqlRepo) makeGetClinicianNonRecurrentBookingsInTimeRange() getClinicianNonRecurrentBookingsInTimeRange {
	return func(ctx context.Context, tr timeRange, clinicianID int) ([]deiz.Booking, error) {
		queryConditions := `WHERE b.clinician_person_id = $1 AND $2 <= upper(b.during) AND lower(b.during) <= $3 AND b.recurrence_id != 2`
		return r.queryBookings(ctx, queryConditions, clinicianID, tr.start, tr.end)
	}
}

func (r *psqlRepo) makeGetClinicianRecurrentBookings() getClinicianRecurrentBookings {
	return func(ctx context.Context, clinicianID int) ([]deiz.Booking, error) {
		queryConditions := `WHERE b.clinician_person_id = $1 AND b.recurrence_id = 2`
		return r.queryBookings(ctx, queryConditions, clinicianID)
	}
}
