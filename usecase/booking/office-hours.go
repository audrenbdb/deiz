package booking

import (
	"context"
	"github.com/audrenbdb/deiz"
	"time"
)

type (
	ClinicianOfficeHoursGetter interface {
		GetClinicianOfficeHours(ctx context.Context, clinicianID int) ([]deiz.OfficeHours, error)
	}
)

//GetOfficeHoursTimeRange converts generic office hours into time range within a given time range
func GetOfficeHoursTimeRange(anchor, end time.Time, h deiz.OfficeHours, loc *time.Location) (time.Time, time.Time) {
	anchorInLoc := anchor.In(loc)
	y, m, d := anchorInLoc.Date()
	if int(anchorInLoc.Weekday()) == h.WeekDay {
		officeOpensAt := time.Date(y, m, d, h.StartMn/60, h.StartMn%60, 0, 0, loc).UTC()
		officeClosesAt := time.Date(y, m, d, h.EndMn/60, h.EndMn%60, 0, 0, loc).UTC()
		return LimitTimeRange(anchor, end, officeOpensAt, officeClosesAt)
	}
	//abort if above given time range
	if anchor.After(end) {
		return time.Time{}, time.Time{}
	}
	nextAnchor := time.Date(y, m, d+1, 0, 0, 0, 0, loc).UTC()
	return GetOfficeHoursTimeRange(nextAnchor, end, h, loc)
}

//FillOfficeHoursWithFreeBookingSlots converts office hours into time ranges withing the given time range
//It then fills these office time ranges with booking slots available
func FillOfficeHoursWithFreeBookingSlots(start, end time.Time, c deiz.Clinician, hours []deiz.OfficeHours,
	bookings []deiz.Booking, m deiz.BookingMotive, loc *time.Location) []deiz.Booking {
	if hours == nil || len(hours) == 0 {
		return bookings
	}
	h := hours[0]
	opening, closing := GetOfficeHoursTimeRange(start, end, h, loc)
	return FillOfficeHoursWithFreeBookingSlots(
		start, end,
		c,
		hours[1:],
		FillTimeRangeWithFreeBookingSlots(opening, closing, c, bookings, h.Address, m),
		m,
		loc)
}

//FillTimeRangeWithFreeBookingSlots fills a time range with free booking slots available with UTC timezone
func FillTimeRangeWithFreeBookingSlots(anchor, end time.Time, c deiz.Clinician, bookings []deiz.Booking,
	a deiz.Address, m deiz.BookingMotive) []deiz.Booking {
	nextAnchor := anchor.Add(time.Minute * time.Duration(m.Duration))

	if nextAnchor.After(end) {
		return bookings
	}
	bookingStart := anchor
	bookingEnd := nextAnchor
	bookings = append(bookings, deiz.Booking{
		Start:     bookingStart,
		End:       bookingEnd,
		Motive:    m,
		Address:   a,
		Clinician: c,
	})
	return FillTimeRangeWithFreeBookingSlots(nextAnchor, end, c, bookings, a, m)
}
