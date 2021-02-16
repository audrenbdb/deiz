package booking

import (
	"context"
	"errors"
	"fmt"
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
	Updater interface {
		UpdateBooking(ctx context.Context, b *deiz.Booking) error
	}
	GetterByID interface {
		GetBookingByID(ctx context.Context, bookingID int) (deiz.Booking, error)
	}
	OverlappingBlockedDeleter interface {
		DeleteOverlappingBlockedBooking(ctx context.Context, start, end time.Time, clinicianID int) error
	}
	ToClinicianMailer interface {
		MailBookingToClinician(ctx context.Context, b *deiz.Booking, tz *time.Location, gCalLink string) error
	}
	ToPatientMailer interface {
		MailBookingToPatient(ctx context.Context, b *deiz.Booking, tz *time.Location, gCalLink, gMapsLink, cancelURL string) error
	}
	CancelBookingToPatientMailer interface {
		MailCancelBookingToPatient(ctx context.Context, b *deiz.Booking, tz *time.Location) error
	}
)

//IsBookingValid ensure minimum fields of booking are valid
func IsBookingValid(b *deiz.Booking) bool {
	if b.Start.After(b.End) || b.Clinician.ID == 0 {
		return false
	}
	return true
}

func NewAvailableBookingSlot(start, end time.Time, address deiz.Address, motive deiz.BookingMotive) deiz.Booking {
	return deiz.Booking{
		Start:   start,
		End:     end,
		Address: address,
		Remote:  address.ID == 0,
		Motive:  motive,
	}
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
	end := start.AddDate(0, 0, 7)
	bookings, err := u.BookingsInTimeRangeGetter.GetBookingsInTimeRange(ctx, start, end, clinicianID)
	if err != nil {
		return nil, err
	}
	bookedSlotsTimeRanges := GetTimeRangesFromBookings(SortBookingsByStart(bookings), [][2]time.Time{})
	officeHoursTimeRange := GetAllOfficeHoursTimeRange(start, end, officeHours, loc)
	for i, timeRange := range officeHoursTimeRange {
		freeTimeRanges := GetTimeRangesNotOverLapping(defaultMotiveDuration, timeRange[0], timeRange[1], bookedSlotsTimeRanges, [][2]time.Time{})
		for _, timeRange := range freeTimeRanges {
			bookings = append(bookings,
				NewAvailableBookingSlot(timeRange[0], timeRange[1],
					officeHours[i].Address,
					deiz.BookingMotive{ID: defaultMotiveID, Duration: defaultMotiveDuration}))
		}
	}
	return bookings, nil

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
	b, err := u.GetterByID.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}
	if b.Clinician.ID != clinicianID {
		return deiz.ErrorUnauthorized
	}
	return u.Deleter.DeleteBooking(ctx, bookingID, clinicianID)
}

func (u *Usecase) ConfirmPreRegisteredBooking(ctx context.Context, b *deiz.Booking, clinicianID int, notifyPatient, notifyClinician bool) error {
	if b.Clinician.ID != clinicianID {
		return deiz.ErrorUnauthorized
	}
	if !b.Confirmed || b.End.Before(b.Start) || b.Blocked || b.Patient.ID == 0 || b.Clinician.ID == 0 {
		return deiz.ErrorStructValidation
	}
	err := u.Updater.UpdateBooking(ctx, b)
	if err != nil {
		return err
	}
	if notifyClinician || notifyPatient {
		return u.NotifyBooking(ctx, b, clinicianID, notifyPatient, notifyClinician)
	}
	return nil
}

func (u *Usecase) PreRegisterBooking(ctx context.Context, b *deiz.Booking, clinicianID int) error {
	if b.Clinician.ID != clinicianID {
		return deiz.ErrorUnauthorized
	}
	if b.Confirmed || b.End.Before(b.Start) || b.Blocked {
		return deiz.ErrorStructValidation
	}
	if err := u.OverlappingBlockedDeleter.DeleteOverlappingBlockedBooking(ctx, b.Start, b.End, clinicianID); err != nil {
		return err
	}
	if err := u.Creater.CreateBooking(ctx, b); err != nil {
		return err
	}
	return nil
}

func (u *Usecase) RegisterBooking(ctx context.Context, b *deiz.Booking, clinicianID int, notifyPatient, notifyClinician bool) error {
	if b.Patient.ID == 0 || b.Clinician.ID == 0 || b.End.Before(b.Start) ||
		b.Blocked || (b.Address.ID == 0 && !b.Remote) || !b.Confirmed {
		return deiz.ErrorStructValidation
	}
	if b.Clinician.ID != clinicianID {
		return deiz.ErrorUnauthorized
	}
	if err := u.OverlappingBlockedDeleter.DeleteOverlappingBlockedBooking(ctx, b.Start, b.End, clinicianID); err != nil {
		return err
	}
	if err := u.Creater.CreateBooking(ctx, b); err != nil {
		return err
	}
	if notifyPatient || notifyClinician {
		return u.NotifyBooking(ctx, b, clinicianID, notifyPatient, notifyClinician)
	}
	return nil
}

func (u *Usecase) NotifyBooking(ctx context.Context, b *deiz.Booking, clinicianID int, notifyPatient, notifyClinician bool) error {
	settings, err := u.CalendarSettingsGetter.GetClinicianCalendarSettings(ctx, clinicianID)
	if err != nil {
		return err
	}
	clinicianTz, err := time.LoadLocation(settings.Timezone.Name)
	if err != nil {
		return err
	}
	if notifyClinician {
		if err := u.MailToClinician(ctx, b, clinicianTz, clinicianID); err != nil {
			return err
		}
	}
	if notifyPatient {
		if err := u.MailToPatient(ctx, b, clinicianTz, clinicianID); err != nil {
			return err
		}
	}
	return nil
}

func (u *Usecase) RemoveBooking(ctx context.Context, bookingID int, clinicianID int, notifyPatient bool) error {
	b, err := u.GetterByID.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}
	if err := u.Deleter.DeleteBooking(ctx, bookingID, clinicianID); err != nil {
		return err
	}
	if notifyPatient {
		s, err := u.CalendarSettingsGetter.GetClinicianCalendarSettings(ctx, clinicianID)
		if err != nil {
			return err
		}
		tz, err := time.LoadLocation(s.Timezone.Name)
		if err != nil {
			return err
		}
		if err := u.CancelToPatientMailer.MailCancelBookingToPatient(ctx, &b, tz); err != nil {
			return err
		}
	}
	return nil
}

func (u *Usecase) MailToClinician(ctx context.Context, b *deiz.Booking, tz *time.Location, clinicianID int) error {
	addressStr := ""
	if !b.Remote {
		addressStr = fmt.Sprintf("%s, %d %s", b.Address.Line, b.Address.PostCode, b.Address.City)
	}
	return u.ToClinicianMailer.MailBookingToClinician(ctx,
		b, tz, u.GCalendarLinkBuilder.BuildGCalendarLink(
			b.Start.In(tz), b.End.In(tz),
			fmt.Sprintf("Consultation avec %s %s", b.Patient.Surname, b.Patient.Name),
			addressStr, b.Note,
		))
}

func (u *Usecase) MailToPatient(ctx context.Context, b *deiz.Booking, tz *time.Location, clinicianID int) error {
	addressStr := ""
	if !b.Remote {
		addressStr = fmt.Sprintf("%s, %d %s", b.Address.Line, b.Address.PostCode, b.Address.City)
	}
	return u.ToPatientMailer.MailBookingToPatient(ctx, b, tz,
		u.GCalendarLinkBuilder.BuildGCalendarLink(
			b.Start.In(tz),
			b.End.In(tz),
			fmt.Sprintf("Consultation avec %s %s", b.Clinician.Surname, b.Clinician.Name),
			addressStr,
			""),
		u.GMapsLinkBuilder.BuildGMapsLink(addressStr),
		b.DeleteID,
	)
}
