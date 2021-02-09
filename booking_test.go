package deiz_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	mockBookingStorer struct {
		err error
	}
)

func (r *mockBookingStorer) StoreBooking(ctx context.Context, b *deiz.Booking) error {
	return r.err
}

func TestRegisterBooking(t *testing.T) {
	var tests = []struct {
		description string

		bookingStorer            *mockBookingStorer
		bookingToClinicianMailer *mockBookingToClinicianMailer
		bookingToPatientMailer   *mockBookingToPatientMailer
		googleCalendarLinkMaker  *mockGCalendarLinkMaker
		googleMapsLinkMaker      *mockGMapsLinkMaker

		inBooking         *deiz.Booking
		inClinicianID     int
		inClinicianTz     string
		inNotifyClinician bool
		inNotifyPatient   bool

		outError error
	}{
		{
			description:   "should fail to authenticate request as emitted from same clinician",
			inBooking:     &deiz.Booking{Clinician: deiz.Clinician{ID: 0}},
			inClinicianID: 1,
			outError:      deiz.ErrorUnauthorized,
		},
		{
			description:   "should fail to store booking into database",
			bookingStorer: &mockBookingStorer{err: errors.New("fail to store booking into db")},
			inBooking:     &deiz.Booking{Clinician: deiz.Clinician{ID: 1}},
			inClinicianID: 1,
			outError:      errors.New("fail to store booking into db"),
		},
		{
			description:   "should fail to load timezone location",
			bookingStorer: &mockBookingStorer{},
			inBooking:     &deiz.Booking{Clinician: deiz.Clinician{ID: 1}},
			inClinicianID: 1,
			inClinicianTz: "fail",
			outError:      errors.New("unknown time zone fail"),
		},
		{
			description:              "should fail to notify clinician of the booking registration through email",
			bookingStorer:            &mockBookingStorer{},
			bookingToClinicianMailer: &mockBookingToClinicianMailer{err: errors.New("fail to mail to clinician")},
			googleCalendarLinkMaker:  &mockGCalendarLinkMaker{},
			inBooking:                &deiz.Booking{Clinician: deiz.Clinician{ID: 1}},
			inNotifyClinician:        true,
			inClinicianID:            1,
			outError:                 errors.New("fail to mail to clinician"),
		},
		{
			description:             "should fail to notify patient of the booking registration through email",
			bookingStorer:           &mockBookingStorer{},
			bookingToPatientMailer:  &mockBookingToPatientMailer{err: errors.New("fail to mail to patient")},
			googleCalendarLinkMaker: &mockGCalendarLinkMaker{},
			googleMapsLinkMaker:     &mockGMapsLinkMaker{},
			inBooking:               &deiz.Booking{Clinician: deiz.Clinician{ID: 1}},
			inNotifyPatient:         true,
			inClinicianID:           1,
			outError:                errors.New("fail to mail to patient"),
		},
		{
			description:              "should succeed to register and send emails to both patient and clinician",
			bookingStorer:            &mockBookingStorer{},
			bookingToPatientMailer:   &mockBookingToPatientMailer{},
			bookingToClinicianMailer: &mockBookingToClinicianMailer{},
			googleCalendarLinkMaker:  &mockGCalendarLinkMaker{},
			googleMapsLinkMaker:      &mockGMapsLinkMaker{},
			inBooking:                &deiz.Booking{Clinician: deiz.Clinician{ID: 1}},
			inNotifyPatient:          true,
			inNotifyClinician:        true,
			inClinicianID:            1,
			outError:                 nil,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			r := deiz.Repo{
				Booking: deiz.BookingRepo{
					Storer: test.bookingStorer,
				},
				Mailing: deiz.MailingService{
					BookingToClinicianMailer: test.bookingToClinicianMailer,
					BookingToPatientMailer:   test.bookingToPatientMailer,
				},
				GoogleCalendar: deiz.GoogleCalendarService{
					LinkMaker: test.googleCalendarLinkMaker,
				},
				GoogleMaps: deiz.GoogleMapsService{
					GoogleMapsLinkMaker: test.googleMapsLinkMaker,
				},
			}
			err := r.RegisterBooking(context.Background(), test.inBooking, test.inClinicianID, test.inClinicianTz, test.inNotifyClinician, test.inNotifyPatient)
			assert.Equal(t, test.outError, err, fmt.Sprintf("got %s, want %s", err, test.outError))

		})
	}
}
