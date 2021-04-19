package booking

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	bookingGetter interface {
		GetBookingByDeleteID(ctx context.Context, deleteID string) (deiz.Booking, error)
		GetBookingByID(ctx context.Context, bookingID int) (deiz.Booking, error)
	}
	bookingDeleter interface {
		DeleteBooking(ctx context.Context, bookingID, clinicianID int) error
	}
	cancelMailer interface {
		MailCancelBookingToClinician(b *deiz.Booking) error
		MailCancelBookingToPatient(b *deiz.Booking) error
	}
)

type SlotDeleter struct {
	bookingGetter  bookingGetter
	bookingDeleter bookingDeleter
	cancelMailer   cancelMailer
}

type SlotDeleterDeps struct {
	BookingGetter  bookingGetter
	BookingDeleter bookingDeleter
	CancelMailer   cancelMailer
}

func NewSlotDeleterUsecase(deps SlotDeleterDeps) *SlotDeleter {
	return &SlotDeleter{
		bookingGetter:  deps.BookingGetter,
		bookingDeleter: deps.BookingDeleter,
		cancelMailer:   deps.CancelMailer,
	}
}

func (d *SlotDeleter) DeleteBlockedSlot(ctx context.Context, bookingID, clinicianID int) error {
	return d.bookingDeleter.DeleteBooking(ctx, bookingID, clinicianID)
}

func (d *SlotDeleter) DeletePreRegisteredSlot(ctx context.Context, bookingID, clinicianID int) error {
	return d.bookingDeleter.DeleteBooking(ctx, bookingID, clinicianID)
}

func (d *SlotDeleter) DeleteBookedSlotFromPatient(ctx context.Context, deleteID string) error {
	booking, err := d.bookingGetter.GetBookingByDeleteID(ctx, deleteID)
	if err != nil {
		return err
	}
	if err := d.bookingDeleter.DeleteBooking(ctx, booking.ID, booking.Clinician.ID); err != nil {
		return err
	}
	return d.cancelMailer.MailCancelBookingToClinician(&booking)
}

func (d *SlotDeleter) DeleteBookedSlotFromClinician(ctx context.Context, bookingID int, notifyPatient bool, clinicianID int) error {
	booking, err := d.bookingGetter.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}
	if err := d.bookingDeleter.DeleteBooking(ctx, booking.ID, clinicianID); err != nil {
		return err
	}
	if notifyPatient {
		return d.cancelMailer.MailCancelBookingToPatient(&booking)
	}
	return nil
}
