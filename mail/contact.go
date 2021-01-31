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
