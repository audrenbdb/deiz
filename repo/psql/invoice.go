package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
	"time"
)

type invoice struct {
	id              int
	personID        int
	createdAt       time.Time
	identifier      string
	sender          []string
	recipient       []string
	cityAndDate     string
	label           string
	priceBeforeTax  int64
	priceAfterTax   int64
	deliveryDate    time.Time
	deliveryDateStr string
	taxFee          float32
	exemption       string
	canceled        bool
	paymentMethodID int
}

type bookingInvoice struct {
	id        int
	personID  int
	invoiceID int
	bookingID int
}

func insertInvoice(ctx context.Context, db db, i *invoice) error {
	const query = `INSERT INTO invoice(person_id, identifier, sender, recipient, city_and_date, label, price_before_tax, price_after_tax, delivery_date, tax_fee, exemption, canceled, payment_method_id, delivery_date_str)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id, created_at`
	row := db.QueryRow(ctx, query, i.personID, i.identifier, i.sender, i.recipient, i.cityAndDate,
		i.label, i.priceBeforeTax, i.priceAfterTax, i.deliveryDate, i.taxFee, i.exemption, i.canceled, i.paymentMethodID, i.deliveryDateStr)
	return row.Scan(&i.id, &i.createdAt)
}

func insertBookingInvoice(ctx context.Context, db db, b *bookingInvoice) error {
	const query = `INSERT INTO clinician_booking_invoice(person_id, clinician_booking_id, invoice_id) VALUES($1, $2, $3) RETURNING id`
	row := db.QueryRow(ctx, query, b.personID, b.bookingID, b.invoiceID)
	return row.Scan(&b.id)
}

func (r *repo) GetBookingsPendingPayment(ctx context.Context, clinicianID int) ([]deiz.Booking, error) {
	const query = `SELECT b.id, lower(b.during), upper(b.during), b.remote, b.note,
	COALESCE(m.id, 0), COALESCE(m.name, ''), COALESCE(m.duration, 0), COALESCE(m.price, 0),
	p.id, p.name, p.surname, p.email, p.phone,
	COALESCE(pa.id, 0), COALESCE(pa.line, ''), COALESCE(pa.post_code, 0), COALESCE(pa.city, ''),
	COALESCE(a.id, 0), COALESCE(a.line, ''), COALESCE(a.post_code, 0), COALESCE(a.city, '')
	FROM clinician_booking b
	LEFT JOIN booking_motive m ON b.booking_motive_id = m.id
	INNER JOIN patient p ON p.id = b.patient_id
	LEFT JOIN address pa ON pa.id = p.address_id
	LEFT JOIN address a ON a.id = b.address_id
	WHERE b.clinician_person_id = $1 AND b.blocked = false AND b.paid = false AND lower(b.during) < NOW()`
	rows, err := r.conn.Query(ctx, query, clinicianID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var bookings []deiz.Booking
	for rows.Next() {
		var b deiz.Booking
		err := rows.Scan(&b.ID, &b.Start, &b.End, &b.Remote, &b.Note,
			&b.Motive.ID, &b.Motive.Name, &b.Motive.Duration, &b.Motive.Price,
			&b.Patient.ID, &b.Patient.Name, &b.Patient.Surname, &b.Patient.Email, &b.Patient.Phone,
			&b.Patient.Address.ID, &b.Patient.Address.Line, &b.Patient.Address.PostCode, &b.Patient.Address.City,
			&b.Address.ID, &b.Address.Line, &b.Address.PostCode, &b.Address.City)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *repo) GetBookingInvoiceByID(ctx context.Context, invoiceID int, clinicianID int) (deiz.BookingInvoice, error) {
	const query = `SELECT i.id, i.created_at, 
	i.identifier, i.sender, i.recipient, i.city_and_date, i.delivery_date, i.delivery_date_str,
	i.label, i.price_before_tax, i.price_after_tax, i.tax_fee, i.exemption, i.canceled,
	pm.id, pm.name,
	COALESCE(b.id, 0), COALESCE(LOWER(b.during), NOW()), COALESCE(UPPER(b.during), NOW()), COALESCE(b.remote, false), COALESCE(b.paid, false), COALESCE(b.blocked, false), COALESCE(b.note, ''),
	COALESCE(bm.id, 0), COALESCE(bm.name, ''), COALESCE(bm.duration, 0), COALESCE(bm.price, 0), COALESCE(bm.public, false),
	COALESCE(ba.id, 0), COALESCE(ba.line, ''), COALESCE(ba.post_code, 0), COALESCE(ba.city, ''),
	COALESCE(bc.id, 0), COALESCE(bc.name, ''), COALESCE(bc.surname, ''), COALESCE(bc.phone, ''), COALESCE(bc.email, ''),
	COALESCE(bp.id, 0), COALESCE(bp.name, ''), COALESCE(bp.surname, ''), COALESCE(bp.phone, ''), COALESCE(bp.email, '')
	FROM clinician_booking_invoice bi
	INNER JOIN invoice i on bi.invoice_id = i.id
	INNER JOIN payment_method pm ON i.payment_method_id = pm.id
	LEFT JOIN clinician_booking b ON bi.clinician_booking_id = b.id
	LEFT JOIN booking_motive bm ON bm.id = b.booking_motive_id
	LEFT JOIN address ba ON ba.id = b.address_id
	LEFT JOIN person bc ON bc.id = b.clinician_person_id
	LEFT JOIN patient bp ON bp.id = b.patient_id
	WHERE i.id = $1 AND clinician_id = $2`
	row := r.conn.QueryRow(ctx, query, invoiceID, clinicianID)
	var i deiz.BookingInvoice
	err := row.Scan(&i.ID, &i.CreatedAt, &i.Identifier, &i.Sender, &i.Recipient, &i.CityAndDate, &i.DeliveryDate, &i.DeliveryDateStr,
		&i.Label, &i.PriceBeforeTax, &i.PriceAfterTax, &i.TaxFee, &i.Exemption, &i.Canceled,
		&i.PaymentMethod.ID, &i.PaymentMethod.Name,
		&i.Booking.ID, &i.Booking.Start, &i.Booking.End, &i.Booking.Remote, &i.Booking.Paid, &i.Booking.Blocked, &i.Booking.Note,
		&i.Booking.Motive.ID, &i.Booking.Motive.Name, &i.Booking.Motive.Duration, &i.Booking.Motive.Price, &i.Booking.Motive.Public,
		&i.Booking.Address.ID, &i.Booking.Address.Line, &i.Booking.Address.PostCode, &i.Booking.Address.City,
		&i.Booking.Clinician.ID, &i.Booking.Clinician.Name, &i.Booking.Clinician.Surname, &i.Booking.Clinician.Phone, &i.Booking.Clinician.Email,
		&i.Booking.Patient.ID, &i.Booking.Patient.Name, &i.Booking.Patient.Surname, &i.Booking.Patient.Phone, &i.Booking.Patient.Email)
	if err != nil {
		return deiz.BookingInvoice{}, err
	}
	return i, nil
}

func (r *repo) GetPeriodBookingInvoices(ctx context.Context, start time.Time, end time.Time, clinicianID int) ([]deiz.BookingInvoice, error) {
	const query = `SELECT i.id, i.created_at,
	i.identifier, i.sender, i.recipient, i.city_and_date, i.delivery_date, i.delivery_date_str,
	i.label, i.price_before_tax, i.price_after_tax, i.tax_fee, i.exemption, i.canceled,
	pm.id, pm.name,
	COALESCE(b.id, 0), COALESCE(LOWER(b.during), NOW()), COALESCE(UPPER(b.during), NOW()), COALESCE(b.remote, false), COALESCE(b.paid, false), COALESCE(b.blocked, false), COALESCE(b.note, ''),
	COALESCE(bm.id, 0), COALESCE(bm.name, ''), COALESCE(bm.duration, 0), COALESCE(bm.price, 0), COALESCE(bm.public, false),
	COALESCE(ba.id, 0), COALESCE(ba.line, ''), COALESCE(ba.post_code, 0), COALESCE(ba.city, ''),
	COALESCE(bc.id, 0), COALESCE(bc.name, ''), COALESCE(bc.surname, ''), COALESCE(bc.phone, ''), COALESCE(bc.email, ''),
	COALESCE(bp.id, 0), COALESCE(bp.name, ''), COALESCE(bp.surname, ''), COALESCE(bp.phone, ''), COALESCE(bp.email, '')
	FROM clinician_booking_invoice bi
	INNER JOIN invoice i on bi.invoice_id = i.id
	INNER JOIN payment_method pm ON i.payment_method_id = pm.id
	LEFT JOIN clinician_booking b ON bi.clinician_booking_id = b.id
	LEFT JOIN booking_motive bm ON bm.id = b.booking_motive_id
	LEFT JOIN address ba ON ba.id = b.address_id
	LEFT JOIN person bc ON bc.id = b.clinician_person_id
	LEFT JOIN patient bp ON bp.id = b.patient_id
	WHERE i.delivery_date >= $1 AND i.delivery_date <= $2 AND person_id = $3`
	rows, err := r.conn.Query(ctx, query, start, end, clinicianID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var invoices []deiz.BookingInvoice
	for rows.Next() {
		var i deiz.BookingInvoice
		err := rows.Scan(&i.ID, &i.CreatedAt, &i.Identifier, &i.Sender, &i.Recipient, &i.CityAndDate, &i.DeliveryDate, &i.DeliveryDateStr,
			&i.Label, &i.PriceBeforeTax, &i.PriceAfterTax, &i.TaxFee, &i.Exemption, &i.Canceled,
			&i.PaymentMethod.ID, &i.PaymentMethod.Name,
			&i.Booking.ID, &i.Booking.Start, &i.Booking.End, &i.Booking.Remote, &i.Booking.Paid, &i.Booking.Blocked, &i.Booking.Note,
			&i.Booking.Motive.ID, &i.Booking.Motive.Name, &i.Booking.Motive.Duration, &i.Booking.Motive.Price, &i.Booking.Motive.Public,
			&i.Booking.Address.ID, &i.Booking.Address.Line, &i.Booking.Address.PostCode, &i.Booking.Address.City,
			&i.Booking.Clinician.ID, &i.Booking.Clinician.Name, &i.Booking.Clinician.Surname, &i.Booking.Clinician.Phone, &i.Booking.Clinician.Email,
			&i.Booking.Patient.ID, &i.Booking.Patient.Name, &i.Booking.Patient.Surname, &i.Booking.Patient.Phone, &i.Booking.Patient.Email)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, i)
	}
	return invoices, nil
}

func (r *repo) CreateBookingInvoice(ctx context.Context, i *deiz.BookingInvoice, clinicianID int) error {
	tx, err := r.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}

	invoice := &invoice{
		personID:        clinicianID,
		identifier:      i.Identifier,
		sender:          i.Sender,
		recipient:       i.Recipient,
		cityAndDate:     i.CityAndDate,
		label:           i.Label,
		priceBeforeTax:  i.PriceBeforeTax,
		priceAfterTax:   i.PriceAfterTax,
		deliveryDate:    i.DeliveryDate,
		taxFee:          i.TaxFee,
		exemption:       i.Exemption,
		canceled:        i.Canceled,
		paymentMethodID: i.PaymentMethod.ID,
		deliveryDateStr: i.DeliveryDateStr,
	}
	err = insertInvoice(ctx, tx, invoice)
	if err != nil {
		return err
	}

	bookingInvoice := &bookingInvoice{personID: clinicianID, bookingID: i.Booking.ID, invoiceID: invoice.id}
	err = insertBookingInvoice(ctx, tx, bookingInvoice)
	if err != nil {
		return err
	}
	i.ID = bookingInvoice.id
	i.CreatedAt = invoice.createdAt

	err = updateBookingPaidStatus(ctx, tx, true, i.Booking.ID, clinicianID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *repo) CountClinicianInvoices(ctx context.Context, clinicianID int) (int, error) {
	const query = `SELECT COUNT(*) FROM clinician_booking_invoice WHERE person_id = $1`
	var count int
	row := r.conn.QueryRow(ctx, query, clinicianID)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
