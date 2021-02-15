package account

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	CalendarSettingsUpdater interface {
		UpdateCalendarSettings(ctx context.Context, s *deiz.CalendarSettings, clinicianID int) error
	}
)

func IsCalendarSettingsValid(s *deiz.CalendarSettings) bool {
	return s.ID != 0 && s.Timezone.ID != 0
}

func (u *Usecase) EditCalendarSettings(ctx context.Context, s *deiz.CalendarSettings, clinicianID int) error {
	if !IsCalendarSettingsValid(s) {
		return deiz.ErrorStructValidation
	}
	return u.CalendarSettingsUpdater.UpdateCalendarSettings(ctx, s, clinicianID)
}
