package booking

import (
	"context"
	"github.com/audrenbdb/deiz"
	"time"
)

type repo interface {
	ClinicianOfficeHoursGetter
	ClinicianBookingsInTimeRangeGetter
	BookingsInTimeRangeGetter
	Creater
	Deleter
	DeleterByDeleteID
	Updater
	GetterByID
	OverlappingBlockedDeleter
	CalendarSettingsGetter
	PatientCreater
	PatientGetterByEmail
	GetterByDeleteID
	BookingSlotAvailableChecker
	AddressGetterByID
}

type mail interface {
	ToClinicianMailer
	ToPatientMailer
	CancelBookingToPatientMailer
	CancelBookingToClinicianMailer
	BookingReminderMailer
}

type GCalendarLinkBuilder interface {
	BuildGCalendarLink(start, end time.Time, subject, addressStr, details string) string
}

type GMapsLinkBuilder interface {
	BuildGMapsLink(addressStr string) string
}

type CalendarSettingsGetter interface {
	GetClinicianCalendarSettings(ctx context.Context, clinicianID int) (deiz.CalendarSettings, error)
}

type Usecase struct {
	loc                                *time.Location
	OfficeHoursGetter                  ClinicianOfficeHoursGetter
	ClinicianBookingsInTimeRangeGetter ClinicianBookingsInTimeRangeGetter
	BookingsInTimeRangeGetter          BookingsInTimeRangeGetter
	Creater                            Creater
	Deleter                            Deleter
	DeleterByDeleteID                  DeleterByDeleteID
	GetterByID                         GetterByID
	OverlappingBlockedDeleter          OverlappingBlockedDeleter
	ToClinicianMailer                  ToClinicianMailer
	ToPatientMailer                    ToPatientMailer
	CancelToPatientMailer              CancelBookingToPatientMailer
	CancelToClinicianMailer            CancelBookingToClinicianMailer
	GCalendarLinkBuilder               GCalendarLinkBuilder
	GMapsLinkBuilder                   GMapsLinkBuilder
	CalendarSettingsGetter             CalendarSettingsGetter
	Updater                            Updater
	PatientCreater                     PatientCreater
	PatientGetterByEmail               PatientGetterByEmail
	GetterByDeleteID                   GetterByDeleteID
	BookingSlotAvailableChecker        BookingSlotAvailableChecker
	AddressGetterByID                  AddressGetterByID
	BookingReminderMailer              BookingReminderMailer
}

func NewUsecase(repo repo, mail mail, gMapsBuilder GMapsLinkBuilder, gCalBuilder GCalendarLinkBuilder, loc *time.Location) *Usecase {
	return &Usecase{
		loc:                                loc,
		OfficeHoursGetter:                  repo,
		ClinicianBookingsInTimeRangeGetter: repo,
		BookingsInTimeRangeGetter:          repo,
		Creater:                            repo,
		Deleter:                            repo,
		DeleterByDeleteID:                  repo,
		Updater:                            repo,
		OverlappingBlockedDeleter:          repo,
		CalendarSettingsGetter:             repo,
		GetterByID:                         repo,
		GetterByDeleteID:                   repo,
		PatientCreater:                     repo,
		PatientGetterByEmail:               repo,
		BookingSlotAvailableChecker:        repo,
		AddressGetterByID:                  repo,

		ToClinicianMailer:       mail,
		ToPatientMailer:         mail,
		CancelToPatientMailer:   mail,
		CancelToClinicianMailer: mail,
		BookingReminderMailer:   mail,

		GCalendarLinkBuilder: gCalBuilder,
		GMapsLinkBuilder:     gMapsBuilder,
	}
}
