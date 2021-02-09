package account

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	ClinicianGetterByEmail interface {
		GetClinicianByEmail(ctx context.Context, email string) (deiz.Clinician, error)
	}
)
