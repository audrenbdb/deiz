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
	Loc *time.Location

	PatientGetter  patientGetter
	PatientCreater patientCreater

	BookingCreater bookingCreater
	BookingUpdater bookingUpdater
	BookingGetter  bookingGetter

	BookingMailer bookingMailer
}

func (r *RegisterUsecase) RegisterBookingFromPatient(ctx context.Context, b *deiz.Booking) error {
	if err := r.setBookingPatient(ctx, b); err != nil {
		return err
	}
	return registerBookings(
		ctx, registrationDependencies{loc: r.Loc,
			getter: r.BookingGetter, creater: r.BookingCreater, mailer: r.BookingMailer},
		[]*deiz.Booking{b}, b.Clinician.ID, true, true)
}

func (r *RegisterUsecase) RegisterBookingsFromClinician(ctx context.Context, bookings []*deiz.Booking, clinicianID int, notifyPatient bool) error {
	return registerBookings(
		ctx, registrationDependencies{loc: r.Loc,
			getter: r.BookingGetter, creater: r.BookingCreater, mailer: r.BookingMailer},
		bookings, clinicianID, notifyPatient, false)
}

type registrationDependencies struct {
	getter  bookingGetter
	creater bookingCreater
	mailer  bookingMailer
	loc     *time.Location
}

func registerBookings(
	ctx context.Context, deps registrationDependencies,
	bookings []*deiz.Booking, clinicianID int, notifyPatient, notifyClinician bool) error {
	if areBookingsInvalid(bookings, clinicianID) {
		return deiz.ErrorStructValidation
	}
	for _, b := range bookings {
		available, err := bookingSlotAvailable(ctx, b, deps.getter, deps.loc)
		if err != nil {
			return err
		}
		if !available {
			return deiz.ErrorBookingSlotAlreadyFilled
		}
		if err := deps.creater.CreateBooking(ctx, b); err != nil {
			return err
		}
		if b.BookingType == deiz.AppointmentBooking {
			if err := notifyRegistration(b, deps.mailer, notifyPatient, notifyClinician); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *RegisterUsecase) RegisterPreRegisteredBooking(ctx context.Context, b *deiz.Booking, clinicianID int, notifyPatient bool) error {
	if areBookingsInvalid([]*deiz.Booking{b}, clinicianID) {
		return deiz.ErrorStructValidation
	}
	available, err := bookingSlotAvailable(ctx, b, r.BookingGetter, r.Loc)
	if err != nil {
		return err
	}
	if !available {
		return deiz.ErrorBookingSlotAlreadyFilled
	}
	if err := r.BookingUpdater.UpdateBooking(ctx, b); err != nil {
		return err
	}
	return notifyRegistration(b, r.BookingMailer, notifyPatient, false)

}

func (r *RegisterUsecase) setBookingPatient(ctx context.Context, b *deiz.Booking) error {
	b.Patient.Sanitize()
	patient, err := r.PatientGetter.GetPatientByEmail(ctx, b.Patient.Email, b.Clinician.ID)
	if err != nil {
		return err
	}
	if patient.IsNotSet() {
		if b.Patient.IsInvalid() {
			return deiz.ErrorStructValidation
		}
		err := r.PatientCreater.CreatePatient(ctx, &b.Patient, b.Clinician.ID)
		if err != nil {
			return err
		}
	} else {
		b.Patient = patient
	}
	return nil
}

func notifyRegistration(b *deiz.Booking, mailer bookingMailer, notifyPatient, notifyClinician bool) error {
	if notifyClinician {
		if err := mailer.MailBookingToClinician(b); err != nil {
			return err
		}
	}
	if notifyPatient && b.Patient.IsEmailSet() {
		if err := mailer.MailBookingToPatient(b); err != nil {
			return err
		}
	}
	return nil
}
