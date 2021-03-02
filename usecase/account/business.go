package account

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	ClinicianBusinessUpdater interface {
		UpdateClinicianBusiness(ctx context.Context, b *deiz.Business, clinicianID int) error
	}
	TaxExemptionCodesGetter interface {
		GetTaxExemptionCodes(ctx context.Context) ([]deiz.TaxExemption, error)
	}
)

func (u *Usecase) EditClinicianBusiness(ctx context.Context, b *deiz.Business, clinicianID int) error {
	return u.BusinessUpdater.UpdateClinicianBusiness(ctx, b, clinicianID)
}

func (u *Usecase) GetTaxExemptionCodes(ctx context.Context) ([]deiz.TaxExemption, error) {
	return u.TaxExemptionCodesGetter.GetTaxExemptionCodes(ctx)
}
