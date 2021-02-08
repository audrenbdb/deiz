package deiz_test

import (
	"context"
	"github.com/audrenbdb/deiz"
	"time"
)

type (
	mockBookingToClinicianMailer struct {
		err error
	}
	mockBookingToPatientMailer struct {
		err error
	}
)

func (r *mockBookingToClinicianMailer) MailBookingToClinician(ctx context.Context, b *deiz.Booking, tz *time.Location, gCalLink string) error {
	return r.err
}

func (r *mockBookingToPatientMailer) MailBookingToPatient(ctx context.Context, b *deiz.Booking, tz *time.Location, gCalLink, gMapsLink, cancelURL string) error {
	return r.err
}
