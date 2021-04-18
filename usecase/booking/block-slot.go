package booking

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type slotBlocker struct {
	blocker bookingCreater
}

func NewSlotBlockerUsecase(blocker bookingCreater) *slotBlocker {
	return &slotBlocker{
		blocker: blocker,
	}
}

func (b *slotBlocker) BlockBookingSlot(ctx context.Context, slot *deiz.Booking, clinicianID int) error {
	if slot.Clinician.ID != clinicianID {
		return deiz.ErrorUnauthorized
	}
	return b.blocker.CreateBooking(ctx, slot.SetBlocked())
}
