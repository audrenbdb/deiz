package settings

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	updater interface {
		UpdateCalendarSettings(ctx context.Context, s *deiz.CalendarSettings, clinicianID int) error
	}
)

type CalendarSettingsUsecase struct {
	SettingsUpdater updater
}

func (u *CalendarSettingsUsecase) EditCalendarSettings(ctx context.Context, s *deiz.CalendarSettings, clinicianID int) error {
	if s.IsInvalid() {
		return deiz.ErrorStructValidation
	}
	return u.SettingsUpdater.UpdateCalendarSettings(ctx, s, clinicianID)
}
