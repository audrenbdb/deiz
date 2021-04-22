package mail

import (
	"fmt"
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/gcal"
	"github.com/audrenbdb/deiz/gmaps"
	"net/url"
	"time"
)

func (m *Mailer) MailBookingReminder(b *deiz.Booking) error {
	details := m.getBookingEmailDetails(b, b.Clinician.FullName())
	template, err := m.htmlTemplate("booking-reminder.html", details)
	if err != nil {
		return err
	}
	plainBody := details.plainBodyToPatient()
	return m.client.Send(createMail(mail{
		to:       b.Patient.Email,
		from:     noReplyAddress,
		subject:  "Rappel de rdv : " + details.BookingDate,
		template: template, plainBody: plainBody,
	}))
}

func (m *Mailer) MailBookingToPatient(b *deiz.Booking) error {
	details := m.getBookingEmailDetails(b, b.Clinician.FullName())
	template, err := m.htmlTemplate("confirmbooking-topatient.html", details)
	if err != nil {
		return err
	}
	plainBody := details.plainBodyToPatient()
	return m.client.Send(createMail(mail{to: b.Patient.Email,
		from:     noReplyAddress,
		subject:  "RDV confirmé " + details.BookingDate,
		template: template, plainBody: plainBody,
	}))
}

func (m *Mailer) MailBookingToClinician(b *deiz.Booking) error {
	details := m.getBookingEmailDetails(b, b.Patient.FullName())
	template, err := m.htmlTemplate("confirmbooking-toclinician.html", details)
	if err != nil {
		return err
	}
	plainBody := details.plainBodyToClinician()
	return m.client.Send(createMail(mail{to: b.Clinician.Email,
		from:     noReplyAddress,
		subject:  fmt.Sprintf("RDV confirmé avec %s %s", details.Patient, details.BookingDate),
		template: template, plainBody: plainBody,
	}))
}

type gCalendarEvent struct {
	start    time.Time
	end      time.Time
	title    string
	location string
	details  string
}

func buildCancelURL(deleteID string) *url.URL {
	cancelURL, _ := url.Parse("https://deiz.fr")
	cancelURL.Path += "bookings/delete"
	params := url.Values{}
	params.Add("id", deleteID)
	cancelURL.RawQuery = params.Encode()
	return cancelURL
}

type bookingEmailDetails struct {
	Clinician     string
	Patient       string
	Phone         string
	BookingDate   string
	GCalendarLink string
	GMapsLink     string
	CancelLink    string
	AddressLine   string
	AddressCity   string
	Remote        bool
	Motive        string
	Email         string
}

func (details *bookingEmailDetails) plainBodyToPatient() string {
	return fmt.Sprintf(`Deiz\n
	RDV confirmé\n
	Avec %s\n
	%s\n
	%s\n
	\n
	Ajouter à Google Calendar : %s\n
	Annuler : %s\n
	\n
	%s\n
	%s\n
	Itinéraire : %s\n
	\n
	Deiz\n
	Agenda pour thérapeutes\n
	https://deiz.fr
	`, details.Clinician, details.Phone, details.BookingDate, details.GCalendarLink,
		details.CancelLink, details.AddressLine, details.AddressCity, details.GMapsLink)
}

func (details *bookingEmailDetails) plainBodyToClinician() string {
	return fmt.Sprintf(`Deiz\n
	RDV confirmé\n
	%s\n
	%s\n
	\n
	%s\n
	%s\n
	%s\n
	Ajouter à Google Calendar : %s\n
	\n
	Deiz\n
	Agenda pour thérapeutes\n
	https://deiz.fr
	`, details.BookingDate, details.Motive, details.Patient, details.Phone, details.Email, details.GCalendarLink)
}

func (m *Mailer) getBookingEmailDetails(b *deiz.Booking, with string) bookingEmailDetails {
	return bookingEmailDetails{
		Clinician:   b.Clinician.FullName(),
		Patient:     b.Patient.FullName(),
		Phone:       b.Clinician.Phone,
		BookingDate: m.intl.Fr.FmtMMMEEEEd(b.Start),
		GCalendarLink: gcal.NewLink(gcal.Event{
			Start:    b.Start,
			End:      b.End,
			Title:    fmt.Sprintf("Consultation avec %s", with),
			Location: b.Address.ToString(),
		}),
		GMapsLink:   gmaps.CreateLink(b.Address.ToString()),
		CancelLink:  buildCancelURL(b.DeleteID).String(),
		AddressLine: b.Address.Line,
		AddressCity: fmt.Sprintf("%d %s", b.Address.PostCode, b.Address.City),
		Remote:      b.Remote,
	}
}
