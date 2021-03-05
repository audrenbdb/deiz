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
	DeleterByDeleteID
	Updater
	GetterByID
	OverlappingBlockedDeleter
	CalendarSettingsGetter
	PatientCreater
	PatientGetterByEmail
	GetterByDeleteID
}

type mail interface {
	ToClinicianMailer
	ToPatientMailer
	CancelBookingToPatientMailer
	CancelBookingToClinicianMailer
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
	DeleterByDeleteID         DeleterByDeleteID
	GetterByID                GetterByID
	OverlappingBlockedDeleter OverlappingBlockedDeleter
	ToClinicianMailer         ToClinicianMailer
	ToPatientMailer           ToPatientMailer
	CancelToPatientMailer     CancelBookingToPatientMailer
	CancelToClinicianMailer   CancelBookingToClinicianMailer
	GCalendarLinkBuilder      GCalendarLinkBuilder
	GMapsLinkBuilder          GMapsLinkBuilder
	CalendarSettingsGetter    CalendarSettingsGetter
	Updater                   Updater
	PatientCreater            PatientCreater
	PatientGetterByEmail      PatientGetterByEmail
	GetterByDeleteID          GetterByDeleteID
}

func NewUsecase(repo repo, mail mail, gMapsBuilder GMapsLinkBuilder, gCalBuilder GCalendarLinkBuilder) *Usecase {
	return &Usecase{
		OfficeHoursGetter:         repo,
		BookingsInTimeRangeGetter: repo,
		Creater:                   repo,
		Deleter:                   repo,
		DeleterByDeleteID:         repo,
		Updater:                   repo,
		OverlappingBlockedDeleter: repo,
		CalendarSettingsGetter:    repo,
		GetterByID:                repo,
		GetterByDeleteID:          repo,
		PatientCreater:            repo,
		PatientGetterByEmail:      repo,

		ToClinicianMailer:       mail,
		ToPatientMailer:         mail,
		CancelToPatientMailer:   mail,
		CancelToClinicianMailer: mail,

		GCalendarLinkBuilder: gCalBuilder,
		GMapsLinkBuilder:     gMapsBuilder,
	}
}
