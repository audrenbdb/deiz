package contact

import (
	"context"
	"errors"
	"fmt"
	"github.com/audrenbdb/deiz/email"
	"strings"
)

type sendContactForm = func(ctx context.Context, form contactForm) error

func sendContactFormFn(getClinician getClinicianByID, sendEmail email.Send) sendContactForm {
	return func(ctx context.Context, form contactForm) error {
		if form.invalid() {
			return errors.New("could not validate form provided")
		}
		c, err := getClinician(ctx, form.ClinicianID)
		if err != nil {
			return err
		}
		m := createContactFormEmail(c, form)
		return sendEmail(m)
	}
}

type sendGetInTouchForm = func(form getInTouchForm) error

func sendGetInTouchFormFn(sendEmail email.Send) sendGetInTouchForm {
	return func(form getInTouchForm) error {
		if form.invalid() {
			return errors.New("could not validate get in touch form provided")
		}
		return sendEmail(createGetInTouchFormEmail(form))
	}
}

type contactFormDataToBind struct {
	Name    string
	Email   string
	Message string
}

func createContactFormEmail(c clinician, f contactForm) email.HTMLMail {
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
	`, f.Name, f.Email, f.Message)
	return email.HTMLMail{
		Header: email.Header{
			To:      c.email,
			From:    f.Email,
			Subject: "Nouvelle question de " + f.Name,
		},
		Body: email.Body{
			HTMLFileName: "contact-form.html",
			DataToBind: contactFormDataToBind{
				Name:    f.Name,
				Email:   f.Email,
				Message: replaceLineJumpsWithBRTags(f.Message),
			},
			Plain: plainBody,
		},
	}
}

func createGetInTouchFormEmail(f getInTouchForm) email.HTMLMail {
	plainBody := fmt.Sprintf(`Deiz\n
		Demande de rappel\n
		\n
		Madame ou monsieur %s souhaite être rappelé!\n
		\n
		Coordonnées :\n
		Nom : %s\n
		EmailUpdater : %s\n
		Téléphone : %s\n
		Ville : %s\n
		Métier : %s\n
		\n
		Deiz\n
		Agenda pour thérapeutes\n
		https://deiz.fr
	`, f.Name, f.Name, f.Email, f.Phone, f.City, f.Job)
	return email.HTMLMail{
		Header: email.Header{
			To:      deizContactAddress,
			From:    f.Email,
			Subject: "Demande de rappel",
		},
		Body: email.Body{
			HTMLFileName: "get-in-touch.html",
			DataToBind:   f,
			Plain:        plainBody,
		},
	}
}

func replaceLineJumpsWithBRTags(txt string) string {
	return strings.Replace(txt, "\n", "<br>", -1)
}

func (f *getInTouchForm) valid() bool {
	return len(f.Email) >= 6 && len(f.Phone) >= 10 && len(f.Name) >= 2 && len(f.City) >= 2 && len(f.Job) >= 2
}

func (f *getInTouchForm) invalid() bool {
	return !f.valid()
}

func (f *contactForm) valid() bool {
	return f.ClinicianID != 0 && len(f.Name) >= 2 && len(f.Message) >= 2 && len(f.Email) >= 6
}

func (f *contactForm) invalid() bool {
	return !f.valid()
}
