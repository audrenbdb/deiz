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
	OverlappingBlockedDeleter
	CalendarSettingsGetter
}

type mail interface {
	ToClinicianMailer
	ToPatientMailer
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
	OverlappingBlockedDeleter OverlappingBlockedDeleter
	ToClinicianMailer         ToClinicianMailer
	ToPatientMailer           ToPatientMailer
	GCalendarLinkBuilder      GCalendarLinkBuilder
	GMapsLinkBuilder          GMapsLinkBuilder
	CalendarSettingsGetter    CalendarSettingsGetter
}

func NewUsecase(repo repo, mail mail, gMapsBuilder GMapsLinkBuilder, gCalBuilder GCalendarLinkBuilder) *Usecase {
	return &Usecase{
		OfficeHoursGetter:         repo,
		BookingsInTimeRangeGetter: repo,
		Creater:                   repo,
		Deleter:                   repo,
		OverlappingBlockedDeleter: repo,
		CalendarSettingsGetter:    repo,

		ToClinicianMailer: mail,
		ToPatientMailer:   mail,

		GCalendarLinkBuilder: gCalBuilder,
		GMapsLinkBuilder:     gMapsBuilder,
	}
}
