package mail

import (
	"bytes"
	"context"
	"fmt"
	"github.com/audrenbdb/deiz"
	"time"
)

func (m *mailer) MailCancelBookingToPatient(ctx context.Context, b *deiz.Booking, tz *time.Location) error {
	var emailBuffer bytes.Buffer
	emailData := struct {
		BookingDateStr string
	}{
		BookingDateStr: b.Start.In(tz).Format("02/01/2006 à 15h04"),
	}
	err := m.tmpl.ExecuteTemplate(&emailBuffer, "cancelappointment-topatient.html", emailData)
	if err != nil {
		return err
	}
	plainBody := fmt.Sprintf(`Annulation\n\n
	La consultation du %s a été supprimée\n
	Pour toute question, veuillez contacter le clinicien concerné.\n
	\n
	Deiz\n
	Application de gestion pour thérapeutes\n
	https://deiz.fr`, emailData.BookingDateStr)
	return m.sender.Send(ctx, createMail(b.Patient.Email,
		b.Clinician.Email,
		fmt.Sprintf("RDV du %s annulé", emailData.BookingDateStr),
		&emailBuffer, plainBody,
		nil,
	))
}

func (m *mailer) MailCancelBookingToClinician(ctx context.Context, b *deiz.Booking, tz *time.Location) error {
	var emailBuffer bytes.Buffer
	emailData := struct {
		BookingDateStr string
		Name           string
		Phone          string
		Email          string
		Motive         string
	}{
		BookingDateStr: b.Start.In(tz).Format("02/01/2006 à 15h04"),
		Name:           b.Patient.Surname + " " + b.Patient.Name,
		Phone:          b.Patient.Phone,
		Email:          b.Patient.Email,
		Motive:         b.Motive.Name,
	}
	err := m.tmpl.ExecuteTemplate(&emailBuffer, "cancelappointment-toclinician.html", emailData)
	if err != nil {
		return err
	}
	plainBody := fmt.Sprintf(`Annulation\n\n
	RDV prévu %s annulé\n
	Motif %s\n
	Patient :\n
	%s\n
	%s\n
	%s\n
	Pour toute question, veuillez contacter le clinicien concerné.\n
	\n
	Deiz\n
	Application de gestion pour thérapeutes\n
	https://deiz.fr`, emailData.BookingDateStr, emailData.Motive, emailData.Name, emailData.Phone, emailData.Email)
	return m.sender.Send(ctx, createMail(b.Patient.Email,
		b.Clinician.Email,
		fmt.Sprintf("RDV du %s avec %s annulé", emailData.BookingDateStr, emailData.Name),
		&emailBuffer, plainBody,
		nil,
	))
}

func (m *mailer) MailBookingToPatient(ctx context.Context, b *deiz.Booking, clinicianTz *time.Location, gCalendarLink, gMapsLink, cancelURL string) error {
	var emailBuffer bytes.Buffer
	emailData := struct {
		Name           string
		Phone          string
		BookingDateStr string
		GCalendarLink  string
		GMapsLink      string
		CancelLink     string
		AddressLine    string
		AddressCity    string
		Remote         bool
	}{
		Name:           b.Clinician.Surname + " " + b.Clinician.Name,
		Phone:          b.Clinician.Phone,
		BookingDateStr: b.Start.In(clinicianTz).Format("le 02/01/2006 à 15h04"),
		GCalendarLink:  gCalendarLink,
		GMapsLink:      gMapsLink,
		CancelLink:     cancelURL,
		AddressLine:    b.Address.Line,
		AddressCity:    fmt.Sprintf("%d %s", b.Address.PostCode, b.Address.City),
		Remote:         b.Remote,
	}

	err := m.tmpl.ExecuteTemplate(&emailBuffer, "confirmbooking-topatient.html", emailData)
	if err != nil {
		return err
	}
	plainBody := fmt.Sprintf(`Deiz\n
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
	`, emailData.Name, emailData.Phone, emailData.BookingDateStr, emailData.GCalendarLink,
		emailData.CancelLink, emailData.AddressLine, emailData.AddressCity, emailData.GMapsLink)
	return m.sender.Send(ctx, createMail(b.Patient.Email,
		b.Clinician.Email,
		"RDV confirmé "+emailData.BookingDateStr,
		&emailBuffer, plainBody,
		nil,
	))
}

func (m *mailer) MailBookingToClinician(ctx context.Context, b *deiz.Booking, clinicianTz *time.Location, gCalendarLink string) error {
	var emailBuffer bytes.Buffer
	emailData := struct {
		Name           string
		Phone          string
		BookingDateStr string
		Motive         string
		Email          string
		GCalendarLink  string
		Remote         bool
	}{
		Name:           b.Patient.Surname + " " + b.Patient.Name,
		Phone:          b.Patient.Phone,
		BookingDateStr: b.Start.In(clinicianTz).Format("le 02/01/2006 à 15h04"),
		Motive:         b.Motive.Name,
		Email:          b.Patient.Email,
		GCalendarLink:  gCalendarLink,
		Remote:         b.Remote,
	}
	err := m.tmpl.ExecuteTemplate(&emailBuffer, "confirmbooking-toclinician.html", emailData)
	if err != nil {
		return err
	}
	plainBody := fmt.Sprintf(`Deiz\n
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
	`, emailData.BookingDateStr, emailData.Motive, emailData.Name, emailData.Phone, emailData.Email, emailData.GCalendarLink)
	return m.sender.Send(ctx, createMail(b.Clinician.Email,
		b.Patient.Email,
		fmt.Sprintf("RDV confirmé avec %s %s", emailData.Name, emailData.BookingDateStr),
		&emailBuffer, plainBody,
		nil,
	))
}
