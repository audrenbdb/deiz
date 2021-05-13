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

func (b *BlockSlotUsecase) BlockBookingSlot(ctx context.Context, slot *deiz.Booking, clinicianID int) error {
	return b.blockSingleSlot(ctx, slot, clinicianID)
}

func (b *BlockSlotUsecase) BlockBookingSlotList(ctx context.Context, slots []*deiz.Booking, cred deiz.Credentials) error {
	for _, slot := range slots {
		err := b.blockSingleSlot(ctx, slot, cred.UserID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BlockSlotUsecase) blockSingleSlot(ctx context.Context, slot *deiz.Booking, clinicianID int) error {
	if slot.Clinician.ID != clinicianID {
		return deiz.ErrorUnauthorized
	}
	return b.Blocker.CreateBooking(ctx, slot.SetBlocked())
}

func (b *BlockSlotUsecase) DeletePastBlockedBookingSlot(ctx context.Context) error {
	return b.Deleter.DeleteBlockedBookingPrior(ctx, time.Now())
}
