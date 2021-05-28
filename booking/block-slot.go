package booking

import (
	"context"
	"github.com/audrenbdb/deiz"
	"time"
)

type BlockSlotUsecase struct {
	Blocker bookingCreater
	Deleter blockedSlotDeleter
}

type blockedSlotDeleter interface {
	DeleteBlockedBookingPrior(ctx context.Context, d time.Time) error
}

func (b *BlockSlotUsecase) BlockBookingSlots(ctx context.Context, slots []*deiz.Booking, cred deiz.Credentials) error {
	if areBookingsInvalid(slots, cred.UserID) {
		return deiz.ErrorUnauthorized
	}
	for _, slot := range slots {
		err := b.Blocker.CreateBooking(ctx, slot)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BlockSlotUsecase) DeletePastBlockedBookingSlot(ctx context.Context) error {
	return b.Deleter.DeleteBlockedBookingPrior(ctx, time.Now())
}
