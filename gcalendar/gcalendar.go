//Package gcalendar builds gcalendar url
package gcalendar

import (
	"fmt"
	"net/url"
)

//Event of google calendar
type Event struct {
	/**
	 * Start of event, acceptable formats are:
	 *
	 * 20200316T010000Z - UTC
	 *
	 * 20200316T010000 - Time local to the user
	 *
	 * 20200316 - All day event
	 */
	Start string
	/**
	 * End of event, acceptable formats are:
	 *
	 * 20200316T010000Z - UTC
	 *
	 * 20200316T010000 - Time local to the user
	 *
	 * 20200316 - All day event
	 */
	End string
	//Title of the event
	Title string
	//Location of the event
	Location string
	//Details of the event
	Details string
}

//NewEventURL generates a google calendar event URL
func NewEventURL(args Event) string {
	baseURL, _ := url.Parse("https://calendar.google.com")
	baseURL.Path += "calendar/event"
	params := url.Values{}
	params.Add("action", "TEMPLATE")
	params.Add("dates", fmt.Sprintf("%s/%s", args.Start, args.End))
	params.Add("text", args.Title)
	params.Add("details", args.Details)
	params.Add("location", args.Location)

	baseURL.RawQuery = params.Encode()
	return baseURL.String()
}
