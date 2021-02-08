package deiz_test

type (
	mockGMapsLinkMaker struct{}
)

func (r *mockGMapsLinkMaker) MakeGoogleMapsLink(address string) string {
	return ""
}
