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

type ReadCalendarUsecase struct {
	Loc               *time.Location
	OfficeHoursGetter officeHoursGetter

	BookingsGetter bookingGetter
}

func (r *ReadCalendarUsecase) GetCalendarSlots(ctx context.Context, start time.Time, defaultDuration int, clinicianID int) ([]deiz.Booking, error) {
	existingBookings, freeBookingSlots, err := r.getBookingSlots(ctx, start, defaultDuration, clinicianID)
	if err != nil {
		return nil, fmt.Errorf("unable to get booking slots: %s", err)
	}
	return append(existingBookings, freeBookingSlots...), nil
}

func (r *ReadCalendarUsecase) GetCalendarFreeSlots(ctx context.Context, start time.Time, defaultDuration int, clinicianID int) ([]deiz.Booking, error) {
	_, freeBookingSlots, err := r.getBookingSlots(ctx, start, defaultDuration, clinicianID)
	if err != nil {
		return nil, fmt.Errorf("unable to get booking slots: %s", err)
	}
	return freeBookingSlots, nil
}

func (r *ReadCalendarUsecase) getBookingSlots(ctx context.Context, start time.Time, defaultDuration int, clinicianID int) ([]deiz.Booking, []deiz.Booking, error) {

	end := start.AddDate(0, 0, 7)
	existingBookings, err := r.BookingsGetter.GetNonRecurrentClinicianBookingsInTimeRange(ctx, start, end, clinicianID)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get bookings in given timerange: %s", err)
	}
	recurrentBookings, err := r.getRecurrentBookingsInTimeRange(ctx, clinicianID)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get existing recurrent bookings: %s", err)
	}
	existingBookings = append(existingBookings, recurrentBookings...)
	freeBookingSlots, err := r.getFreeBookingSlots(ctx, timeRange{start, end}, deiz.SortBookingByDate(existingBookings), defaultDuration, clinicianID)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get free booking slots: %s", err)
	}
	return existingBookings, freeBookingSlots, nil
}

func (r *ReadCalendarUsecase) getRecurrentBookingsInTimeRange(ctx context.Context, clinicianID int) ([]deiz.Booking, error) {
	recurrentBookings, err := r.BookingsGetter.GetClinicianWeeklyRecurrentBookings(ctx, clinicianID)
	if err != nil {
		return nil, err
	}
	for _, b := range recurrentBookings {
		r.setRecurrentBookingInTimeRange(&b)
	}
	return recurrentBookings, nil
}

func (r *ReadCalendarUsecase) setRecurrentBookingInTimeRange(b *deiz.Booking) {
	bookingTr := timeRange{
		start: b.Start.In(r.Loc),
		end:   b.End.In(r.Loc),
	}
	tr := convertCalEventToTimeRange(bookingTr, calEvent{
		weekday: int(bookingTr.start.Weekday()),
		startMn: bookingTr.start.Hour()*60 + bookingTr.start.Minute(),
		endMn:   bookingTr.end.Hour()*60 + bookingTr.end.Minute(),
		loc:     r.Loc,
	}, false)
	b.Start = tr.start
	b.End = tr.end
}

func (r *ReadCalendarUsecase) getFreeBookingSlots(ctx context.Context, timeRange timeRange, existingBookings []deiz.Booking, defaultDuration int, clinicianID int) ([]deiz.Booking, error) {
	availabilities, err := r.getOfficeHoursAvailabilities(ctx, timeRange, clinicianID)
	if err != nil {
		return nil, fmt.Errorf("unable to get clinician availabilities: %s", err)
	}
	bookingSlots := []deiz.Booking{}
	for _, availability := range availabilities {
		bookingSlots = append(bookingSlots,
			splitAvailabilityInFreeBookingSlots(availability, existingBookings,
				defaultDuration, []deiz.Booking{})...)
	}
	return bookingSlots, nil
}

func splitAvailabilityInFreeBookingSlots(availability officeHoursAvailability, existingBookings []deiz.Booking, defaultDuration int, freeBookings []deiz.Booking) []deiz.Booking {
	nextFreeBooking := deiz.Booking{
		BookingType: deiz.AppointmentBooking,
		Start:       availability.availableTimeRange.start,
		End:         availability.availableTimeRange.start.Add(time.Minute * time.Duration(defaultDuration)),
		Address:     availability.hours.Address.ToString(),
		MeetingMode: availability.hours.MeetingMode,
	}
	//make sure next free booking time range do not overlaps with existing bookings
	for _, booking := range existingBookings {
		if bookingsOverlap(&nextFreeBooking, &booking) {
			nextFreeBooking.Start = booking.End
			nextFreeBooking.End = nextFreeBooking.Start.Add(time.Minute * time.Duration(defaultDuration))
		}
	}
	if nextFreeBooking.End.After(availability.availableTimeRange.end) {
		return freeBookings
	}
	availability.availableTimeRange.start = nextFreeBooking.End
	return splitAvailabilityInFreeBookingSlots(
		availability,
		existingBookings,
		defaultDuration, append(freeBookings, nextFreeBooking))
}

func (r *ReadCalendarUsecase) getOfficeHoursAvailabilities(ctx context.Context, timeRange timeRange, clinicianID int) ([]officeHoursAvailability, error) {
	officeHours, err := r.OfficeHoursGetter.GetClinicianOfficeHours(ctx, clinicianID)
	if err != nil {
		return nil, err
	}
	var officeHoursRanges []officeHoursAvailability
	for _, h := range officeHours {
		officeHoursRanges = append(officeHoursRanges,
			officeHoursAvailability{
				hours: h,
				availableTimeRange: convertCalEventToTimeRange(timeRange, calEvent{
					weekday: h.WeekDay,
					startMn: h.StartMn,
					endMn:   h.EndMn,
					loc:     r.Loc,
				}, true),
			})
	}
	return officeHoursRanges, nil
}

//check if given dates in an array have same weekday
func datesShareSameWeekday(dates []time.Time, weekday int) bool {
	for _, d := range dates {
		if int(d.Weekday()) != weekday {
			return false
		}
	}
	return true
}

type calEvent struct {
	weekday int
	startMn int
	endMn   int
	loc     *time.Location
}

func convertCalEventToTimeRange(limit timeRange, ev calEvent, shouldConstraintWithingLimits bool) timeRange {
	y, m, d := limit.start.In(ev.loc).Date()
	if datesShareSameWeekday([]time.Time{limit.start.In(ev.loc)}, ev.weekday) {
		officeOpensAt := time.Date(y, m, d, ev.startMn/60, ev.startMn%60, 0, 0, ev.loc).UTC()
		officeClosesAt := time.Date(y, m, d, ev.endMn/60, ev.endMn%60, 0, 0, ev.loc).UTC()
		if shouldConstraintWithingLimits {
			return constraintTimeRangeWithinLimit(limit, timeRange{start: officeOpensAt, end: officeClosesAt})
		}
		return timeRange{start: officeOpensAt, end: officeClosesAt}
	}
	if limit.start.After(limit.end) {
		return timeRange{}
	}
	return convertCalEventToTimeRange(timeRange{
		start: time.Date(y, m, d+1, 0, 0, 0, 0, time.UTC),
		end:   limit.end}, ev, shouldConstraintWithingLimits)
}

func constraintTimeRangeWithinLimit(limit timeRange, tr timeRange) timeRange {
	if !timeRangesOverlaps(timeRange{limit.start, limit.end}, tr) {
		return timeRange{}
	}
	if tr.start.Before(limit.start) {
		tr.start = limit.start
	}
	if tr.end.After(limit.end) {
		tr.end = limit.end
	}
	return tr
}

func timeRangesOverlaps(trA, trB timeRange) bool {
	return trA.start.Before(trB.end) && trB.start.Before(trA.end)
}

type officeHoursAvailability struct {
	hours              deiz.OfficeHours
	availableTimeRange timeRange
}
