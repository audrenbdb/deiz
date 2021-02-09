package deiz

const ErrorUnauthorized Error = "unauthorized"
const ErrorStructValidation Error = "unable to validate struct"

type Error string

func (e Error) Error() string {
	return string(e)
}
