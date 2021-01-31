package deiz

import (
	"context"
	"fmt"
	"github.com/audrenbdb/deiz/gcalendar"
	"github.com/audrenbdb/deiz/gmaps"
	"time"
)

type Booking struct {
	ID        int           `json:"id" validate:"required"`
	DeleteID  string        `json:"deleteId"`
	Start     time.Time     `json:"start" validate:"required"`
	End       time.Time     `json:"end" validate:"required"`
	Motive    BookingMotive `json:"motive" validate:"required"`
	Clinician Clinician     `json:"clinician" validate:"required"`
	Patient   Patient       `json:"patient" validate:"required"`
	Address   Address       `json:"address" validate:"required"`
	Remote    bool          `json:"remote"`
	Paid      bool          `json:"paid"`
	Blocked   bool          `json:"blocked"`
	Note      string        `json:"note"`
}

const (
	ErrBookingSlotAlreadyFilled Error = "booking slot already filled"
)

//driver functions
type (
	freeBookingSlotFiller interface {
		FillFreeBookingSlot(ctx context.Context, b *Booking) error
	}
	bookingSlotRemover interface {
		RemoveBookingSlot(ctx context.Context, b *Booking) error
	}
	bookingsInTimeRangeGetter interface {
		GetBookingsInTimeRange(ctx context.Context, from, to time.Time, clinicianID int) ([]Booking, error)
	}
	bookingMailer interface {
		MailBookingToPatient(ctx context.Context, b *Booking, tz *time.Location, gCalLink, gMapsLink, cancelURL string) error
		MailBookingToClinician(ctx context.Context, b *Booking, tz *time.Location, gCalLink string) error
	}
	bookingCancelMailer interface {
		MailCancelBookingToPatient(ctx context.Context, b *Booking, tz *time.Location) error
		MailCancelBookingToClinician(ctx context.Context, b *Booking, tz *time.Location) error
	}
)

//core functions
type (
	//FillFreeBookingSlot fills a free booking slot with a blocked one or a clinician appointment
	FillFreeBookingSlot func(ctx context.Context, s *Booking, clinicianID int) error
	//FreeBookingSlot marks a blocked / booked booking slot as free
	FreeBookingSlot func(ctx context.Context, s *Booking, clinicianID int) error
	//GetAllBookingSlotsFromWeek returns all booking slots of a given week for a given clinician
	GetAllBookingSlotsFromWeek func(ctx context.Context, from time.Time, clinicianID int) ([]Booking, error)
	//GetFreeBookingSlotsFromWeek returns free slots for a given week and a given clinician
	GetFreeBookingSlotsFromWeek func(ctx context.Context, from time.Time, clinicianID int) ([]Booking, error)
	//MailBooking send an email with booking details
	MailBooking func(ctx context.Context, b *Booking, sendToPatient, sendToClinician bool) error
	//MailCancelBooking send an email confirming booking cancel
	MailCancelBooking func(ctx context.Context, b *Booking, sendToPatient, sendToClinician bool) error
)

func fillFreeBookingSlotFunc(filler freeBookingSlotFiller) FillFreeBookingSlot {
	return func(ctx context.Context, b *Booking, clinicianID int) error {
		b.Clinician.ID = clinicianID
		return filler.FillFreeBookingSlot(ctx, b)
	}
}

func mailCancelBookingFunc(mailer bookingCancelMailer, tz clinicianTimezoneGetter) MailCancelBooking {
	return func(ctx context.Context, b *Booking, sendToPatient, sendToClinician bool) error {
		loc, err := getClinicianTimezoneLoc(ctx, b.Clinician.ID, tz)
		if err != nil {
			return err
		}
		if sendToPatient {
			err := mailer.MailCancelBookingToPatient(ctx, b, loc)
			if err != nil {
				return err
			}
		}
		if sendToClinician {
			err := mailer.MailCancelBookingToClinician(ctx, b, loc)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func mailBookingFunc(mailer bookingMailer, tz clinicianTimezoneGetter) MailBooking {
	return func(ctx context.Context, b *Booking, sendToPatient, sendToClinician bool) error {
		loc, err := getClinicianTimezoneLoc(ctx, b.Clinician.ID, tz)
		if err != nil {
			return err
		}
		b.Start = b.Start.In(loc)
		b.End = b.End.In(loc)
		var gCalLink string
		var gMapsLink string
		gCalEvent := gcalendar.Event{
			Start: fmt.Sprintf("%d%02d%02dT%02d%02d00", b.Start.Year(), b.Start.Month(), b.Start.Day(), b.Start.Hour(), b.Start.Minute()),
			End:   fmt.Sprintf("%d%02d%02dT%02d%02d00", b.End.Year(), b.End.Month(), b.End.Day(), b.End.Hour(), b.End.Minute()),
		}
		if !b.Remote {
			addressStr := fmt.Sprintf("%s, %d %s", b.Address.Line, b.Address.PostCode, b.Address.City)
			gCalEvent.Location = addressStr
			gCalLink = gcalendar.NewEventURL(gCalEvent)
			gMapsLink = gmaps.NewQueryAddressURL(addressStr)
		}

		if sendToPatient {
			err := mailer.MailBookingToPatient(ctx, b, loc, gCalLink, gMapsLink, getCancelBookingURL(b.DeleteID))
			if err != nil {
				return err
			}
		}
		if sendToClinician {
			err := mailer.MailBookingToClinician(ctx, b, loc, gCalLink)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func freeBookingSlotFunc(remover bookingSlotRemover) FreeBookingSlot {
	return func(ctx context.Context, b *Booking, clinicianID int) error {
		b.Clinician.ID = clinicianID
		return remover.RemoveBookingSlot(ctx, b)
	}
}

func getAllBookingSlotsFromWeekFunc(getter bookingsInTimeRangeGetter, settings calendarSettingsGetter, officeHours officeHoursGetter) GetAllBookingSlotsFromWeek {
	return func(ctx context.Context, start time.Time, clinicianID int) ([]Booking, error) {
		const daysToFetch = 6
		var bookings []Booking
		end := start.AddDate(0, 0, daysToFetch)
		s, err := settings.GetCalendarSettings(ctx, clinicianID)
		if err != nil {
			return nil, err
		}
		h, err := officeHours.GetOfficeHours(ctx, clinicianID)
		if err != nil {
			return nil, err
		}
		loc, err := time.LoadLocation(s.Timezone.Name)
		if err != nil {
			return nil, err
		}
		freeBookings := fillOfficeHoursWithFreeSlots(start.In(loc), end.In(loc), Clinician{ID: clinicianID}, h, []Booking{}, s.DefaultMotive, loc)
		bookedSlots, err := getter.GetBookingsInTimeRange(ctx, start, end, clinicianID)
		for _, b := range bookedSlots {
			bookings = append(bookings, removeOverlappingFreeSlots(b, []Booking{b}, freeBookings)...)
		}
		return bookings, nil
	}
}

func getFreeBookingSlotsFromWeekFunc(getter bookingsInTimeRangeGetter, settings calendarSettingsGetter, officeHours officeHoursGetter) GetFreeBookingSlotsFromWeek {
	return func(ctx context.Context, start time.Time, clinicianID int) ([]Booking, error) {
		const daysToFetch = 6
		var bookings []Booking
		end := start.AddDate(0, 0, daysToFetch)
		s, err := settings.GetCalendarSettings(ctx, clinicianID)
		if err != nil {
			return nil, err
		}
		h, err := officeHours.GetOfficeHours(ctx, clinicianID)
		if err != nil {
			return nil, err
		}
		loc, err := time.LoadLocation(s.Timezone.Name)
		if err != nil {
			return nil, err
		}
		freeBookings := fillOfficeHoursWithFreeSlots(start.In(loc), end.In(loc), Clinician{ID: clinicianID}, h, []Booking{}, s.DefaultMotive, loc)
		bookedSlots, err := getter.GetBookingsInTimeRange(ctx, start, end, clinicianID)
		for _, b := range bookedSlots {
			bookings = append(bookings, removeOverlappingFreeSlots(b, []Booking{}, freeBookings)...)
		}
		return bookings, nil
	}
}

func removeOverlappingFreeSlots(b Booking, bookings []Booking, freeSlots []Booking) []Booking {
	if freeSlots == nil || len(freeSlots) == 0 {
		return bookings
	}
	if !timeRangesOverlaps(b.Start, b.End, freeSlots[0].Start, freeSlots[0].End) {
		bookings = append(bookings, freeSlots[0])
	}
	return removeOverlappingFreeSlots(b, bookings, freeSlots[1:])
}

//fillTimeRangeWithFreeSlots fills a time range with free booking slots available with UTC timezone
func fillTimeRangeWithFreeSlots(end, anchor time.Time, c Clinician, b []Booking, a Address, m BookingMotive) []Booking {
	nextAnchor := anchor.Add(time.Minute * time.Duration(m.Duration))
	if nextAnchor.After(end) {
		return b
	}
	bookingStart := anchor
	bookingEnd := nextAnchor
	b = append(b, Booking{
		Start:     bookingStart.UTC(),
		End:       bookingEnd.UTC(),
		Motive:    m,
		Address:   a,
		Clinician: c,
	})
	return fillTimeRangeWithFreeSlots(end, nextAnchor, c, b, a, m)
}

func fillOfficeHoursWithFreeSlots(start, end time.Time, c Clinician, hours []OfficeHours, b []Booking, m BookingMotive, loc *time.Location) []Booking {
	if hours == nil || len(hours) == 0 {
		return b
	}
	h := hours[0]
	opening, closing := getOfficeHoursTimeRange(start, end, h, loc)
	b = fillTimeRangeWithFreeSlots(closing, opening, c, b, h.Address, m)
	return fillOfficeHoursWithFreeSlots(start, end, c, hours[1:], b, m, loc)
}

//getOfficeHoursTimeRange converts generic office hours into time range within a given time range
func getOfficeHoursTimeRange(anchor, end time.Time, h OfficeHours, loc *time.Location) (time.Time, time.Time) {
	if int(anchor.Weekday()) == h.WeekDay {
		officeOpensAt := time.Date(anchor.Year(), anchor.Month(), anchor.Day(), h.StartMn/60, h.StartMn%60, 0, 0, loc)
		officeClosesAt := time.Date(anchor.Year(), anchor.Month(), anchor.Day(), h.EndMn/60, h.EndMn%60, 0, 0, loc)
		return limitTimeRange(anchor, end, officeOpensAt, officeClosesAt)
	}
	//abort if above given time range
	if anchor.After(end) {
		return time.Time{}, time.Time{}
	}
	nextDay := anchor.Add(time.Hour * time.Duration(24))
	nextAnchor := time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, loc)
	return getOfficeHoursTimeRange(nextAnchor, end, h, loc)
}

//merge two time ranges together such as second one cannot overlaps first one
func limitTimeRange(lowerLimit, upperLimit, rangeStart, rangeEnd time.Time) (time.Time, time.Time) {
	if !timeRangesOverlaps(lowerLimit, upperLimit, rangeStart, rangeEnd) {
		return time.Time{}, time.Time{}
	}
	if rangeStart.Before(lowerLimit) {
		rangeStart = lowerLimit
	}
	if rangeEnd.After(upperLimit) {
		rangeEnd = upperLimit
	}
	return rangeStart, rangeEnd
}

//timeRangesOverlaps checks if two time ranges overlaps
func timeRangesOverlaps(startA, endA, startB, endB time.Time) bool {
	return startA.Before(endB) && startB.Before(endA)
}
