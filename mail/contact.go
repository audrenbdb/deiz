package mail

import (
	"bytes"
	"context"
	"fmt"
	"github.com/audrenbdb/deiz"
	"strings"
)

func (m *mailer) MailContactForm(ctx context.Context, to string, form deiz.ContactForm) error {
	var emailBuffer bytes.Buffer
	emailData := struct {
		Name    string
		Email   string
		Message string
	}{
		Name:    form.Name,
		Email:   form.Email,
		Message: strings.Replace(form.Message, "\n", "<br>", -1),
	}
	err := m.tmpl.ExecuteTemplate(&emailBuffer, "contact-form.html", emailData)
	if err != nil {
		return err
	}
	plainBody := fmt.Sprintf(`Deiz\n
	Formulaire de contact\n
	\n
	De : %s| %s\n
	Message :\n
	%s\n
	\n
	Deiz\n
	Agenda pour thérapeutes\n
	https://deiz.fr
	`, form.Name, form.Email, form.Message)
	return m.sender.Send(ctx, createMail(to,
		form.Email, "Nouvelle question de "+emailData.Name,
		&emailBuffer, plainBody, nil))
}

func (m *mailer) MailGetInTouchForm(ctx context.Context, form deiz.GetInTouchForm) error {
	var emailBuffer bytes.Buffer
	err := m.tmpl.ExecuteTemplate(&emailBuffer, "get-in-touch.html", form)
	if err != nil {
		return err
	}
	plainBody := fmt.Sprintf(`Deiz\n
	Demande de rappel\n
	\n
	Madame ou monsieur %s souhaite être rappelé!\n
	\n
	Coordonnées :\n
	Nom : %s\n
	Email : %s\n
	Téléphone : %s\n
	Ville : %s\n
	Métier : %s\n
	\n
	Deiz\n
	Agenda pour thérapeutes\n
	https://deiz.fr`, form.Name, form.Name, form.Email, form.Phone, form.City, form.Job)
	return m.sender.Send(ctx, createMail("contact@deiz.fr", form.Email, "Demande de rappel", &emailBuffer, plainBody, nil))
}
