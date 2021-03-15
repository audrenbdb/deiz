package contact

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	ClinicianGetterByID interface {
		GetClinicianByID(ctx context.Context, clinicianID int) (deiz.Clinician, error)
	}
	ContactFormMailer interface {
		MailContactForm(ctx context.Context, to string, form deiz.ContactForm) error
	}
	GetInTouchMailer interface {
		MailGetInTouchForm(ctx context.Context, form deiz.GetInTouchForm) error
	}
)

type repo interface {
	ClinicianGetterByID
}

type mail interface {
	ContactFormMailer
	GetInTouchMailer
}

type Usecase struct {
	ClinicianGetter   ClinicianGetterByID
	ContactFormMailer ContactFormMailer
	GetInTouchMailer  GetInTouchMailer
}

func NewUsecase(repo repo, mail mail) *Usecase {
	return &Usecase{
		ClinicianGetter:   repo,
		ContactFormMailer: mail,
		GetInTouchMailer:  mail,
	}
}

func (u *Usecase) SendContactFormToClinician(ctx context.Context, f deiz.ContactForm) error {
	if f.ClinicianID == 0 || len(f.Name) < 2 || len(f.Message) < 2 || len(f.Email) < 6 {
		return deiz.ErrorUnauthorized
	}
	c, err := u.ClinicianGetter.GetClinicianByID(ctx, f.ClinicianID)
	if err != nil {
		return err
	}
	return u.ContactFormMailer.MailContactForm(ctx, c.Email, f)
}

func (u *Usecase) SendGetInTouchForm(ctx context.Context, f deiz.GetInTouchForm) error {
	if len(f.Email) < 6 || len(f.Phone) < 10 || len(f.Name) < 2 || len(f.City) < 2 || len(f.Job) < 2 {
		return deiz.ErrorUnauthorized
	}
	return u.GetInTouchMailer.MailGetInTouchForm(ctx, f)
}
