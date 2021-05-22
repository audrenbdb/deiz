package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
	"time"
)

func setInvoiceCanceled(ctx context.Context, db db, invoiceID int, clinicianID int) error {
	const query = `UPDATE booking_invoice SET canceled = true WHERE id = $1 AND person_id = $2`
	cmdTag, err := db.Exec(ctx, query, invoiceID, clinicianID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errNoRowsUpdated
	}
	return nil
}

func (r *Repo) GetBookingsPendingPayment(ctx context.Context, clinicianID int) ([]deiz.Booking, error) {
	const query = `SELECT b.id, lower(b.during), upper(b.during), b.booking_type_id, b.note,
	COALESCE(b.description, ''),
	p.id, p.name, p.surname, COALESCE(p.email, ''), p.phone,
	COALESCE(pa.id, 0), COALESCE(pa.line, ''), COALESCE(pa.post_code, 0), COALESCE(pa.city, ''),
	COALESCE(b.address, ''),
	FROM clinician_booking b
	LEFT JOIN booking_motive m ON b.booking_motive_id = m.id
	INNER JOIN patient p ON p.id = b.patient_id
	LEFT JOIN address pa ON pa.id = p.address_id
	WHERE b.clinician_person_id = $1 AND b.booking_type_id = 1 AND b.paid = false AND lower(b.during) < NOW()`
	rows, err := r.conn.Query(ctx, query, clinicianID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var bookings []deiz.Booking
	for rows.Next() {
		var b deiz.Booking
		err := rows.Scan(&b.ID, &b.Start, &b.End, &b.BookingType, &b.Note, &b.Description,
			&b.Patient.ID, &b.Patient.Name, &b.Patient.Surname, &b.Patient.Email, &b.Patient.Phone,
			&b.Patient.Address.ID, &b.Patient.Address.Line, &b.Patient.Address.PostCode, &b.Patient.Address.City)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *Repo) GetPeriodBookingInvoices(ctx context.Context, start time.Time, end time.Time, clinicianID int) ([]deiz.BookingInvoice, error) {
	const query = `SELECT
	i.id, i.created_at, i.identifier, i.sender, i.recipient,
	i.city_and_date, i.delivery_date,
	i.delivery_date_str, i.label, i.price_before_tax, i.price_after_tax, i.tax_fee,
	COALESCE(i.exemption, ''), i.canceled,
	pm.id, pm.name,
	COALESCE(bp.id, 0), COALESCE(bp.name, ''), COALESCE(bp.surname, ''), COALESCE(bp.phone, ''), COALESCE(bp.email, '')
	FROM booking_invoice i INNER JOIN payment_method pm ON i.payment_method_id = pm.id
	LEFT JOIN clinician_booking b ON i.booking_id = b.id
	LEFT JOIN patient bp ON bp.id = b.patient_id
	WHERE i.person_id = $3 AND i.delivery_date >= $1 AND i.delivery_date <= $2`
	rows, err := r.conn.Query(ctx, query, start, end, clinicianID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var invoices []deiz.BookingInvoice
	for rows.Next() {
		var i deiz.BookingInvoice
		err := rows.Scan(&i.ID, &i.CreatedAt, &i.Identifier, &i.Sender, &i.Recipient, &i.CityAndDate, &i.DeliveryDate, &i.DeliveryDateStr,
			&i.Label, &i.PriceBeforeTax, &i.PriceAfterTax, &i.TaxFee, &i.Exemption, &i.Canceled, &i.PaymentMethod.ID, &i.PaymentMethod.Name,
			&i.Booking.Patient.ID, &i.Booking.Patient.Name, &i.Booking.Patient.Surname, &i.Booking.Patient.Phone, &i.Booking.Patient.Email)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, i)
	}
	return invoices, nil
}

func insertBookingInvoice(ctx context.Context, db db, i *deiz.BookingInvoice) error {
	const query = `INSERT INTO booking_invoice
	(person_id, booking_id, created_at, identifier, sender, recipient,
	city_and_date, label, price_before_tax, price_after_tax, delivery_date,
	delivery_date_str, tax_fee, exemption, payment_method_id)
	VALUES($1, NULLIF($2, 0), $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id`
	row := db.QueryRow(ctx, query, i.ClinicianID, i.Booking.ID, i.CreatedAt, i.Identifier, i.Sender,
		i.Recipient, i.CityAndDate, i.Label, i.PriceBeforeTax, i.PriceAfterTax, i.DeliveryDate,
		i.DeliveryDateStr, i.TaxFee, i.Exemption, i.PaymentMethod.ID)
	err := row.Scan(&i.ID)
	return err
}

func (r *Repo) SaveBookingInvoice(ctx context.Context, i *deiz.BookingInvoice) error {
	tx, err := r.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}
	err = insertBookingInvoice(ctx, tx, i)
	if err != nil {
		return err
	}
	err = updateBookingPaidStatus(ctx, tx, true, i.Booking.ID, i.ClinicianID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *Repo) SaveCorrectingBookingInvoice(ctx context.Context, i *deiz.BookingInvoice) error {
	originalInvoiceID := i.ID
	tx, err := r.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}
	err = insertBookingInvoice(ctx, tx, i)
	if err != nil {
		return err
	}
	err = setInvoiceCanceled(ctx, tx, originalInvoiceID, i.ClinicianID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *Repo) CountClinicianInvoices(ctx context.Context, clinicianID int) (int, error) {
	const query = `SELECT COUNT(*) FROM booking_invoice WHERE person_id = $1`
	var count int
	row := r.conn.QueryRow(ctx, query, clinicianID)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
