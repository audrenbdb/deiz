package mail

import (
	"fmt"
	"github.com/audrenbdb/deiz"
	"net/url"
	"time"
)

func (m *mailer) MailBookingReminder(b *deiz.Booking) error {
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

func (m *mailer) MailBookingToPatient(b *deiz.Booking) error {
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

func (m *mailer) MailBookingToClinician(b *deiz.Booking) error {
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

func (m *mailer) buildGCalendarLink(event gCalendarEvent) string {
	startStr := fmt.Sprintf("%d%02d%02dT%02d%02d00", event.start.Year(), event.start.Month(), event.start.Day(), event.start.Hour(), event.start.Minute())
	endStr := fmt.Sprintf("%d%02d%02dT%02d%02d00", event.end.Year(), event.end.Month(), event.end.Day(), event.end.Hour(), event.end.Minute())
	baseURL, _ := url.Parse("https://calendar.google.com")
	baseURL.Path += "calendar/event"
	params := url.Values{}
	params.Add("action", "TEMPLATE")
	params.Add("dates", fmt.Sprintf("%s/%s", startStr, endStr))
	params.Add("text", event.title)
	params.Add("details", event.details)
	params.Add("location", event.location)

	baseURL.RawQuery = params.Encode()
	return baseURL.String()
}

func buildGMapsLink(address string) string {
	baseURL, _ := url.Parse("https://www.google.com")
	baseURL.Path += "maps/search/"
	params := url.Values{}
	params.Add("api", "1")
	params.Add("query", address)

	baseURL.RawQuery = params.Encode()
	return baseURL.String()
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

func (m *mailer) getBookingEmailDetails(b *deiz.Booking, with string) bookingEmailDetails {
	return bookingEmailDetails{
		Clinician:   b.Clinician.FullName(),
		Patient:     b.Patient.FullName(),
		Phone:       b.Clinician.Phone,
		BookingDate: m.intl.Fr.FmtMMMEEEEd(b.Start),
		GCalendarLink: m.buildGCalendarLink(gCalendarEvent{
			start:    b.Start,
			end:      b.End,
			title:    fmt.Sprintf("Consultation avec %s", with),
			location: b.Address.ToString(),
		}),
		GMapsLink:   buildGMapsLink(b.Address.ToString()),
		CancelLink:  buildCancelURL(b.DeleteID).String(),
		AddressLine: b.Address.Line,
		AddressCity: fmt.Sprintf("%d %s", b.Address.PostCode, b.Address.City),
		Remote:      b.Remote,
	}
}
