package gcal

import (
	"fmt"
	"net/url"
	"time"
)

type Event struct {
	Title    string
	Location string
	Details  string

	Start time.Time
	End   time.Time
}

func NewLink(event Event) string {
	startStr := fmt.Sprintf("%d%02d%02dT%02d%02d00", event.Start.Year(), event.Start.Month(), event.Start.Day(), event.Start.Hour(), event.Start.Minute())
	endStr := fmt.Sprintf("%d%02d%02dT%02d%02d00", event.End.Year(), event.End.Month(), event.End.Day(), event.End.Hour(), event.End.Minute())
	baseURL, _ := url.Parse("https://calendar.google.com")
	baseURL.Path += "calendar/event"
	params := url.Values{}
	params.Add("action", "TEMPLATE")
	params.Add("dates", fmt.Sprintf("%s/%s", startStr, endStr))
	params.Add("text", event.Title)
	params.Add("details", event.Details)
	params.Add("location", event.Location)

	baseURL.RawQuery = params.Encode()
	return baseURL.String()
}
