package deiz

import (
	"context"
	"time"
)

type MailingService struct {
	BookingToClinicianMailer BookingToClinicianMailer
	BookingToPatientMailer   BookingToPatientMailer
}

type (
	BookingToClinicianMailer interface {
		MailBookingToClinician(ctx context.Context, b *Booking, tz *time.Location, gCalLink string) error
	}
	BookingToPatientMailer interface {
		MailBookingToPatient(ctx context.Context, b *Booking, tz *time.Location, gCalLink, gMapsLink, cancelURL string) error
	}
)
