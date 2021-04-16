package account

import (
	"context"
	"errors"
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	mockCalendarSettingsUpdater struct {
		err error
	}
)

func (m *mockCalendarSettingsUpdater) UpdateCalendarSettings(ctx context.Context, s *deiz.CalendarSettings, clinicianID int) error {
	return m.err
}

func TestIsCalendarSettingsValid(t *testing.T) {
	var tests = []struct {
		description string
		inSettings  deiz.CalendarSettings
		valid       bool
	}{
		{
			description: "should succeed validating",
			inSettings: deiz.CalendarSettings{
				ID:       1,
				Timezone: deiz.Timezone{ID: 1},
			},
			valid: true,
		},
		{
			description: "should return false on validation",
			inSettings:  deiz.CalendarSettings{},
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			valid := isCalendarSettingsValid(&test.inSettings)
			assert.Equal(t, test.valid, valid)
		})
	}
}

func TestEditCalendarSettings(t *testing.T) {
	var tests = []struct {
		description string

		updater mockCalendarSettingsUpdater

		inSettings    deiz.CalendarSettings
		inClinicianID int

		outError error
	}{
		{
			description: "should fail to validate calendar settings",
			outError:    deiz.ErrorStructValidation,
		},
		{
			description: "should fail to update",
			inSettings: deiz.CalendarSettings{
				ID:       1,
				Timezone: deiz.Timezone{ID: 1},
			},
			updater:  mockCalendarSettingsUpdater{err: errors.New("failed to update")},
			outError: errors.New("failed to update"),
		},
		{
			description: "should fail to update",
			inSettings: deiz.CalendarSettings{
				ID:       1,
				Timezone: deiz.Timezone{ID: 1},
			},
			updater: mockCalendarSettingsUpdater{},
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			u := Usecase{
				CalendarSettingsUpdater: &test.updater,
			}
			err := u.EditCalendarSettings(context.Background(), &test.inSettings, test.inClinicianID)
			assert.Equal(t, test.outError, err)
		})
	}
}
