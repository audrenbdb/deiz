/* build gmap url link to given location */
package gmaps

import (
	"net/url"
)

//NewQueryAddressURL generates a google maps search query
func NewQueryAddressURL(address string) string {
	baseURL, _ := url.Parse("https://www.google.com")
	baseURL.Path += "maps/search/"
	params := url.Values{}
	params.Add("api", "1")
	params.Add("query", address)

	baseURL.RawQuery = params.Encode()
	return baseURL.String()
}
