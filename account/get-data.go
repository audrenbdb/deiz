package account

import (
	"context"
	"github.com/audrenbdb/deiz"
)

func (u *GetDataUsecase) GetClinicianAccountData(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error) {
	return u.AccountDataGetter.GetClinicianAccount(ctx, clinicianID)
}

func (u *GetDataUsecase) GetClinicianAccountPublicData(ctx context.Context, clinicianID int) (deiz.ClinicianAccountPublicData, error) {
	acc, err := u.AccountDataGetter.GetClinicianAccount(ctx, clinicianID)
	if err != nil {
		return deiz.ClinicianAccountPublicData{}, err
	}
	return deiz.ClinicianAccountPublicData{
		Clinician:       acc.Clinician,
		StripePublicKey: acc.StripePublicKey,
		PublicMotives:   filterPublicMotives(acc.BookingMotives),
		RemoteAllowed:   acc.CalendarSettings.RemoteAllowed,
	}, nil
}

func filterPublicMotives(motives []deiz.BookingMotive) []deiz.BookingMotive {
	publicMotives := []deiz.BookingMotive{}
	for _, m := range motives {
		if m.Public {
			publicMotives = append(publicMotives, m)
		}
	}
	return publicMotives
}

type (
	accountGetter interface {
		GetClinicianAccount(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error)
	}
)

type GetDataUsecase struct {
	AccountDataGetter accountGetter
}
