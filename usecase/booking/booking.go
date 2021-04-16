package booking

import (
	"context"
	"fmt"
	"github.com/audrenbdb/deiz"
	"net/url"
	"time"
)

type (
	ClinicianBookingsInTimeRangeGetter interface {
		GetClinicianBookingsInTimeRange(ctx context.Context, start, end time.Time, clinicianID int) ([]deiz.Booking, error)
	}
	BookingsInTimeRangeGetter interface {
		GetBookingsInTimeRange(ctx context.Context, start, end time.Time) ([]deiz.Booking, error)
	}
	Creater interface {
		CreateBooking(ctx context.Context, b *deiz.Booking) error
	}
	Deleter interface {
		DeleteBooking(ctx context.Context, bookingID, clinicianID int) error
	}
	DeleterByDeleteID interface {
		DeleteBookingByDeleteID(ctx context.Context, deleteID string) error
	}
	Updater interface {
		UpdateBooking(ctx context.Context, b *deiz.Booking) error
	}
	GetterByID interface {
		GetBookingByID(ctx context.Context, bookingID int) (deiz.Booking, error)
	}
	GetterByDeleteID interface {
		GetBookingByDeleteID(ctx context.Context, deleteID string) (deiz.Booking, error)
	}
	OverlappingBlockedDeleter interface {
		DeleteOverlappingBlockedBooking(ctx context.Context, start, end time.Time, clinicianID int) error
	}
	BookingSlotAvailableChecker interface {
		IsBookingSlotAvailable(ctx context.Context, from, to time.Time, clinicianID int) (bool, error)
	}
	ToClinicianMailer interface {
		MailBookingToClinician(ctx context.Context, b *deiz.Booking, gCalLink string) error
	}
	BookingReminderMailer interface {
		MailBookingReminder(ctx context.Context, b *deiz.Booking, gCalLink, gMapsLink, cancelURL string) error
	}
	ToPatientMailer interface {
		MailBookingToPatient(ctx context.Context, b *deiz.Booking, gCalLink, gMapsLink, cancelURL string) error
	}
	CancelBookingToPatientMailer interface {
		MailCancelBookingToPatient(ctx context.Context, b *deiz.Booking) error
	}
	CancelBookingToClinicianMailer interface {
		MailCancelBookingToClinician(ctx context.Context, b *deiz.Booking) error
	}
	PatientCreater interface {
		CreatePatient(ctx context.Context, p *deiz.Patient, clinicianID int) error
	}
	PatientGetterByEmail interface {
		GetPatientByEmail(ctx context.Context, email string, clinicianID int) (deiz.Patient, error)
	}
	AddressGetterByID interface {
		GetAddressByID(ctx context.Context, addressID int) (deiz.Address, error)
	}
)

func NewAvailableBookingSlot(start, end time.Time, address deiz.Address, motive deiz.BookingMotive, clinicianID int) deiz.Booking {
	return deiz.Booking{
		Start:     start,
		End:       end,
		Address:   address,
		Remote:    address.ID == 0,
		Motive:    motive,
		Clinician: deiz.Clinician{ID: clinicianID},
	}
}

func (u *Usecase) GetFreeBookingSlots(ctx context.Context, start time.Time, tzName string, motiveID, motiveDuration, clinicianID int) ([]deiz.Booking, error) {
	officeHours, err := u.OfficeHoursGetter.GetClinicianOfficeHours(ctx, clinicianID)
	if err != nil {
		return nil, err
	}
	end := start.AddDate(0, 0, 7)
	bookings, err := u.ClinicianBookingsInTimeRangeGetter.GetClinicianBookingsInTimeRange(ctx, start, end, clinicianID)
	if err != nil {
		return nil, err
	}
	var freeSlots []deiz.Booking
	bookedSlotsTimeRanges := GetTimeRangesFromBookings(SortBookingsByStart(bookings), [][2]time.Time{})
	officeHoursTimeRange := GetAllOfficeHoursTimeRange(start, end, officeHours, u.loc)
	for i, timeRange := range officeHoursTimeRange {
		freeTimeRanges := GetTimeRangesNotOverLapping(motiveDuration, timeRange[0], timeRange[1], bookedSlotsTimeRanges, [][2]time.Time{})
		for _, timeRange := range freeTimeRanges {
			freeSlots = append(freeSlots,
				NewAvailableBookingSlot(timeRange[0], timeRange[1],
					officeHours[i].Address,
					deiz.BookingMotive{ID: motiveID, Duration: motiveDuration}, clinicianID))
		}
	}
	return freeSlots, nil
}

func (u *Usecase) GetBookingSlots(ctx context.Context, start time.Time, tzName string, defaultMotiveID, defaultMotiveDuration, clinicianID int) ([]deiz.Booking, error) {
	officeHours, err := u.OfficeHoursGetter.GetClinicianOfficeHours(ctx, clinicianID)
	if err != nil {
		return nil, err
	}
	end := start.AddDate(0, 0, 7)
	bookings, err := u.ClinicianBookingsInTimeRangeGetter.GetClinicianBookingsInTimeRange(ctx, start, end, clinicianID)
	if err != nil {
		return nil, err
	}
	bookedSlotsTimeRanges := GetTimeRangesFromBookings(SortBookingsByStart(bookings), [][2]time.Time{})
	officeHoursTimeRange := GetAllOfficeHoursTimeRange(start, end, officeHours, u.loc)
	for i, timeRange := range officeHoursTimeRange {
		freeTimeRanges := GetTimeRangesNotOverLapping(defaultMotiveDuration, timeRange[0], timeRange[1], bookedSlotsTimeRanges, [][2]time.Time{})
		for _, timeRange := range freeTimeRanges {
			bookings = append(bookings,
				NewAvailableBookingSlot(timeRange[0], timeRange[1],
					officeHours[i].Address,
					deiz.BookingMotive{ID: defaultMotiveID, Duration: defaultMotiveDuration}, clinicianID))
		}
	}
	return bookings, nil
}

func isBookingToClinician(b *deiz.Booking, clinicianID int) bool {
	return b.Clinician.ID == clinicianID
}

func (u *Usecase) BlockBookingSlot(ctx context.Context, b *deiz.Booking, clinicianID int) error {
	if !isBookingToClinician(b, clinicianID) {
		return deiz.ErrorUnauthorized
	}
	b.SetBlocked()
	return u.Creater.CreateBooking(ctx, b)
}

func (u *Usecase) UnlockBookingSlot(ctx context.Context, bookingID, clinicianID int) error {
	return u.Deleter.DeleteBooking(ctx, bookingID, clinicianID)
}

func (u *Usecase) ConfirmPreRegisteredBooking(ctx context.Context, b *deiz.Booking, clinicianID int, notifyPatient, notifyClinician bool) error {
	if !isBookingToClinician(b, clinicianID) {
		return deiz.ErrorUnauthorized
	}
	if bookingRegistrationInvalid(b) {
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
	if !isBookingToClinician(b, clinicianID) {
		return deiz.ErrorUnauthorized
	}
	if bookingPreRegistrationInvalid(b) {
		return deiz.ErrorStructValidation
	}
	available, err := u.BookingSlotAvailableChecker.IsBookingSlotAvailable(ctx, b.Start, b.End, clinicianID)
	if err != nil {
		return err
	}
	if !available {
		return deiz.ErrorBookingSlotAlreadyFilled
	}
	if err := u.OverlappingBlockedDeleter.DeleteOverlappingBlockedBooking(ctx, b.Start, b.End, clinicianID); err != nil {
		return err
	}
	if err := u.Creater.CreateBooking(ctx, b); err != nil {
		return err
	}
	return nil
}

func (u *Usecase) setBookingPatient(ctx context.Context, b *deiz.Booking) error {
	b.Patient.Sanitize()
	patient, err := u.PatientGetterByEmail.GetPatientByEmail(ctx, b.Patient.Email, b.Clinician.ID)
	if err != nil {
		return fmt.Errorf("unable to get patient by email: %s", err)
	}
	if !patient.IsSet() {
		err := u.PatientCreater.CreatePatient(ctx, &b.Patient, b.Clinician.ID)
		if err != nil {
			return fmt.Errorf("unable to create patient: %s", err)
		}
	} else {
		b.SetPatient(patient)
	}
	return nil
}

func (u *Usecase) setBookingAddress(ctx context.Context, b *deiz.Booking) error {
	if b.Address.IsSet() {
		var err error
		b.Address, err = u.AddressGetterByID.GetAddressByID(ctx, b.ID)
		if err != nil {
			return fmt.Errorf("unable to get address: %s", err)
		}
	}
	return nil
}

func (u *Usecase) RegisterPublicBooking(ctx context.Context, b *deiz.Booking) error {
	if bookingPreRegistrationInvalid(b) {
		return deiz.ErrorStructValidation
	}
	if err := u.setBookingPatient(ctx, b); err != nil {
		return err
	}
	if err := u.setBookingAddress(ctx, b); err != nil {
		return err
	}
	if err := u.Creater.CreateBooking(ctx, b); err != nil {
		return err
	}
	return u.NotifyBooking(ctx, b, b.Clinician.ID, true, true)
}

func bookingRegistrationValid(b *deiz.Booking) bool {
	return !(b.EndBeforeStart() ||
		b.ClinicianNotSet() ||
		b.Blocked ||
		!b.RemoteStatusMatchAddress())
}

func bookingRegistrationInvalid(b *deiz.Booking) bool {
	return !bookingRegistrationValid(b)
}

func bookingPreRegistrationValid(b *deiz.Booking) bool {
	return !(b.Confirmed || b.EndBeforeStart() || b.Blocked)
}

func bookingPreRegistrationInvalid(b *deiz.Booking) bool {
	return !bookingPreRegistrationValid(b)
}

func (u *Usecase) RegisterBooking(ctx context.Context, b *deiz.Booking, clinicianID int, notifyPatient, notifyClinician bool) error {
	if bookingRegistrationInvalid(b) || b.PatientNotSet() {
		return deiz.ErrorStructValidation
	}
	if !isBookingToClinician(b, clinicianID) {
		return deiz.ErrorUnauthorized
	}
	available, err := u.BookingSlotAvailableChecker.IsBookingSlotAvailable(ctx, b.Start, b.End, clinicianID)
	if err != nil {
		return fmt.Errorf("la disponibilité n'a pu être vérifiée : %s", err)
	}
	if !available {
		return deiz.ErrorBookingSlotAlreadyFilled
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
	if notifyClinician {
		if err := u.MailToClinician(ctx, b, clinicianID); err != nil {
			return err
		}
	}
	if notifyPatient {
		if err := u.MailToPatient(ctx, b, clinicianID); err != nil {
			return err
		}
	}
	return nil
}

func (u *Usecase) RemovePublicBooking(ctx context.Context, deleteID string) error {
	b, err := u.GetterByDeleteID.GetBookingByDeleteID(ctx, deleteID)
	if err != nil {
		return err
	}
	err = u.DeleterByDeleteID.DeleteBookingByDeleteID(ctx, deleteID)
	if err != nil {
		return err
	}
	return u.CancelToClinicianMailer.MailCancelBookingToClinician(ctx, &b)
}

func (u *Usecase) RemoveBooking(ctx context.Context, bookingID int, notifyPatient bool, clinicianID int) error {
	b, err := u.GetterByID.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}
	if err := u.Deleter.DeleteBooking(ctx, bookingID, clinicianID); err != nil {
		return err
	}
	if notifyPatient {
		return u.NotifyRemoveBookingToPatient(ctx, &b)
	}
	return nil
}

func (u *Usecase) NotifyRemoveBookingToPatient(ctx context.Context, b *deiz.Booking) error {
	if err := u.CancelToPatientMailer.MailCancelBookingToPatient(ctx, b); err != nil {
		return fmt.Errorf("unable to send cancel booking email to patient %s", err)
	}
	return nil
}

func (u *Usecase) MailToClinician(ctx context.Context, b *deiz.Booking, clinicianID int) error {
	return u.ToClinicianMailer.MailBookingToClinician(ctx,
		b, u.GCalendarLinkBuilder.BuildGCalendarLink(
			b.Start.In(u.loc), b.End.In(u.loc),
			fmt.Sprintf("Consultation avec %s %s", b.Patient.Surname, b.Patient.Name),
			b.Address.ToString(), b.Note,
		))
}

func createCancelURL(deleteID string) *url.URL {
	cancelURL, _ := url.Parse("https://deiz.fr")
	cancelURL.Path += "bookings/delete"
	params := url.Values{}
	params.Add("id", deleteID)
	cancelURL.RawQuery = params.Encode()
	return cancelURL
}

func (u *Usecase) createGCalendarLink(b *deiz.Booking) string {
	return u.GCalendarLinkBuilder.BuildGCalendarLink(
		b.Start.In(u.loc),
		b.End.In(u.loc),
		fmt.Sprintf("Consultation avec %s %s", b.Clinician.Surname, b.Clinician.Name),
		b.Address.ToString(),
		"",
	)
}

func (u *Usecase) MailToPatient(ctx context.Context, b *deiz.Booking, clinicianID int) error {
	cancelURL := createCancelURL(b.DeleteID)
	return u.ToPatientMailer.MailBookingToPatient(ctx, b,
		u.createGCalendarLink(b),
		u.GMapsLinkBuilder.BuildGMapsLink(b.Address.ToString()),
		cancelURL.String(),
	)
}

func (u *Usecase) SendReminders(ctx context.Context) error {
	start := resetTimeMn(time.Now().AddDate(0, 0, 2).UTC())
	end := start.Add(time.Hour * time.Duration(1)).UTC()

	bookings, err := u.BookingsInTimeRangeGetter.GetBookingsInTimeRange(ctx, start, end)
	if err != nil {
		return fmt.Errorf("unable to get bookings in time range: %s", err)
	}
	for _, b := range bookings {
		err := u.BookingReminderMailer.MailBookingReminder(ctx,
			&b, u.createGCalendarLink(&b), u.GMapsLinkBuilder.BuildGMapsLink(b.Address.ToString()), createCancelURL(b.DeleteID).String())
		if err != nil {
			return err
		}
	}
	return nil
}
