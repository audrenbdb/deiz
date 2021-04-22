package billing

import (
	"context"
	"github.com/audrenbdb/deiz"
	"time"
)

func (c *GetPeriodInvoicesUsecase) GetPeriodInvoices(ctx context.Context, start, end time.Time, clinicianID int) ([]deiz.BookingInvoice, error) {
	return c.getter.GetPeriodBookingInvoices(ctx, start, end, clinicianID)
}

type GetPeriodInvoicesUsecase struct {
	getter periodInvoicesGetter
}

func NewPeriodInvoicesGetterUsecase(getter periodInvoicesGetter) *GetPeriodInvoicesUsecase {
	return &GetPeriodInvoicesUsecase{
		getter: getter,
	}
}
