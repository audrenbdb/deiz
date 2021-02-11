package patient

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	Searcher interface {
		SearchPatient(ctx context.Context, search string, clinicianID int) ([]deiz.Patient, error)
	}
)

func (u *Usecase) SearchPatient(ctx context.Context, search string, clinicianID int) ([]deiz.Patient, error) {
	return u.Searcher.SearchPatient(ctx, search, clinicianID)
}
