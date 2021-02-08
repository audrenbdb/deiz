package deiz

const ErrorUnauthorized Error = "unauthorized"

type Error string

func (e Error) Error() string {
	return string(e)
}
