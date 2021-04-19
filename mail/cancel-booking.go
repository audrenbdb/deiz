package mail

import (
	"fmt"
	"github.com/audrenbdb/deiz"
)

func (m *mailer) MailCancelBookingToPatient(b *deiz.Booking) error {
	details := m.getCancelEmailDetails(b)
	template, err := m.htmlTemplate("cancelappointment-topatient.html", details)
	if err != nil {
		return err
	}
	return m.client.Send(createMail(mail{to: b.Patient.Email,
		from:     noReplyAddress,
		subject:  fmt.Sprintf("RDV du %s annulé", details.BookingDate),
		template: template, plainBody: details.plainBodyToPatient(),
	}))
}

func (m *mailer) MailCancelBookingToClinician(b *deiz.Booking) error {
	details := m.getCancelEmailDetails(b)
	template, err := m.htmlTemplate("cancelappointment-toclinician.html", details)
	if err != nil {
		return err
	}
	plainBody := details.plainBodyToClinician()
	return m.client.Send(createMail(mail{
		to:        b.Clinician.Email,
		from:      noReplyAddress,
		subject:   fmt.Sprintf("RDV du %s avec %s annulé", details.BookingDate, details.Name),
		template:  template,
		plainBody: plainBody,
	}))
}

func (m *mailer) getCancelEmailDetails(b *deiz.Booking) cancelEmailDetails {
	return cancelEmailDetails{
		BookingDate: m.intl.Fr.FmtMMMEEEEd(b.Start),
		Name:        b.Patient.Surname + " " + b.Patient.Name,
		Phone:       b.Patient.Phone,
		Email:       b.Patient.Email,
		Motive:      b.Motive.Name,
	}
}

func (details *cancelEmailDetails) plainBodyToClinician() string {
	return fmt.Sprintf(`Annulation\n\n
		RDV prévu %s annulé\n
		Motif %s\n
	Patient :\n
		%s\n
		%s\n
		%s\n
		Pour toute question, veuillez contacter le patient concerné.\n
		\n
		Deiz\n
		Application de gestion pour thérapeutes\n
	https://deiz.fr`,
		details.BookingDate,
		details.Motive,
		details.Name, details.Phone, details.Email)
}

func (details *cancelEmailDetails) plainBodyToPatient() string {
	return fmt.Sprintf(`Annulation\n\n
	La consultation du %s a été supprimée\n
	Pour toute question, veuillez contacter le clinicien concerné.\n
	\n
	Deiz\n
	Application de gestion pour thérapeutes\n
	https://deiz.fr`, details.BookingDate)
}

type cancelEmailDetails struct {
	BookingDate string
	Name        string
	Phone       string
	Email       string
	Motive      string
}
