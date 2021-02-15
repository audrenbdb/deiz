//Package gcalendar builds gcalendar url
package gcalendar

import (
	"fmt"
	"net/url"
	"time"
)

type GCalendarService struct{}

//NewEventURL generates a google calendar event URL
func (s *GCalendarService) BuildGCalendarLink(start, end time.Time, title, location, details string) string {
	startStr := fmt.Sprintf("%d%02d%02dT%02d%02d00", start.Year(), start.Month(), start.Day(), start.Hour(), start.Minute())
	endStr := fmt.Sprintf("%d%02d%02dT%02d%02d00", end.Year(), end.Month(), end.Day(), end.Hour(), end.Minute())
	baseURL, _ := url.Parse("https://calendar.google.com")
	baseURL.Path += "calendar/event"
	params := url.Values{}
	params.Add("action", "TEMPLATE")
	params.Add("dates", fmt.Sprintf("%s/%s", startStr, endStr))
	params.Add("text", title)
	params.Add("details", details)
	params.Add("location", location)

	baseURL.RawQuery = params.Encode()
	return baseURL.String()
}

func NewService() *GCalendarService {
	return &GCalendarService{}
}
