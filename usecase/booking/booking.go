package booking

import (
	"context"
	"errors"
	"github.com/audrenbdb/deiz"
	"time"
)

type (
	BookingsInTimeRangeGetter interface {
		GetBookingsInTimeRange(ctx context.Context, start, end time.Time, clinicianID int) ([]deiz.Booking, error)
	}
	Creater interface {
		CreateBooking(ctx context.Context, b *deiz.Booking) error
	}
	Deleter interface {
		DeleteBooking(ctx context.Context, bookingID, clinicianID int) error
	}
)

//IsBookingValid ensure minimum fields of booking are valid
func IsBookingValid(b *deiz.Booking) bool {
	if b.Start.After(b.End) || b.Clinician.ID == 0 {
		return false
	}
	return true
}

func (u *Usecase) GetBookingSlots(ctx context.Context, start time.Time, tzName string, defaultMotiveID, defaultMotiveDuration, clinicianID int) ([]deiz.Booking, error) {
	loc, err := time.LoadLocation(tzName)
	if err != nil || tzName == "" {
		return nil, deiz.ErrorParsingTimezone
	}
	officeHours, err := u.OfficeHoursGetter.GetClinicianOfficeHours(ctx, clinicianID)
	if err != nil {
		return nil, err
	}
	var bookings []deiz.Booking
	end := start.AddDate(0, 0, 6)
	freeBookings := FillOfficeHoursWithFreeBookingSlots(
		start, end,
		deiz.Clinician{ID: clinicianID}, officeHours, []deiz.Booking{}, deiz.BookingMotive{ID: defaultMotiveID, Duration: defaultMotiveDuration}, loc)
	bookedSlots, err := u.BookingsInTimeRangeGetter.GetBookingsInTimeRange(ctx, start, end, clinicianID)
	if err != nil {
		return nil, err
	}
	bookings = RemoveOverlappingFreeBookingSlots(freeBookings, bookedSlots, []deiz.Booking{})
	return append(bookings, bookedSlots...), nil
}

func RemoveOverlappingFreeBookingSlots(freeSlots, bookedSlots, slotsToKeep []deiz.Booking) []deiz.Booking {
	if freeSlots == nil || len(freeSlots) == 0 {
		return slotsToKeep
	}
	slot := freeSlots[0]
	overlaps := false
	for _, b := range bookedSlots {
		if TimeRangesOverlaps(slot.Start, slot.End, b.Start, b.End) {
			overlaps = true
			break
		}
	}
	if !overlaps {
		slotsToKeep = append(slotsToKeep, slot)
	}
	return RemoveOverlappingFreeBookingSlots(freeSlots[1:], bookedSlots, slotsToKeep)
}

func (u *Usecase) BlockBookingSlot(ctx context.Context, b *deiz.Booking, clinicianID int) error {
	if b.Clinician.ID != clinicianID {
		return deiz.ErrorUnauthorized
	}
	if b.Patient.ID != 0 || b.Address.ID != 0 || b.Motive.ID != 0 || b.Note != "" {
		return errors.New("booking is not empty")
	}
	return u.Creater.CreateBooking(ctx, b)
}

func (u *Usecase) UnlockBookingSlot(ctx context.Context, bookingID, clinicianID int) error {
	return u.Deleter.DeleteBooking(ctx, bookingID, clinicianID)
}
