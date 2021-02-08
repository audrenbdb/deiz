package deiz

import "time"

type GoogleCalendarService struct {
	LinkMaker GoogleCalendarLinkMaker
}

type (
	GoogleCalendarLinkMaker interface {
		MakeGoogleCalendarLink(start, end time.Time, title, location, details string) string
	}
)
