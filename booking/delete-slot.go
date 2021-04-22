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

type DeleteSlotUsecase struct {
	BookingGetter  bookingGetter
	BookingDeleter bookingDeleter
	CancelMailer   cancelMailer
}

func (d *DeleteSlotUsecase) DeleteBlockedSlot(ctx context.Context, bookingID, clinicianID int) error {
	return d.BookingDeleter.DeleteBooking(ctx, bookingID, clinicianID)
}

func (d *DeleteSlotUsecase) DeletePreRegisteredSlot(ctx context.Context, bookingID, clinicianID int) error {
	return d.BookingDeleter.DeleteBooking(ctx, bookingID, clinicianID)
}

func (d *DeleteSlotUsecase) DeleteBookedSlotFromPatient(ctx context.Context, deleteID string) error {
	booking, err := d.BookingGetter.GetBookingByDeleteID(ctx, deleteID)
	if err != nil {
		return err
	}
	if err := d.BookingDeleter.DeleteBooking(ctx, booking.ID, booking.Clinician.ID); err != nil {
		return err
	}
	return d.CancelMailer.MailCancelBookingToClinician(&booking)
}

func (d *DeleteSlotUsecase) DeleteBookedSlotFromClinician(ctx context.Context, bookingID int, notifyPatient bool, clinicianID int) error {
	booking, err := d.BookingGetter.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}
	if err := d.BookingDeleter.DeleteBooking(ctx, booking.ID, clinicianID); err != nil {
		return err
	}
	if notifyPatient {
		return d.CancelMailer.MailCancelBookingToPatient(&booking)
	}
	return nil
}
