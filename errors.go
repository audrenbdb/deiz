package deiz

const ErrorUnauthorized Error = "unauthorized"
const ErrorStructValidation Error = "unable to validate struct"
const ErrorBookingSlotAlreadyFilled Error = "booking slot already filled"
const ErrorParsingTimezone Error = "unable to parse timzeone"

type Error string

func (e Error) Error() string {
	return string(e)
}
