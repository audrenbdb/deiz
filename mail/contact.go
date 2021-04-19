package mail

import (
	"context"
	"fmt"
	"github.com/audrenbdb/deiz"
)

func (m *Mailer) MailContactForm(ctx context.Context, to string, form deiz.ContactForm) error {
	details := contactEmailDetails{
		Name:    form.Name,
		Email:   form.Email,
		Message: form.HtmlMessage(),
	}
	template, err := m.htmlTemplate("contact-form.html", details)
	if err != nil {
		return err
	}
	plainBody := details.plainBody()
	return m.client.Send(createMail(mail{
		to:       to,
		from:     form.Email,
		subject:  "Nouvelle question de " + details.Name,
		template: template, plainBody: plainBody,
	}))
}

func (m *Mailer) MailGetInTouchForm(ctx context.Context, form deiz.GetInTouchForm) error {
	template, err := m.htmlTemplate("get-in-touch.html", form)
	if err != nil {
		return err
	}
	plainBody := getInTouchPlainBody(form)
	return m.client.Send(createMail(mail{
		to:        "contact@deiz.fr",
		from:      form.Email,
		subject:   "Demande de rappel",
		template:  template,
		plainBody: plainBody,
	}))
}

func getInTouchPlainBody(form deiz.GetInTouchForm) string {
	return fmt.Sprintf(`Deiz\n
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
}

type contactEmailDetails struct {
	Name    string
	Email   string
	Message string
}

func (details *contactEmailDetails) plainBody() string {
	return fmt.Sprintf(`Deiz\n
	Formulaire de contact\n
	\n
	De : %s| %s\n
	Message :\n
	%s\n
	\n
	Deiz\n
	Agenda pour thérapeutes\n
	https://deiz.fr
	`, details.Name, details.Email, details.Message)
}
