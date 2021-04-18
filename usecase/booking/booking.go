package booking

import (
	"context"
	"fmt"
	"github.com/audrenbdb/deiz"
	"net/url"
	"time"
)

type (
	ClinicianBookingsInTimeRangeGetter interface {
		GetClinicianBookingsInTimeRange(ctx context.Context, start, end time.Time, clinicianID int) ([]deiz.Booking, error)
	}
	BookingsInTimeRangeGetter interface {
		GetBookingsInTimeRange(ctx context.Context, start, end time.Time) ([]deiz.Booking, error)
	}

	ToClinicianMailer interface {
		MailBookingToClinician(ctx context.Context, b *deiz.Booking, gCalLink string) error
	}
	BookingReminderMailer interface {
		MailBookingReminder(ctx context.Context, b *deiz.Booking, gCalLink, gMapsLink, cancelURL string) error
	}
	ToPatientMailer interface {
		MailBookingToPatient(ctx context.Context, b *deiz.Booking, gCalLink, gMapsLink, cancelURL string) error
	}
	CancelBookingToPatientMailer interface {
		MailCancelBookingToPatient(ctx context.Context, b *deiz.Booking) error
	}
	CancelBookingToClinicianMailer interface {
		MailCancelBookingToClinician(ctx context.Context, b *deiz.Booking) error
	}
	PatientCreater interface {
		CreatePatient(ctx context.Context, p *deiz.Patient, clinicianID int) error
	}
	PatientGetterByEmail interface {
		GetPatientByEmail(ctx context.Context, email string, clinicianID int) (deiz.Patient, error)
	}
	AddressGetterByID interface {
		GetAddressByID(ctx context.Context, addressID int) (deiz.Address, error)
	}
)

func (u *Usecase) MailToClinician(ctx context.Context, b *deiz.Booking, clinicianID int) error {
	return u.ToClinicianMailer.MailBookingToClinician(ctx,
		b, u.GCalendarLinkBuilder.BuildGCalendarLink(
			b.Start.In(u.loc), b.End.In(u.loc),
			fmt.Sprintf("Consultation avec %s %s", b.Patient.Surname, b.Patient.Name),
			b.Address.ToString(), b.Note,
		))
}

func (u *Usecase) MailToPatient(ctx context.Context, b *deiz.Booking, clinicianID int) error {
	cancelURL := createCancelURL(b.DeleteID)
	return u.ToPatientMailer.MailBookingToPatient(ctx, b,
		createGCalendarBookingLinkToPatient(b, u.loc, u.GCalendarLinkBuilder),
		u.GMapsLinkBuilder.BuildGMapsLink(b.Address.ToString()),
		cancelURL.String(),
	)
}

//get time in 48h and resets its mn value.
//IE now => 2020-01-01 10:30
//Reminder => 2020-01-03 10:00
func getReminderStartTime() time.Time {
	anchor := time.Now().AddDate(0, 0, 2).UTC()
	return time.Date(anchor.Year(), anchor.Month(), anchor.Day(), anchor.Hour(), 0, 0, 0, time.UTC)
}

func (u *Usecase) SendReminders(ctx context.Context) error {
	start := getReminderStartTime()
	end := start.Add(time.Hour * time.Duration(1)).UTC()

	bookings, err := u.BookingsInTimeRangeGetter.GetBookingsInTimeRange(ctx, start, end)
	if err != nil {
		return fmt.Errorf("unable to get bookings in time range: %s", err)
	}
	for _, b := range bookings {
		gCalLink := createGCalendarBookingLinkToPatient(&b, u.loc, u.GCalendarLinkBuilder)
		err := u.BookingReminderMailer.MailBookingReminder(ctx,
			&b, gCalLink, u.GMapsLinkBuilder.BuildGMapsLink(b.Address.ToString()), createCancelURL(b.DeleteID).String())
		if err != nil {
			return err
		}
	}
	return nil
}

type (
	clinicianBookingsInTimeRangeGetter interface {
		GetClinicianBookingsInTimeRange(ctx context.Context, start, end time.Time, clinicianID int) ([]deiz.Booking, error)
	}
	bookingCreater interface {
		CreateBooking(ctx context.Context, b *deiz.Booking) error
	}
	gCalendarLinkBuilder interface {
		BuildGCalendarLink(start, end time.Time, subject, addressStr, details string) string
	}
	bookingMailer interface {
		MailBookingToClinician(ctx context.Context, b *deiz.Booking, gCalLink string) error
		MailBookingToPatient(ctx context.Context, b *deiz.Booking, gCalLink, gMapsLink, cancelURL string) error
	}
)

func bookingsOverlap(booking1, booking2 *deiz.Booking) bool {
	return booking1.Start.Before(booking2.End) && booking2.Start.Before(booking1.End)
}

func bookingSlotAvailable(ctx context.Context, b *deiz.Booking, getter clinicianBookingsInTimeRangeGetter) (bool, error) {
	bookings, err := getter.GetClinicianBookingsInTimeRange(ctx, b.Start, b.End, b.Clinician.ID)
	if err != nil {
		return false, err
	}
	for _, booking := range bookings {
		if bookingsOverlap(b, &booking) {
			return false, nil
		}
	}
	return true, nil
}

func createGCalendarBookingLinkToPatient(b *deiz.Booking, loc *time.Location, builder gCalendarLinkBuilder) string {
	return createGCalendarBookingLink(b, loc, builder,
		fmt.Sprintf("Consultation avec %s %s", b.Clinician.Surname, b.Clinician.Name))
}

func createGCalendarBookingLinkToClinician(b *deiz.Booking, loc *time.Location, builder gCalendarLinkBuilder) string {
	return createGCalendarBookingLink(b, loc, builder,
		fmt.Sprintf("Consultation avec %s %s", b.Patient.Surname, b.Patient.Name))
}

func createGCalendarBookingLink(b *deiz.Booking, loc *time.Location, builder gCalendarLinkBuilder, subject string) string {
	return builder.BuildGCalendarLink(
		b.Start.In(loc),
		b.End.In(loc),
		subject,
		b.Address.ToString(),
		"",
	)
}

func mailBookingToClinician(
	ctx context.Context,
	b *deiz.Booking,
	loc *time.Location,
	mailer bookingMailer,
	builder gCalendarLinkBuilder) error {
	gCalendarLink := createGCalendarBookingLinkToClinician(b, loc, builder)
	return mailer.MailBookingToClinician(ctx, b, gCalendarLink)
}

func mailBookingToPatient(
	ctx context.Context,
	b *deiz.Booking,
	loc *time.Location,
	mailer bookingMailer,
	gCalendarLinkBuilder gCalendarLinkBuilder,
	gMapsLinkBuilder gMapsLinkBuilder,
) error {
	gCalendarLink := createGCalendarBookingLinkToPatient(b, loc, gCalendarLinkBuilder)
	gMapsLink := gMapsLinkBuilder.BuildGMapsLink(b.Address.ToString())
	cancelURL := createCancelURL(b.DeleteID).String()
	return mailer.MailBookingToPatient(ctx, b, gCalendarLink, gMapsLink, cancelURL)
}

func createCancelURL(deleteID string) *url.URL {
	cancelURL, _ := url.Parse("https://deiz.fr")
	cancelURL.Path += "bookings/delete"
	params := url.Values{}
	params.Add("id", deleteID)
	cancelURL.RawQuery = params.Encode()
	return cancelURL
}
