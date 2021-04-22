package officehours

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	creater interface {
		CreateOfficeHours(ctx context.Context, h *deiz.OfficeHours, clinicianID int) error
	}
	deleter interface {
		DeleteOfficeHours(ctx context.Context, hoursID, clinicianID int) error
	}
)

type Usecase struct {
	Creater creater
	Deleter deleter
}

func (u *Usecase) AddOfficeHours(ctx context.Context, h *deiz.OfficeHours, clinicianID int) error {
	if h.IsInvalid() {
		return deiz.ErrorStructValidation
	}
	return u.Creater.CreateOfficeHours(ctx, h, clinicianID)
}

func (u *Usecase) RemoveOfficeHours(ctx context.Context, hoursID, clinicianID int) error {
	return u.Deleter.DeleteOfficeHours(ctx, hoursID, clinicianID)
}
