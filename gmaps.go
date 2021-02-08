package deiz

type GoogleMapsService struct {
	GoogleMapsLinkMaker GoogleMapsLinkMaker
}

type (
	GoogleMapsLinkMaker interface {
		MakeGoogleMapsLink(address string) string
	}
)
