package deiz

import (
	"context"
	"time"
)

type CalendarSettings struct {
	ID            int           `json:"id" validate:"required"`
	DefaultMotive BookingMotive `json:"defaultMotive" validate:"required"`
	Timezone      Timezone      `json:"timezone" validate:"required"`
	Step          int           `json:"step" validate:"required"`
}

type Timezone struct {
	ID   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

//repo functions
type (
	calendarSettingsEditer interface {
		EditCalendarSettings(ctx context.Context, settings *CalendarSettings, clinicianID int) error
	}
	calendarSettingsGetter interface {
		GetCalendarSettings(ctx context.Context, clinicianID int) (CalendarSettings, error)
	}
	clinicianTimezoneGetter interface {
		GetClinicianTimezone(ctx context.Context, clinicianID int) (Timezone, error)
	}
)

type (
	//EditCalendarSettings edit clinician calendar settings
	EditCalendarSettings func(ctx context.Context, settings *CalendarSettings, clinicianID int) error
)

func editCalendarSettingsFunc(editer calendarSettingsEditer) EditCalendarSettings {
	return func(ctx context.Context, settings *CalendarSettings, clinicianID int) error {
		return editer.EditCalendarSettings(ctx, settings, clinicianID)
	}
}

func getClinicianTimezoneLoc(ctx context.Context, clinicianID int, getter clinicianTimezoneGetter) (*time.Location, error) {
	tz, err := getter.GetClinicianTimezone(ctx, clinicianID)
	if err != nil {
		return nil, err
	}
	return time.LoadLocation(tz.Name)
}
