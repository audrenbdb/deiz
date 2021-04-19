package booking

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type SlotBlocker struct {
	blocker bookingCreater
}

type SlotBlockerDeps struct {
	Blocker bookingCreater
}

func NewSlotBlockerUsecase(deps SlotBlockerDeps) *SlotBlocker {
	return &SlotBlocker{
		blocker: deps.Blocker,
	}
}

func (b *SlotBlocker) BlockBookingSlot(ctx context.Context, slot *deiz.Booking, clinicianID int) error {
	if slot.Clinician.ID != clinicianID {
		return deiz.ErrorUnauthorized
	}
	return b.blocker.CreateBooking(ctx, slot.SetBlocked())
}
