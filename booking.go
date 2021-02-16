package deiz

import (
	"time"
)

type Booking struct {
	ID        int           `json:"id"`
	DeleteID  string        `json:"deleteId"`
	Start     time.Time     `json:"start"`
	End       time.Time     `json:"end"`
	Motive    BookingMotive `json:"motive"`
	Clinician Clinician     `json:"clinician"`
	Patient   Patient       `json:"patient"`
	Address   Address       `json:"address"`
	Remote    bool          `json:"remote"`
	Paid      bool          `json:"paid"`
	Blocked   bool          `json:"blocked"`
	Confirmed bool          `json:"confirmed"`
	Note      string        `json:"note"`
}

/*
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
		/*
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
*/

/*
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
		if err != nil {
			return nil, err
		}
		bookings = removeOverlappingFreeSlots(freeBookings, bookedSlots, []Booking{})
		return bookings, nil
	}
}

*/
