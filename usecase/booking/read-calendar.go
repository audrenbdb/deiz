package booking

import (
	"context"
	"fmt"
	"github.com/audrenbdb/deiz"
	"time"
)

type (
	officeHoursGetter interface {
		GetClinicianOfficeHours(ctx context.Context, clinicianID int) ([]deiz.OfficeHours, error)
	}
)

type readCalendar struct {
	loc               *time.Location
	officeHoursGetter officeHoursGetter

	bookingsGetter clinicianBookingsInTimeRangeGetter
}

func NewCalendarReaderUsecase(loc *time.Location, officeHoursGetter officeHoursGetter, bookingsGetter clinicianBookingsInTimeRangeGetter) *readCalendar {
	return &readCalendar{
		loc:               loc,
		officeHoursGetter: officeHoursGetter,
		bookingsGetter:    bookingsGetter,
	}
}

func (r *readCalendar) GetClinicianBookingSlots(ctx context.Context, start time.Time, motiveID, motiveDuration, clinicianID int) ([]deiz.Booking, error) {
	existingBookings, freeBookingSlots, err := r.getBookingSlots(ctx, start, motiveID, motiveDuration, clinicianID)
	if err != nil {
		return nil, fmt.Errorf("unable to get booking slots: %s", err)
	}
	return append(existingBookings, freeBookingSlots...), nil
}

func (r *readCalendar) GetPublicBookingSlots(ctx context.Context, start time.Time, motiveID, motiveDuration, clinicianID int) ([]deiz.Booking, error) {
	_, freeBookingSlots, err := r.getBookingSlots(ctx, start, motiveID, motiveDuration, clinicianID)
	if err != nil {
		return nil, fmt.Errorf("unable to get booking slots: %s", err)
	}
	return freeBookingSlots, nil
}

func (r *readCalendar) getBookingSlots(ctx context.Context, start time.Time, motiveID, motiveDuration, clinicianID int) ([]deiz.Booking, []deiz.Booking, error) {
	end := start.AddDate(0, 0, 7)
	existingBookings, err := r.bookingsGetter.GetClinicianBookingsInTimeRange(ctx, start, end, clinicianID)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get bookings in given timerange: %s", err)
	}
	freeBookingSlots, err := r.getFreeBookingSlots(ctx, start, end, existingBookings, deiz.BookingMotive{ID: motiveID, Duration: motiveDuration}, clinicianID)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get free booking slots: %s", err)
	}
	return existingBookings, freeBookingSlots, nil
}

func (r *readCalendar) getFreeBookingSlots(ctx context.Context, start, end time.Time, existingBookings []deiz.Booking, defaultMotive deiz.BookingMotive, clinicianID int) ([]deiz.Booking, error) {
	availabilities, err := r.getOfficeHoursAvailabilities(ctx, start, end, clinicianID)
	if err != nil {
		return nil, fmt.Errorf("unable to get clinician availabilities: %s", err)
	}
	bookingSlots := []deiz.Booking{}
	for _, availability := range availabilities {
		bookingSlots = append(bookingSlots,
			splitAvailabilityInFreeBookingSlots(availability, existingBookings,
				defaultMotive, clinicianID, []deiz.Booking{})...)
	}
	return bookingSlots, nil
}

func splitAvailabilityInFreeBookingSlots(availability officeHoursAvailability, existingBookings []deiz.Booking, defaultMotive deiz.BookingMotive, clinicianID int, freeBookings []deiz.Booking) []deiz.Booking {
	nextFreeBooking := deiz.Booking{
		Start:     availability.availableTimeRange[0],
		End:       availability.availableTimeRange[0].Add(time.Minute * time.Duration(defaultMotive.Duration)),
		Address:   availability.hours.Address,
		Remote:    availability.hours.Address.IsNotSet(),
		Motive:    defaultMotive,
		Clinician: deiz.Clinician{ID: clinicianID},
	}
	for _, booking := range existingBookings {
		if timeRangesOverlaps([2]time.Time{nextFreeBooking.Start, nextFreeBooking.End}, [2]time.Time{booking.Start, booking.End}) {
			nextFreeBooking.Start = booking.End
			nextFreeBooking.End = nextFreeBooking.Start.Add(time.Minute * time.Duration(defaultMotive.Duration))
		}
	}
	if nextFreeBooking.End.After(availability.availableTimeRange[1]) {
		return freeBookings
	}
	availability.availableTimeRange[0] = nextFreeBooking.End
	return splitAvailabilityInFreeBookingSlots(
		availability,
		existingBookings,
		defaultMotive,
		clinicianID, append(freeBookings, nextFreeBooking))
}

func (r *readCalendar) getOfficeHoursAvailabilities(ctx context.Context, start, end time.Time, clinicianID int) ([]officeHoursAvailability, error) {
	officeHours, err := r.officeHoursGetter.GetClinicianOfficeHours(ctx, clinicianID)
	if err != nil {
		return nil, err
	}
	var officeHoursRanges []officeHoursAvailability
	for _, h := range officeHours {
		officeHoursRanges = append(officeHoursRanges,
			officeHoursAvailability{
				hours:              h,
				availableTimeRange: r.convertOfficeHoursToTimeRange(start, end, h),
			})
	}
	return officeHoursRanges, nil
}

func (r *readCalendar) convertOfficeHoursToTimeRange(start, end time.Time, h deiz.OfficeHours) [2]time.Time {
	y, m, d := start.In(r.loc).Date()
	if h.IsWithinDate(start.In(r.loc)) {
		officeOpensAt := time.Date(y, m, d, h.StartMn/60, h.StartMn%60, 0, 0, r.loc).UTC()
		officeClosesAt := time.Date(y, m, d, h.EndMn/60, h.EndMn%60, 0, 0, r.loc).UTC()
		return constraintTimeRangeWithinLimit(start, end, [2]time.Time{officeOpensAt, officeClosesAt})
	}
	if start.After(end) {
		return [2]time.Time{time.Time{}, time.Time{}}
	}
	return r.convertOfficeHoursToTimeRange(time.Date(y, m, d+1, 0, 0, 0, 0, time.UTC), end, h)
}

func constraintTimeRangeWithinLimit(lowerLimit, upperLimit time.Time, timeRange [2]time.Time) [2]time.Time {
	if !timeRangesOverlaps([2]time.Time{lowerLimit, upperLimit}, timeRange) {
		return [2]time.Time{time.Time{}, time.Time{}}
	}
	if timeRange[0].Before(lowerLimit) {
		timeRange[0] = lowerLimit
	}
	if timeRange[1].After(upperLimit) {
		timeRange[1] = upperLimit
	}
	return timeRange
}

func timeRangesOverlaps(timeRangeA, timeRangeB [2]time.Time) bool {
	return timeRangeA[0].Before(timeRangeB[1]) && timeRangeB[0].Before(timeRangeA[1])
}

type officeHoursAvailability struct {
	hours              deiz.OfficeHours
	availableTimeRange [2]time.Time
}
