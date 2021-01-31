package mail

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"testing"
	"text/template"
	"time"
)

func TestMailBookingToClinician(t *testing.T) {
	//Commented to prevent fail email send while testing
	/*
		m := NewService(template.Must(template.ParseGlob("templates/*.html")), NewGmailClient())
		err := m.MailBookingToClinician(context.Background(), &deiz.Booking{
			Patient: deiz.Patient{
				Email:   "audren.bdb@gmail.com",
				Name:    "PATIENT",
				Surname: "patient",
				Phone:   "TESTPHONE",
			},
			Clinician: deiz.Clinician{
				Email: "martin.nicolas7@gmail.com",
			},
			Start: time.Now(),
			Motive: deiz.BookingMotive{
				Name: "Test",
			},
		}, "testlink.fr", time.UTC)
		assert.NoError(t, err)
	*/
}

func TestMailBookingToPatient(t *testing.T) {
	//Commented to prevent fail email send while testing

	m := NewService(template.Must(template.ParseGlob("templates/*.html")), NewGmailClient())
	err := m.MailBookingToPatient(context.Background(), &deiz.Booking{
		Patient: deiz.Patient{Email: "martin.nicolas7@gmail.com"},
		Clinician: deiz.Clinician{
			Surname: "Clinician",
			Name:    "CLINICIAN",
			Phone:   "TESTPHONE",
			Email:   "audren.bdb@gmail.com",
		},
		Address: deiz.Address{
			Line:     "Test LINE",
			PostCode: 10000,
			City:     "Test city",
		},
		Start:  time.Now(),
		Remote: true,
	}, time.UTC, "testgcalendar.fr", "testgmaps.fr", "testcancel.fr")
	assert.NoError(t, err)

}
