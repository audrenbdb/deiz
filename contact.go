package deiz

import "context"

type ContactForm struct {
	Name    string `json:"name" validate:"min=2"`
	Email   string `json:"email" validate:"email"`
	Message string `json:"message" validate:"min=2"`
}

type (
	contactFormMailer interface {
		MailContactForm(ctx context.Context, to string, form ContactForm) error
	}
)

type (
	MailContactForm func(ctx context.Context, clinicianEmail string, form ContactForm) error
)

func mailContactFormFunc(mailer contactFormMailer) MailContactForm {
	return func(ctx context.Context, clinicianEmail string, form ContactForm) error {
		mailer.MailContactForm(ctx, clinicianEmail, form)
		return nil
	}
}
