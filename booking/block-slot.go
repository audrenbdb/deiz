package booking

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type BlockSlotUsecase struct {
	Blocker bookingCreater
}

func (b *BlockSlotUsecase) BlockBookingSlot(ctx context.Context, slot *deiz.Booking, clinicianID int) error {
	if slot.Clinician.ID != clinicianID {
		return deiz.ErrorUnauthorized
	}
	return b.Blocker.CreateBooking(ctx, slot.SetBlocked())
}
