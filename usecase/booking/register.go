package booking

import (
	"context"
	"github.com/audrenbdb/deiz"
	"time"
)

type (
	bookingUpdater interface {
		UpdateBooking(ctx context.Context, b *deiz.Booking) error
	}
	gMapsLinkBuilder interface {
		BuildGMapsLink(addressStr string) string
	}
	patientGetter interface {
		GetPatientByEmail(ctx context.Context, email string, clinicianID int) (deiz.Patient, error)
	}
	patientCreater interface {
		CreatePatient(ctx context.Context, p *deiz.Patient, clinicianID int) error
	}
)

type register struct {
	loc *time.Location

	patientGetter  patientGetter
	patientCreater patientCreater

	bookingCreater bookingCreater
	bookingUpdater bookingUpdater
	bookingGetter  clinicianBookingsInTimeRangeGetter

	bookingMailer        bookingMailer
	gCalendarLinkBuilder gCalendarLinkBuilder
	gMapsLinkBuilder     gMapsLinkBuilder
}

func NewRegisterUsecase(
	loc *time.Location,
	patientGetter patientGetter,
	patientCreater patientCreater,
	bookingCreater bookingCreater,
	bookingUpdater bookingUpdater,
	bookingGetter clinicianBookingsInTimeRangeGetter,
	bookingMailer bookingMailer,
	gCalendarLinkBuilder gCalendarLinkBuilder,
	gMapsLinkBuilder gMapsLinkBuilder,
) *register {
	return &register{
		loc:                  loc,
		patientGetter:        patientGetter,
		patientCreater:       patientCreater,
		bookingCreater:       bookingCreater,
		bookingUpdater:       bookingUpdater,
		bookingGetter:        bookingGetter,
		bookingMailer:        bookingMailer,
		gCalendarLinkBuilder: gCalendarLinkBuilder,
		gMapsLinkBuilder:     gMapsLinkBuilder,
	}
}

func (r *register) RegisterBookingFromPatient(ctx context.Context, b *deiz.Booking) error {
	if err := r.setBookingPatient(ctx, b); err != nil {
		return err
	}
	if r.registrationInvalid(b, b.Clinician.ID) {
		return deiz.ErrorStructValidation
	}
	if err := r.bookingCreater.CreateBooking(ctx, b); err != nil {
		return err
	}
	return r.notifyRegistration(ctx, b, true, true)
}

func (r *register) RegisterBookingFromClinician(ctx context.Context, b *deiz.Booking, clinicianID int, notifyPatient bool) error {
	if r.registrationInvalid(b, clinicianID) {
		return deiz.ErrorStructValidation
	}
	available, err := bookingSlotAvailable(ctx, b, r.bookingGetter)
	if err != nil {
		return err
	}
	if !available {
		return deiz.ErrorBookingSlotAlreadyFilled
	}
	if err := r.bookingCreater.CreateBooking(ctx, b); err != nil {
		return err
	}
	return r.notifyRegistration(ctx, b, notifyPatient, false)
}

func (r *register) RegisterPreRegisteredBooking(ctx context.Context, b *deiz.Booking, clinicianID int, notifyPatient bool) error {
	if r.registrationInvalid(b, clinicianID) {
		return deiz.ErrorStructValidation
	}
	available, err := bookingSlotAvailable(ctx, b, r.bookingGetter)
	if err != nil {
		return err
	}
	if !available {
		return deiz.ErrorBookingSlotAlreadyFilled
	}
	if err := r.bookingUpdater.UpdateBooking(ctx, b); err != nil {
		return err
	}
	return r.notifyRegistration(ctx, b, notifyPatient, false)

}

func (r *register) setBookingPatient(ctx context.Context, b *deiz.Booking) error {
	b.Patient.Sanitize()
	patient, err := r.patientGetter.GetPatientByEmail(ctx, b.Patient.Email, b.Clinician.ID)
	if err != nil {
		return err
	}
	if patient.IsNotSet() {
		err := r.patientCreater.CreatePatient(ctx, &b.Patient, b.Clinician.ID)
		if err != nil {
			return err
		}
	} else {
		b.Patient = patient
	}
	return nil
}

func (r *register) notifyRegistration(ctx context.Context, b *deiz.Booking, notifyPatient, notifyClinician bool) error {
	if notifyClinician {
		if err := mailBookingToClinician(ctx, b, r.loc, r.bookingMailer, r.gCalendarLinkBuilder); err != nil {
			return err
		}
	}
	if notifyPatient {
		if err := mailBookingToPatient(ctx, b, r.loc, r.bookingMailer, r.gCalendarLinkBuilder, r.gMapsLinkBuilder); err != nil {
			return err
		}
	}
	return nil
}

func (r *register) registrationValid(b *deiz.Booking, clinicianID int) bool {
	return b.Start.Before(b.End) &&
		b.ClinicianSet() && !b.Blocked && b.RemoteStatusMatchAddress() &&
		b.PatientSet() && b.Clinician.ID == clinicianID
}

func (r *register) registrationInvalid(b *deiz.Booking, clinicianID int) bool {
	return !r.registrationValid(b, clinicianID)
}
