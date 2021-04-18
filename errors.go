package deiz

const ErrorUnauthorized Error = "unauthorized"
const ErrorStructValidation Error = "unable to validate struct"
const ErrorBookingSlotAlreadyFilled Error = "Le créneau choisit en chevauche un autre"

type Error string

func (e Error) Error() string {
	return string(e)
}
