package booking2

import (
	"context"
	"fmt"
	"sort"
	"time"
)

type booking struct {
	ID          int         `json:"id"`
	Description string      `json:"description"`
	Start       time.Time   `json:"start"`
	End         time.Time   `json:"end"`
	Patient     person      `json:"patient"`
	Clinician   person      `json:"clinician"`
	Confirmed   bool        `json:"confirmed"`
	Timezone    string      `json:"timezone"`
	Recurrence  recurrence  `json:"recurrence"`
	Address     string      `json:"address"`
	MeetingMode meetingMode `json:"meetingMode"`
	Type        bookingType `json:"bookingType"`
}

type person struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

type recurrence uint8

type bookingType uint8

const (
	blockedBooking bookingType = iota
	appointmentBooking
	eventBooking
)

const (
	nonRecurrent recurrence = iota
	dailyRecurrent
	weeklyRecurrent
	monthlyRecurrent
)

type officeHours struct {
	start       time.Time
	end         time.Time
	address     string
	meetingMode meetingMode
	tRangeCfg   timeRangeCfg
}

type meetingMode uint8

const (
	remoteMeeting meetingMode = iota
	inOfficeMeeting
	atExternalAddressMeeting
)

//timeRangeCfg is a generic struct to define a time range
type timeRangeCfg struct {
	weekDay int
	startMn int
	endMn   int
	loc     *time.Location
}
type timeRange struct {
	start time.Time
	end   time.Time
}

type address struct {
	id       int
	line     string
	postCode int
	city     string
}

func (a address) toString() string {
	if a.id == 0 {
		return ""
	}
	return fmt.Sprintf("%s, %d %s", a.line, a.postCode, a.city)
}

//getClinicianWeek retrieves list of week bookings and available time slots
type getClinicianWeek = func(ctx context.Context, from time.Time, bookingDuration time.Duration, clinicianID int) ([]booking, error)

func createGetClinicianWeekFunc(
	getClinicianBookingsInWeek getClinicianBookingsInWeek,
	getOfficeHoursInWeek getOfficeHoursInWeek,
) getClinicianWeek {
	return func(ctx context.Context, weekStart time.Time, bookingDuration time.Duration, clinicianID int) ([]booking, error) {
		w := getWeekTimeRangeFromStartDate(weekStart)
		existingBookings, err := getClinicianBookingsInWeek(ctx, w, clinicianID)
		if err != nil {
			return nil, err
		}
		officeHours, err := getOfficeHoursInWeek(ctx, w, clinicianID)
		freeBookings := getFreeBookingSlotsInWeek(existingBookings, officeHours, bookingDuration)
		if err != nil {
			return nil, err
		}
		return append(existingBookings, freeBookings...), nil
	}
}

func getFreeBookingSlotsInWeek(existingBookings []booking, officeHours []officeHours, slotDuration time.Duration) []booking {
	freeBookingSlots := []booking{}
	for _, h := range officeHours {
		freeBookingSlots = append(freeBookingSlots,
			splitOfficeHoursInFreeBookingSlots(h, existingBookings, slotDuration)...)
	}
	return freeBookingSlots
}

func splitOfficeHoursInFreeBookingSlots(h officeHours, existingBookings []booking, slotDuration time.Duration) []booking {
	freeBookingSlots := []booking{}
	availability := timeRange{h.start, h.end}
	for isEnoughTimeAvailableForSlot(slotDuration, availability) {
		nextSlot := getNextSlotAvailableInTimeRange(availability, slotDuration, existingBookings)
		freeBookingSlots = append(freeBookingSlots, booking{
			Start:       nextSlot.start,
			End:         nextSlot.end,
			Address:     h.address,
			MeetingMode: h.meetingMode,
		})
		availability.start = nextSlot.end
	}
	return freeBookingSlots
}

func getNextSlotAvailableInTimeRange(r timeRange, duration time.Duration, existingBookings []booking) timeRange {
	slot := timeRange{start: r.start, end: r.start.Add(duration)}
	for _, b := range existingBookings {
		if timeRangesOverlaps(slot, timeRange{b.Start, b.End}) {
			slot.start = b.End
			slot.end = slot.start.Add(duration)
		}
	}
	return slot
}

func sortBookingsByDate(bookings []booking) {
	sort.SliceStable(bookings, func(i, j int) bool {
		return bookings[i].Start.Before(bookings[j].Start)
	})
}

func isEnoughTimeAvailableForSlot(duration time.Duration, r timeRange) bool {
	return durationInMnBetweenTwoDates(r.start, r.end) > int(duration.Minutes())
}

func durationInMnBetweenTwoDates(d1, d2 time.Time) int {
	return int(d2.Sub(d1).Minutes())
}

//getClinicianBookingsInWeek retrieves all bookings withing a given week
type getClinicianBookingsInWeek = func(ctx context.Context, w timeRange, clinicianID int) ([]booking, error)

func createGetClinicianBookingsInWeek(
	getClinicianNonRecurrentBookingsInTimeRange getClinicianNonRecurrentBookingsInTimeRange,
	getClinicianRecurrentBookingsInTimeRange getClinicianRecurrentBookingsInTimeRange,
) getClinicianBookingsInWeek {
	return func(ctx context.Context, w timeRange, clinicianID int) ([]booking, error) {
		nonRecurrentBookings, err := getClinicianNonRecurrentBookingsInTimeRange(ctx, w, clinicianID)
		if err != nil {
			return nil, err
		}
		recurrentBookings, err := getClinicianRecurrentBookingsInTimeRange(ctx, w, clinicianID)
		if err != nil {
			return nil, err
		}
		bookings := append(nonRecurrentBookings, recurrentBookings...)
		sortBookingsByDate(bookings)
		return bookings, nil
	}
}

type getClinicianRecurrentBookingsInTimeRange = func(ctx context.Context, tr timeRange, clinicianID int) ([]booking, error)

func createGetClinicianRecurrentBookingsInTimeRange(
	getClinicianRecurrentBookings getClinicianRecurrentBookings,
) getClinicianRecurrentBookingsInTimeRange {
	return func(ctx context.Context, tr timeRange, clinicianID int) ([]booking, error) {
		recurrentBookings, err := getClinicianRecurrentBookings(ctx, clinicianID)
		if err != nil {
			return nil, err
		}
		for _, b := range recurrentBookings {
			setRecurrentBookingTimeRange(&b, tr.start)
		}
		return recurrentBookings, nil
	}
}

func setRecurrentBookingTimeRange(b *booking, anchor time.Time) {
	trCfg := getTimeRangeCfgFromBooking(b)
	tr := getFirstTimeRangeMatchingTimeRangeCfg(anchor, trCfg)
	b.Start = tr.start
	b.End = tr.end
}

func getTimeRangeCfgFromBooking(b *booking) timeRangeCfg {
	var trCfg timeRangeCfg
	loc := parseTimezone(b.Timezone)
	startInLocal := b.Start.In(loc)
	endInLocal := b.End.In(loc)
	trCfg.weekDay = int(startInLocal.Weekday())
	trCfg.startMn = getTotalMinutesFromDay(startInLocal)
	trCfg.endMn = getTotalMinutesFromDay(endInLocal)
	trCfg.loc = loc
	return trCfg
}

//attempts to parse a timezone and returns utc if it fails
func parseTimezone(tz string) *time.Location {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.UTC
	}
	return loc
}

//getTotalMinutesFromDay retrieves total minutes elapses in a given day.
//example : 08:30 is 510mn
func getTotalMinutesFromDay(d time.Time) int {
	return d.Hour()*60 + d.Minute()
}

//getNonRecurrentBookingsInTimeRange retrieves a list of non recurrent bookings in time range
type getClinicianNonRecurrentBookingsInTimeRange = func(ctx context.Context, tr timeRange, clinicianID int) ([]booking, error)

func createGetClinicianNonRecurrentBookingsInTimeRangeFunc(
	getClinicianBookingsInTimeRange getClinicianBookingsInTimeRange) getClinicianNonRecurrentBookingsInTimeRange {
	return func(ctx context.Context, tr timeRange, clinicianID int) ([]booking, error) {
		bookings, err := getClinicianBookingsInTimeRange(ctx, tr, clinicianID)
		if err != nil {
			return nil, err
		}
		return filterNonRecurrentBookings(bookings), nil
	}
}

func filterNonRecurrentBookings(bookings []booking) []booking {
	var nonRecurrentBookings []booking
	for _, b := range bookings {
		if b.Recurrence == nonRecurrent {
			nonRecurrentBookings = append(nonRecurrentBookings, b)
		}
	}
	return nonRecurrentBookings
}

//getOfficeHoursInWeek retrieves clinician office hours within a given time range.
type getOfficeHoursInWeek = func(ctx context.Context, tr timeRange, clinicianID int) ([]officeHours, error)

func createGetOfficeHoursInWeekFunc(
	getClinicianOfficeHours getClinicianOfficeHours,
) getOfficeHoursInWeek {
	return func(ctx context.Context, tr timeRange, clinicianID int) ([]officeHours, error) {
		officeHoursList, err := getClinicianOfficeHours(ctx, clinicianID)
		if err != nil {
			return nil, err
		}
		setOfficeHoursTimeRange(tr, officeHoursList)
		return officeHoursList, nil
	}
}

func setOfficeHoursTimeRange(tr timeRange, hours []officeHours) {
	for _, h := range hours {
		availableTimeRange :=
			constraintTimeRangeWithinLimit(tr,
				getFirstTimeRangeMatchingTimeRangeCfg(tr.start, h.tRangeCfg))
		h.start = availableTimeRange.start
		h.end = availableTimeRange.end
	}
}

func constraintTimeRangeWithinLimit(limit timeRange, tr timeRange) timeRange {
	tr.start = constraintTimeWithinLimit(limit, tr.start)
	tr.end = constraintTimeWithinLimit(limit, tr.end)
	return tr
}

func constraintTimeWithinLimit(limit timeRange, t time.Time) time.Time {
	if t.Before(limit.start) {
		return limit.start
	}
	if t.After(limit.end) {
		return limit.end
	}
	return t
}

func timeRangesOverlaps(trA, trB timeRange) bool {
	return trA.start.Before(trB.end) && trB.start.Before(trA.end)
}

//getFirstTimeRangeMatchingTimeRangeCfg
func getFirstTimeRangeMatchingTimeRangeCfg(anchor time.Time, cfg timeRangeCfg) timeRange {
	if timeRangeCfgMatchDateWeekday(anchor, cfg) {
		return createTimeRangeFromTimeRangeCfg(anchor, cfg)
	}
	return getFirstTimeRangeMatchingTimeRangeCfg(
		getNextLocalDayToUTC(anchor, cfg.loc), cfg)
}

func createTimeRangeFromTimeRangeCfg(date time.Time, cfg timeRangeCfg) timeRange {
	y, m, d := date.In(cfg.loc).Date()
	return timeRange{
		start: newDateWithGivenMn(y, m, d, cfg.startMn, cfg.loc),
		end:   newDateWithGivenMn(y, m, d, cfg.endMn, cfg.loc),
	}
}

//getNextLocalDay returns next local day time at midnight in UTC.
//Example : initial date is 02/01/2006 15:20 in local tz
//returning date will be 03/01/2006 00:00 in local tz
func getNextLocalDayToUTC(date time.Time, tz *time.Location) time.Time {
	y, m, d := date.In(tz).Date()
	return time.Date(y, m, d+1, 0, 0, 0, 0, tz).UTC()
}

func newDateWithGivenMn(year int, month time.Month, day int, mn int, tz *time.Location) time.Time {
	return time.Date(year, month, day, mn/60, mn%60, 0, 0, tz).UTC()
}

func timeRangeCfgMatchDateWeekday(date time.Time, timeRangeCfg timeRangeCfg) bool {
	return int(date.In(timeRangeCfg.loc).Weekday()) == timeRangeCfg.weekDay
}

func getWeekTimeRangeFromStartDate(start time.Time) timeRange {
	return timeRange{
		start: start,
		end:   start.AddDate(0, 0, 7),
	}
}
