/*
Package usecase references all usecases to be implemented
*/
package usecase

import (
	"context"
	"github.com/audrenbdb/deiz"
	"time"
)

type (
	BookingUsecases struct {
		Register       BookingRegister
		PreRegister    BookingPreRegister
		CalendarReader CalendarReader
		SlotDeleter    BookingSlotDeleter
		SlotBlocker    BookingSlotBlocker
	}
)

type (
	BookingRegister interface {
		RegisterBookingFromClinician(ctx context.Context, b *deiz.Booking, clinicianID int, notifyPatient bool) error
		RegisterBookingFromPatient(ctx context.Context, b *deiz.Booking) error
		RegisterPreRegisteredBooking(ctx context.Context, b *deiz.Booking, clinicianID int, notifyPatient bool) error
	}
	BookingSlotDeleter interface {
		DeleteBlockedSlot(ctx context.Context, bookingID, clinicianID int) error
		DeletePreRegisteredSlot(ctx context.Context, bookingID, clinicianID int) error
		DeleteBookedSlotFromPatient(ctx context.Context, deleteID string) error
		DeleteBookedSlotFromClinician(ctx context.Context, bookingID int, notifyPatient bool, clinicianID int) error
	}
	BookingPreRegister interface {
		PreRegisterBooking(ctx context.Context, b *deiz.Booking, clinicianID int) error
	}
	BookingSlotBlocker interface {
		BlockBookingSlot(ctx context.Context, slot *deiz.Booking, clinicianID int) error
		BlockBookingSlotList(ctx context.Context, slots []*deiz.Booking, credentials deiz.Credentials) error
	}
	CalendarReader interface {
		GetCalendarFreeSlots(ctx context.Context, start time.Time, motive deiz.BookingMotive, clinicianID int) ([]deiz.Booking, error)
		GetCalendarSlots(ctx context.Context, start time.Time, motive deiz.BookingMotive, clinicianID int) ([]deiz.Booking, error)
	}
)
