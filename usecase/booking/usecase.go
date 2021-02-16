package booking

import (
	"context"
	"github.com/audrenbdb/deiz"
	"time"
)

type repo interface {
	ClinicianOfficeHoursGetter
	BookingsInTimeRangeGetter
	Creater
	Deleter
	Updater
	GetterByID
	OverlappingBlockedDeleter
	CalendarSettingsGetter
}

type mail interface {
	ToClinicianMailer
	ToPatientMailer
	CancelBookingToPatientMailer
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
	OfficeHoursGetter         ClinicianOfficeHoursGetter
	BookingsInTimeRangeGetter BookingsInTimeRangeGetter
	Creater                   Creater
	Deleter                   Deleter
	GetterByID                GetterByID
	OverlappingBlockedDeleter OverlappingBlockedDeleter
	ToClinicianMailer         ToClinicianMailer
	ToPatientMailer           ToPatientMailer
	CancelToPatientMailer     CancelBookingToPatientMailer
	GCalendarLinkBuilder      GCalendarLinkBuilder
	GMapsLinkBuilder          GMapsLinkBuilder
	CalendarSettingsGetter    CalendarSettingsGetter
	Updater                   Updater
}

func NewUsecase(repo repo, mail mail, gMapsBuilder GMapsLinkBuilder, gCalBuilder GCalendarLinkBuilder) *Usecase {
	return &Usecase{
		OfficeHoursGetter:         repo,
		BookingsInTimeRangeGetter: repo,
		Creater:                   repo,
		Deleter:                   repo,
		Updater:                   repo,
		OverlappingBlockedDeleter: repo,
		CalendarSettingsGetter:    repo,
		GetterByID:                repo,

		ToClinicianMailer:     mail,
		ToPatientMailer:       mail,
		CancelToPatientMailer: mail,

		GCalendarLinkBuilder: gCalBuilder,
		GMapsLinkBuilder:     gMapsBuilder,
	}
}
