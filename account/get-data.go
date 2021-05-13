package account

import (
	"context"
	"github.com/audrenbdb/deiz"
)

//GetClinicianAccountData retrieves data about a clinician account for a client application to function properly
func (u *GetDataUsecase) GetClinicianAccountData(ctx context.Context, cred deiz.Credentials) (deiz.ClinicianAccount, error) {
	return u.AccountDataGetter.GetClinicianAccount(ctx, cred.UserID)
}

func (u *GetDataUsecase) GetClinicianAccountPublicData(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error) {
	acc, err := u.AccountDataGetter.GetClinicianAccount(ctx, clinicianID)
	if err != nil {
		return deiz.ClinicianAccount{}, err
	}
	return deiz.ClinicianAccount{
		Clinician:       acc.Clinician,
		StripePublicKey: acc.StripePublicKey,
		BookingMotives:  filterPublicMotives(acc.BookingMotives),
		CalendarSettings: deiz.CalendarSettings{
			RemoteAllowed:     acc.CalendarSettings.RemoteAllowed,
			NewPatientAllowed: acc.CalendarSettings.NewPatientAllowed,
		},
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
