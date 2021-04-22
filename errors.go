package deiz

const GenericError Error = "fail"
const ErrorClinicianDoesNotExist = "ce clinicien n'existe pas"
const ErrorUnauthorized Error = "unauthorized"
const ErrorStructValidation Error = "unable to validate struct"
const ErrorBookingSlotAlreadyFilled Error = "Le créneau choisit en chevauche un autre"

type Error string

func (e Error) Error() string {
	return string(e)
}
