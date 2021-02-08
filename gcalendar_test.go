package deiz_test

import "time"

type (
	mockGCalendarLinkMaker struct{}
)

func (r *mockGCalendarLinkMaker) MakeGoogleCalendarLink(start, end time.Time, title, location, details string) string {
	return ""
}
