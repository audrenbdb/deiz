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
	patientGetter interface {
		GetPatientByEmail(ctx context.Context, email string, clinicianID int) (deiz.Patient, error)
	}
	patientCreater interface {
		CreatePatient(ctx context.Context, p *deiz.Patient, clinicianID int) error
	}
)

type RegisterUsecase struct {
	loc *time.Location

	patientGetter  patientGetter
	patientCreater patientCreater

	bookingCreater bookingCreater
	bookingUpdater bookingUpdater
	bookingGetter  clinicianBookingsInTimeRangeGetter

	bookingMailer bookingMailer
}

type RegisterDeps struct {
	Loc            *time.Location
	PatientGetter  patientGetter
	PatientCreater patientCreater
	BookingCreater bookingCreater
	BookingUpdater bookingUpdater
	BookingGetter  clinicianBookingsInTimeRangeGetter
	BookingMailer  bookingMailer
}

func NewRegisterUsecase(deps RegisterDeps) *RegisterUsecase {
	return &RegisterUsecase{
		loc:            deps.Loc,
		patientGetter:  deps.PatientGetter,
		patientCreater: deps.PatientCreater,
		bookingCreater: deps.BookingCreater,
		bookingUpdater: deps.BookingUpdater,
		bookingGetter:  deps.BookingGetter,
		bookingMailer:  deps.BookingMailer,
	}
}

func (r *RegisterUsecase) RegisterBookingFromPatient(ctx context.Context, b *deiz.Booking) error {
	if err := r.setBookingPatient(ctx, b); err != nil {
		return err
	}
	if r.registrationInvalid(b, b.Clinician.ID) {
		return deiz.ErrorStructValidation
	}
	if err := r.bookingCreater.CreateBooking(ctx, b); err != nil {
		return err
	}
	return r.notifyRegistration(b, true, true)
}

func (r *RegisterUsecase) RegisterBookingFromClinician(ctx context.Context, b *deiz.Booking, clinicianID int, notifyPatient bool) error {
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
	return r.notifyRegistration(b, notifyPatient, false)
}

func (r *RegisterUsecase) RegisterPreRegisteredBooking(ctx context.Context, b *deiz.Booking, clinicianID int, notifyPatient bool) error {
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
	return r.notifyRegistration(b, notifyPatient, false)

}

func (r *RegisterUsecase) setBookingPatient(ctx context.Context, b *deiz.Booking) error {
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

func (r *RegisterUsecase) notifyRegistration(b *deiz.Booking, notifyPatient, notifyClinician bool) error {
	if notifyClinician {
		if err := r.bookingMailer.MailBookingToClinician(b); err != nil {
			return err
		}
	}
	if notifyPatient {
		if err := r.bookingMailer.MailBookingToPatient(b); err != nil {
			return err
		}
	}
	return nil
}

func (r *RegisterUsecase) registrationValid(b *deiz.Booking, clinicianID int) bool {
	return b.Start.Before(b.End) &&
		b.ClinicianSet() && !b.Blocked && b.RemoteStatusMatchAddress() &&
		b.PatientSet() && b.Clinician.ID == clinicianID
}

func (r *RegisterUsecase) registrationInvalid(b *deiz.Booking, clinicianID int) bool {
	return !r.registrationValid(b, clinicianID)
}
